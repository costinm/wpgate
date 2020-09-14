package mesh

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/costinm/wpgate/pkg/streams"
)

const (
	TopicConnectUP = "connectUP"
)

// Information about a node.
// Sent periodically, signed by the origin - for example as a JWT, or UDP
// proto
type NodeAnnounce struct {
	UA string `json:"UA,omitempty"`

	// Non-link local IPs from all interfaces. Includes public internet addresses
	// and Wifi IP4 address. Used to determine if a node is directly connected.
	IPs []*net.UDPAddr `json:"IPs,omitempty"`

	// Set if the node is an active Android AP.
	SSID string `json:"ssid,omitempty"`

	// True if the node is an active Android AP on the interface sending the message.
	// Will trigger special handling in link-local - if the receiving interface is also
	// an android client.
	AP bool `json:"AP,omitempty"`

	Ack bool `json:"ACK,omitempty"`

	// VIP of the direct parent, if this node is connected.
	// Used to determine the mesh topology.
	Vpn string `json:"Vpn,omitempty"`
}

// Node information, based on registration info or discovery.
// Map of nodes, keyed by interface address is stored in Gateway.nodes.
type DMNode struct {
	// VIP is the mesh specific IP6 address. The 'network' identifies the master node, the
	// link part is the sha of the public key. This is a byte[16].
	// Last 8 bytes as uint64 are the primary key in the map.
	VIP net.IP `json:"vip,omitempty"`

	// Pub
	PublicKey []byte `json:"pub,omitempty"`

	// Information from the node - from an announce or message.
	NodeAnnounce *NodeAnnounce

	Labels map[string]string `json:"l,omitempty"`

	Bacokff time.Duration `json:"-"`

	// Last LL GW address used by the peer.
	// Public IP addresses are stored in Reg.IPs.
	// If set, returned as the first address in GWs, which is used to connect.
	// This is not sent in the registration - but extracted from the request
	// remote address.
	GW *net.UDPAddr `json:"gw,omitempty"`

	// Set if the gateway has an active incoming connection from this
	// node, with the node acting as client.
	// Streams will be forwarded to the node using special 'accept' mode.
	// This is similar with PUSH in H2.
	TunSrv MuxSession `json:"-"`

	// Existing tun to the remote node, previously dialed.
	TunClient MuxSession `json:"-"`

	// IP4 address of last announce
	Last4 *net.UDPAddr `json:"-"`

	// IP6 address of last announce
	Last6 *net.UDPAddr `json:"-"`

	FirstSeen time.Time

	// Last packet or registration from the peer.
	LastSeen time.Time `json:"t"`

	// In seconds since first seen, last 100
	Seen []int `json:"-"`

	// LastSeen in a multicast announce
	LastSeen4 time.Time

	// LastSeen in a multicast announce
	LastSeen6 time.Time `json:"-"`

	// Number of multicast received
	Announces int

	// Numbers of announces received from that node on the P2P interface
	AnnouncesOnP2P int

	// Numbers of announces received from that node on the P2P interface
	AnnouncesFromP2P int
}

func NewDMNode() *DMNode {
	now := time.Now()
	return &DMNode{
		Labels:       map[string]string{},
		FirstSeen:    now,
		LastSeen:     now,
		NodeAnnounce: &NodeAnnounce{},
	}
}

// Dial a stream over a multiplexed connection.
//
// A node connects to a VPN or an intermediary node.
// The connectcan can multiplex multiple streams.
// This function connects a Stream object, using a multiplexed
// stream.
//
// In addition to the VPN, the node may have sessions (direct
// or proxied) to specific nodes.
//
// For example SSHClient, SSHServer, Quic can support this.
type MuxSession interface {
	// DialProxy will use the remote gateway to jump to
	// a different destination, indicated by stream.
	// On return, the stream ServerOut and ServerIn will be
	// populated, and connected to stream Dest.
	DialProxy(tp *streams.Stream) error

	// The VIP of the remote host, after authentication.
	RemoteVIP() net.IP

	// Wait for the stream to finish.
	Wait() error

	// RemoteAccept requests the other end to forward 'remote port' on the
	// gateway to a local port (forwardDest). This is primarily used
	// to create a 'plain TCP' listener on the other end.
	//
	// It is also useful with legacy SSH servers.
	//
	// SNI/HTTP/connect/etc do not use this mechanism.
	//
	// Equivalent with -R
	//
	// Using ":0" as remoteListener will open any port.
	//
	//RemoteAccept(remoteListenAddr, forwardDest string) error
}

// MUXDialer is implemented by a transport that can be
// used for egress for streams. SSHGate creating SSHClients is an example.
//
// On the server side, MuxSession are created when a client connects
type MUXDialer interface {
	// Dial one TCP/mux connection to the IP:port.
	// The destination is a mesh node - port typically 5222, or 22 for 'regular' SSH serves.
	//
	// After handshake, an initial message is sent, including informations about the current node.
	//
	// The remote can be a trusted VPN, an untrusted AP/Gateway, a peer (link local or with public IP),
	// or a child. The subsriptions are used to indicate what messages will be forwarded to the server.
	// Typically VPN will receive all events, AP will receive subset of events related to topology while
	// child/peer only receive directed messages.
	DialMUX(addr string, pub []byte, subs []string) (MuxSession, error)
}

// IPResolver uses DNS cache or lookups to return the name
// associated with an IP, for metrics/stats/logs
type IPResolver interface {
	IPResolve(ip string) string
}

// Textual representation of the node registration data.
func (n *DMNode) String() string {
	b, _ := json.Marshal(n)
	return string(b)
}

// Return the list of gateways for the node, starting with the link local if any.
func (n *DMNode) GWs() []*net.UDPAddr {
	res := []*net.UDPAddr{}

	if n.GW != nil {
		res = append(res, n.GW)
	}
	if n.Last4 != nil {
		res = append(res, n.Last4)
	}
	if n.Last6 != nil {
		res = append(res, n.Last6)
	}
	return res
}

// Called when receiving a registration or regular valid message via a different gateway.
// - HandleRegistrationRequest - after validating the VIP
//
//
// For VPN, the srcPort is assigned by the NAT, can be anything
// For direct, the port will be 5228 or 5229
func (n *DMNode) UpdateGWDirect(addr net.IP, zone string, srcPort int, onRes bool) {
	n.LastSeen = time.Now()
	n.GW = &net.UDPAddr{IP: addr, Port: srcPort, Zone: zone}
}
func (n *DMNode) BackoffReset() {
	n.Bacokff = 0
}
func (n *DMNode) BackoffSleep() {
	if n.Bacokff == 0 {
		n.Bacokff = 5 * time.Second
	}
	time.Sleep(n.Bacokff)
	if n.Bacokff < 5*time.Minute {
		n.Bacokff = n.Bacokff * 2
	}
}

// Track one interface.
type ActiveInterface struct {
	// Interface name. Name containing 'p2p' results in specific behavior.
	Name string

	// IP6 link local address. May be nil if IPPub is set.
	// One or the other must be set.
	IP6LL net.IP

	// IP4 address - may be a routable address, nil or private address.
	// If public address - may be included in the register, but typically not
	// useful.
	IP4 net.IP

	// Public addresses. IP6 address may be used for direct connections (in some
	// cases)
	IPPub []net.IP

	// Port for the UDP unicast link-local listener.
	Port int
	// Port for the UDP unicast link-local listener.
	Port4 int

	// True if this interface is an Android AP
	AndroidAP bool

	// True if this interface is connected to an Android DM node.
	AndroidAPClient bool
}

type ScanResults struct {
	// Visible devices at this moment
	Scan []*MeshDevice `json:"scan,omitempty"`

	Stats string `json:"stat,omitempty"`

	// Visible wifi networks (all kinds)
	Visible int `json:"visible,omitempty"`

	// My SSID and PSK
	SSID          string `json:"s,omitempty"`
	PSK           string `json:"p,omitempty"`
	ConnectedWifi string `json:"w,omitempty"`
	Freq          int    `json:"f,omitempty"`
	Level         int    `json:"l,omitempty"`
}

// WifiRegistrationInfo contains information about the wifi node sent to the
// other nodes, to sync up visibility info.
//
type WifiRegistrationInfo struct {
	// Visible P2P devices in the mesh. This includes active APs as well as devices announcing via
	// BLE or NAN (or other means).
	Devices map[string]*MeshDevice `json:"devices,omitempty"`

	SSID string `json:"ssid,omitempty"`
	PSK  string `json:"psk,omitempty"`

	// Network we are connected to.
	// TODO: In case of chained P2P networks, should be either the path, or a separate field should include the path
	// and the net should be the 'top level' network of the root.
	Net string `json:"net,omitempty"`

	// Number of visible wifi networks (all kinds)
	VisibleWifi int `json:"scanCnt,omitempty"`
}

// Info about a device from the P2P info.
type MeshDevice struct {
	SSID string `json:"s,omitempty"`
	PSK  string `json:"p,omitempty"`

	// MAC is used with explicit P2P connect ( i.e. no hacks )
	// User input required on the receiving end ( PBC )
	MAC string `json:"d,omitempty"`

	Name string `json:"N,omitempty"`

	// Set only if the device is currently visible in scan
	Level int `json:"l,omitempty"`
	Freq  int `json:"f,omitempty"`

	// Extracted from DIRECT DNSSD
	UserAgent string `json:"ua,omitempty"`
	Net       string `json:"n,omitempty"`

	Cap   string `json:"c,omitempty"`
	BSSID string `json:"b,omitempty"`

	LastSeen time.Time `json:"lastSeen,omitempty"`

	Self int `json:"self,omitempty"`
	// Only on supplicant,not on android
	ServiceUpdateInd int `json:"sui,omitempty"`
}

func (md *MeshDevice) String() string { return fmt.Sprintf("%s/%d", md.SSID, md.Level) }

// Keyed by Hostname:port (if found in dns tables) or IP:port
type HostStats struct {
	// First open
	Open time.Time

	// Last usage
	Last time.Time

	SentBytes   int
	RcvdBytes   int
	SentPackets int
	RcvdPackets int
	Count       int

	LastLatency time.Duration
	LastBPS     int
}

type ListenerConf struct {
	// Local address (ex :8080). This is the requested address - if busy :0 will be used instead, and Port
	// will be the actual port
	// TODO: UDS
	// TODO: indicate TLS SNI binding.
	Local string

	// Real port the listener is listening on, or 0 if the listener is not bound to a port (virtual, using mesh).
	Port int

	// Remote where to forward the proxied connections
	// IP:port format, where IP can be a mesh VIP
	Remote string `json:"Remote,omitempty"`
}


type Host struct {
	// Address and port of a HTTP server to forward the domain.
	Addr string

	// Directory to serve static files. Used if Addr not set.
	Dir string
	Mux http.Handler `json:"-"`
}

// Configuration for the Gateway.
//
type GateCfg struct {

	// Port proxies: will register a listener for each port, forwarding to the
	// given address.
	Listeners []*ListenerConf `json:"TcpProxy,omitempty"`

	// Set of hosts with certs to configure in the h2 server.
	// The cert is expected in CertDir/HOSTNAME.[key,crt]
	// The server will terminate TLS and HTTP, forward to the host as plain text.
	Hosts map[string]*Host `json:"Hosts,omitempty"`
}


// UdpWriter is implemented by capture, provides a way to send back packets to
// the captured app.
type UdpWriter interface {
	WriteTo(data []byte, dstAddr *net.UDPAddr, srcAddr *net.UDPAddr) (int, error)
}
