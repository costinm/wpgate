package msgs

import (
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/costinm/wpgate/pkg/auth"
	"github.com/costinm/wpgate/pkg/mesh"
)

type Backoff interface {
	BackoffSleep()
	BackoffReset()
}

var ReceiveBaseUrl = "https://127.0.0.1:5228/"


// Used to push a message from a remote sender.
//
// Mapped to /s/[DESTID]?...
// Local
//
// q or path can be used to pass command. Body and query string are sent.
// TODO: compatibility with cloud events and webpush
// TODO: RBAC (including admin check for system notifications)
//
func HTTPHandlerSend(w http.ResponseWriter, r *http.Request) {
	//transport.GetPeerCertBytes(r)

	r.ParseForm()

	var cmd string
	var parts []string
	q := r.Form.Get("q")

	if q != "" {
		parts = strings.Split(q, "/")
		cmd = q
	} else {
		parts = strings.Split(r.URL.Path, "/")
		parts = parts[2:]
		cmd = strings.Join(parts, " ")

		log.Println("MSG_SEND: ", parts, "--", cmd)
	}

	params := map[string]string{}
	for k, v := range r.Form {
		params[k] = v[0]
	}
	var err error
	var body []byte
	if r.Method == "POST" {
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
	}

	DefaultMux.HandleMessageForNode(NewMessage(cmd, params).SetDataJSON(body))
	w.WriteHeader(200)
}

var SharedWPAuth = []byte{1}

// Webpush handler.
// Auth: VAPID or client cert - results in VIP of sender
//
func (mux *Mux) HTTPHandlerWebpush(w http.ResponseWriter, r *http.Request) {
	h2ctx := r.Context().Value(mesh.H2Info).(*mesh.ReqContext)

	parts := strings.Split(r.RequestURI, "/")
	if len(parts) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid URL"))
		return
	}

	log.Println("WEBPUSH over HTTP ", parts)

	dest := parts[2]
	if dest == "" || dest == mux.Auth.Name || dest == mux.Auth.VIP6.String() {
		ec := auth.NewContextUA(mux.Auth.Priv, mux.Auth.Pub, SharedWPAuth)
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Body error"))
			return
		}

		authH := r.Header.Get("Authorization")
		if authH == "" {
			log.Println("Missing auth ", r.Header)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Missing auth"))
			return
		}
		jwt, pub, err := auth.CheckVAPID(authH, time.Now())

		msgb, err := ec.Decrypt(b)
		if err != nil {
			log.Println("Decrypt error ", err)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Decrypt error"))
			return
		}

		ev := mux.ProcessMessage(msgb, h2ctx)
		log.Println("GOT WEBPUSH: ", jwt, string(msgb), ev)

		if ev == nil {
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Invalid format"))
				return
			}
		}

		role := mux.Auth.Authorized[string(pub)]
		if role == "" {
			log.Println("Unauthorized ")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Unauthorized"))
			return
		}

		mux.HandleMessageForNode(ev)
	} else {
		// Dest is remote, we're just forwarding

	}

	w.WriteHeader(201)
}

// Currently mapped to /dmesh/uds - sends a message to a specific connection, defaults to the UDS connection
// to the android or root dmwifi app.
func (mux *Mux) HTTPUDS(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var cmd string
	var parts []string
	q := r.Form.Get("q")
	h := r.Form.Get("h")
	if h == "" {
		h = "dmesh"
	}

	if q != "" {
		parts = strings.Split(q, " ")
		cmd = q
	} else {
		parts = strings.Split(r.URL.Path, "/")
		parts = parts[3:]
		cmd = strings.Join(parts, " ")

		log.Println("UDS: ", parts, "--", cmd)
	}

	params := map[string]string{}
	for k, v := range r.Form {
		params[k] = v[0]
	}
	var err error
	var body []byte
	if r.Method == "POST" {
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
	}

	ch := mux.connections[h]
	if ch != nil {
		ch.SendMessageToRemote(NewMessage(cmd, params).SetDataJSON(body))
		w.WriteHeader(200)
	} else {
		w.WriteHeader(404)
		return
	}
}

// MonitorEvents will connect to a mesh address and monitor the messages.
//
// base is used for forwarding.
//
//func (w *Mux) MonitorEvents(node Backoff, idhex string, path []string) {
//	hc := transport.NewSocksHttpClient("")
//	hc.Timeout = 1 * time.Hour
//
//	if idhex == "" {
//		hc = http.DefaultClient
//	}
//
//	for {
//		t0 := time.Now()
//
//		err := w.MonitorNode(hc, idhex, path)
//		if err != nil {
//			log.Println("WATCH_ERR", idhex, err, time.Since(t0))
//			node.BackoffSleep()
//			continue
//		}
//		node.BackoffReset()
//
//		log.Println("WATCH_CLOSE", idhex, time.Since(t0))
//		node.BackoffSleep()
//	}
//
//}

// UA represents a "user agent" - or client using the webpush protocol
type UA struct {
	// URL of the subscribe for the push service
	PushService string
}

// Create a subscription, using the Webpush standard protocol.
//
// URL is "/subscribe", no header required ( but passing a VAPID or mtls),
// response in 'location' for read and Link for sub endpoint.
func (ua *UA) Subscribe() (sub *auth.Subscription, err error) {
	res, err := http.Post(ua.PushService+"/subscribe", "text/plain", nil)

	if err != nil {
		return
	}
	sub = &auth.Subscription{}
	sub.Location = res.Header.Get("location")
	links := textproto.MIMEHeader(res.Header)["Link"]
	for _, l := range links {
		for _, link := range strings.Split(l, ",") {
			parts := strings.Split(link, ";")
			if len(parts) > 1 &&
				strings.TrimSpace(parts[1]) == "rel=\"urn:ietf:params:push\"" {
				sub.Endpoint = parts[0]
			}
		}
	}

	// generate encryption key and authenticator

	return
}

// Subscribe creates a subscription. Initial version is just a
// random - some interface will be added later, to allow sets.
func SubscribeHandler(res http.ResponseWriter, req *http.Request) {
	// For simple testing we ignore sender auth, as well as subscription sets
	token := make([]byte, 16)
	rand.Read(token)

	id := base64.RawURLEncoding.EncodeToString(token)

	res.WriteHeader(201)

	// TODO: try to use a different server, to verify UA is
	// parsing both

	// Used for send - on same server as subscribe
	res.Header().Add("Link", "</p/"+
		id+
		">;rel=\"urn:ietf:params:push\"")

	// May provide support for set: should be enabled if a
	// set interface is present, want to test without set as well
	//res.Header().Add("Link", "</p/" +
	//	"JzLQ3raZJfFBR0aqvOMsLrt54w4rJUsV" +
	//	">;rel=\"urn:ietf:params:push:set\"")

	res.Header().Add("Location", ReceiveBaseUrl+"/r/"+id)

	return
}
