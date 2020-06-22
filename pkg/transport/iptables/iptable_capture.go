package iptables

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"unsafe"

	"github.com/costinm/dmesh/dm/gate"
	"github.com/costinm/wpgate/pkg/mesh"
	"golang.org/x/net/ipv4"
	"golang.org/x/sys/unix"
)

//

// Status:
// - TCP capture with redirect works
// - capture with TPROXY is not possible - TPROXY is on the PREROUTING chain only,
// not touched by output packets.
//   https://upload.wikimedia.org/wikipedia/commons/3/37/Netfilter-packet-flow.svg
//
// It works great for transparent proxy in a gateway/router - however same can be also done using the TUN
// and routing to the TUN using iptables or other means.
//
// TODO: combination of tun and tproxy to avoid user-space IP ?
//

// Initialize a port as a TPROXY socket. This can be sent over UDS from the root, and used for
// UDP capture.
func StartUDPTProxyListener(port int) (*os.File, error) {
	// TPROXY mode for UDP - alternative is to use REDIRECT and parse
	// /proc/net/nf_conntrack
	s, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}

	err = unix.SetsockoptInt(s, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
	if err != nil {
		return nil, err
	}

	// NET_CAP
	err = unix.SetsockoptInt(s, unix.SOL_IP, unix.IP_TRANSPARENT, 1)
	if err != nil {
		fmt.Println("TRANSPARENT err ", err)
		//return err
	}

	err = unix.SetsockoptInt(s, unix.SOL_IP, unix.IP_FREEBIND, 1)
	if err != nil {
		fmt.Println("FREEBIND err ", err)
		//return err
	}

	err = unix.SetsockoptInt(s, unix.IPPROTO_IP, unix.IP_RECVORIGDSTADDR, 1)
	if err != nil {
		return nil, err
	}
	log.Println("Openned TPROXY capture port in TRANSPARENT mode ", port)
	err = unix.Bind(s, &unix.SockaddrInet4{
		Port: port,
	})
	if err != nil {
		log.Println("Error binding ", err)
		return nil, err
	}

	f := os.NewFile(uintptr(s), "TProxy")
	return f, nil
}

// Start listening for TCP and UDP addresses.
// For TCP, will use the OriginalDst
func StartIstioCapture(p *mesh.Gateway, addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go iptablesServeConn(p, conn)
		}
	}()

	var f *os.File
	if p.UDPListener == nil {
		// Requires root, not starting UDP proxy
		//return nil
		f, err = StartUDPTProxyListener(15006)
		if err != nil {
			log.Println("Error starting TPROXY", err)
			return err
		}
	} else {
		f = p.UDPListener
	}
	c, err := net.FileConn(f)

	lu, ok := c.(*net.UDPConn)
	if !ok {
		return errors.New("failed to cast")
	}

	p.UDPWriter = &transparentUdp{con: lu}

	//lu, err := net.ListenUDP("udp", &net.UDPAddr{
	//	Port: 15000,
	//})
	//if err != nil {
	//	return err
	//}
	//int flags = 1;
	//ret = setsockopt(udp_socket, IPPROTO_IP, IP_RECVORIGDSTADDR, &flags, sizeof(flags));

	go func() {
		data := make([]byte, 1600)
		oob := ipv4.NewControlMessage(ipv4.FlagDst)
		//oob := make([]byte, 256)
		for {

			n, noob, _, addr, err := lu.ReadMsgUDP(data[0:], oob)
			if err != nil {
				continue
			}

			cm4, err := syscall.ParseSocketControlMessage(oob[0:noob])
			origPort := uint16(0)
			var origIP net.IP
			for _, cm := range cm4 {
				if cm.Header.Type == unix.IP_RECVORIGDSTADDR {
					// \attention: IPv4 only!!!
					// address type, 1 - IPv4, 4 - IPv6, 3 - hostname, only IPv4 is supported now
					rawaddr := make([]byte, 4)
					// raw IP address, 4 bytes for IPv4 or 16 bytes for IPv6, only IPv4 is supported now
					copy(rawaddr, cm.Data[4:8])
					origIP = net.IP(rawaddr)

					// Bigendian is the network bit order, seems to be used here.
					origPort = binary.BigEndian.Uint16(cm.Data[2:])

				}
			}
			//if cm4.Parse(oob) == nil {
			//dst = cm4.Dst
			//}
			//log.Printf("NOOB %d %d %V %x", noob, flags, cm4, oob[0:noob])
			//if ((cmsg->cmsg_level == SOL_IP) && (cmsg->cmsg_type == IP_RECVORIGDSTADDR))
			//{
			//	memcpy (&dstaddr, CMSG_DATA(cmsg), sizeof (dstaddr));
			//	dstaddr.sin_family = AF_INET;
			//}
			go gate.HandleUdp(p, origIP, origPort, addr.IP, uint16(addr.Port), data[0:n])
		}
	}()
	return nil
}

//
type transparentUdp struct {
	con *net.UDPConn
}

// UDP write with source address control.
func (tudp *transparentUdp) WriteTo(data []byte, dstAddr *net.UDPAddr, srcAddr *net.UDPAddr) (int, error) {

	// Attempt to write as UDP
	cm4 := new(ipv4.ControlMessage)
	cm4.Src = srcAddr.IP
	oob := cm4.Marshal()
	n, _, err := tudp.con.WriteMsgUDP(data, oob, dstAddr)
	if err != nil {
		n, err = tudp.con.WriteToUDP(data, dstAddr)
		if err != nil {
			log.Print("Failed to send DNS ", dstAddr, srcAddr)
		}
	}

	return n, err // tudp.con.WriteTo(data, dstAddr)
}

// ServeConn is used to serve a single TCP UdpNat.
// See https://github.com/cybozu-go/transocks
// https://github.com/ryanchapman/go-any-proxy/blob/master/any_proxy.go,
// and other examples.
// Based on REDIRECT.
func iptablesServeConn(p *mesh.Gateway, conn net.Conn) error {
	addr, port, conn1, err := getOriginalDst(conn.(*net.TCPConn))
	if err != nil {
		conn.Close()
		return err
	}

	iaddr := net.IP(addr)

	proxy := p.NewTcpProxy(conn1.RemoteAddr(), "IPT", nil, conn1, conn1)
	defer proxy.Close()

	err = proxy.Dial("", &net.TCPAddr{IP: iaddr, Port: int(port)})
	if err != nil {
		return err
	}
	proxy.Proxy()

	return nil
}

const (
	SO_ORIGINAL_DST      = 80
	IP6T_SO_ORIGINAL_DST = 80
)

// Should be used only for REDIRECT capture.
func getOriginalDst(clientConn *net.TCPConn) (rawaddr []byte, port uint16, newTCPConn *net.TCPConn, err error) {

	if clientConn == nil {
		err = errors.New("ERR: clientConn is nil")
		return
	}

	// test if the underlying fd is nil
	remoteAddr := clientConn.RemoteAddr()
	if remoteAddr == nil {
		err = errors.New("ERR: clientConn.fd is nil")
		return
	}

	//srcipport := fmt.Sprintf("%v", clientConn.RemoteAddr())

	newTCPConn = nil
	// net.TCPConn.File() will cause the receiver's (clientConn) socket to be placed in blocking mode.
	// The workaround is to take the File returned by .File(), do getsockopt() to get the original
	// destination, then create a new *net.TCPConn by calling net.Conn.FileConn().  The new TCPConn
	// will be in non-blocking mode.  What a pain.
	clientConnFile, err := clientConn.File()
	if err != nil {
		//common.Errorf("GETORIGINALDST|%v->?->FAILEDTOBEDETERMINED|ERR: could not get a copy of the client UdpNat's file object", srcipport)
		return
	} else {
		clientConn.Close()
	}

	// Get original destination
	// this is the only syscall in the Golang libs that I can find that returns 16 bytes
	// Example result: &{Multiaddr:[2 0 31 144 206 190 36 45 0 0 0 0 0 0 0 0] Interface:0}
	// port starts at the 3rd byte and is 2 bytes long (31 144 = port 8080)
	// IPv6 version, didn't find a way to detect network family
	//addr, err := syscall.GetsockoptIPv6Mreq(int(clientConnFile.Fd()), syscall.IPPROTO_IPV6, IP6T_SO_ORIGINAL_DST)
	// IPv4 address starts at the 5th byte, 4 bytes long (206 190 36 45)
	addr, err := syscall.GetsockoptIPv6Mreq(int(clientConnFile.Fd()), syscall.IPPROTO_IP, SO_ORIGINAL_DST)
	if err != nil {
		return
	}
	newConn, err := net.FileConn(clientConnFile)
	if err != nil {
		return
	}
	if _, ok := newConn.(*net.TCPConn); ok {
		newTCPConn = newConn.(*net.TCPConn)
		clientConnFile.Close()
	} else {
		errmsg := fmt.Sprintf("ERR: newConn is not a *net.TCPConn, instead it is: %T (%v)", newConn, newConn)
		err = errors.New(errmsg)
		return
	}

	// \attention: IPv4 only!!!
	// address type, 1 - IPv4, 4 - IPv6, 3 - hostname, only IPv4 is supported now
	rawaddr = make([]byte, 4)
	// raw IP address, 4 bytes for IPv4 or 16 bytes for IPv6, only IPv4 is supported now
	copy(rawaddr, addr.Multiaddr[4:8])

	// Bigendian is the network bit order, seems to be used here.
	port = binary.BigEndian.Uint16(addr.Multiaddr[2:])

	return
}

func isLittleEndian() bool {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	return (b == 0x04)
}

var (
	NativeOrder binary.ByteOrder
)

func init() {
	if isLittleEndian() {
		NativeOrder = binary.LittleEndian
	} else {
		NativeOrder = binary.BigEndian
	}
}
