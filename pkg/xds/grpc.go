//go:generate protoc --gogofaster_out=plugins=grpc:$GOPATH/src xds.proto
package xds

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"log"
	"math/big"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// To install the tools:
// go get -u github.com/gogo/protobuf/protoc-gen-gogofaster
// go get -u github.com/gogo/protobuf/gogoproto

// alt: --go_out=plugins=grpc:.

type GrpcService struct {
}

// Subscribe maps the the webpush subscribe request
func (s *GrpcService) StreamAggregatedResources(AggregatedDiscoveryService_StreamAggregatedResourcesServer) error {
	return nil
}


func Connect(addr string, clientPem string) (*grpc.ClientConn, AggregatedDiscoveryServiceClient, error) {
	opts := []grpc.DialOption{}

	// Cert file is a PEM, it is loaded into a x509.CertPool,
	// will be set as RootCAs in the tls.Config
	//creds, err := credentials.NewClientTLSFromFile("../testdata/ca.pem", "x.test.youtube.com")

	cp := x509.NewCertPool()
	cp.AppendCertsFromPEM([]byte(clientPem))

	creds := credentials.NewTLS(&tls.Config{
		// name will override ":authority" request headers,
		// also used to validate the cert ( it seems to make it to the
		// SNI header as well ?). Must match one of the cert entries
		ServerName: "x.test.youtube.com",
		// pub keys that signed the cert
		RootCAs: cp})

	if false {
		// example code: docker/libtrust/certificates.go
		cert := &x509.Certificate{
			SerialNumber: big.NewInt(0),
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			NotBefore:    time.Now().Add(-time.Hour * 24 * 7),
			NotAfter:     time.Now().Add(time.Hour * 24 * 365 * 10),
		}
		issCert := &x509.Certificate{
			SerialNumber: big.NewInt(0),
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			NotBefore:    time.Now().Add(-time.Hour * 24 * 7),
			NotAfter:     time.Now().Add(time.Hour * 24 * 365 * 10),
		}
		pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			fmt.Println(err)
		}

		certDER, err := x509.CreateCertificate(rand.Reader,
			cert, issCert, pk.Public(), pk)
		if err != nil {
			fmt.Println(err)
		}

		cp := x509.NewCertPool()
		c, err := x509.ParseCertificate(certDER)
		if err != nil {
			fmt.Println(err)
		}
		cp.AddCert(c)

		creds = credentials.NewTLS(&tls.Config{
			RootCAs: cp,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	opts = append(opts, grpc.WithTransportCredentials(creds))

	// WithInsecure() == no auth, plain http. Either that or TransportCred
	// required.

	// grpc/credentials/oauth defines a number of options - it's
	// an interface called on each request, returning headers
	//
	// GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
	// TODO: vapid option

	// json has client_email, private_key.
	// Aud: "https://accounts.google.com/o/oauth2/token"
	// Scope: ...
	// Iss: Email
	// Signed and posted to google, which returns an oauth2 token

	// Uses 2-legged auth
	// RS256/JWT, with Iss = email, Sub = email, Aud = (google url)
	// Iat, Exp = 1h
	//json := []byte("{}")
	//jwtCreds, err := oauth.NewJWTAccessFromKey(json)

	// Other option: NewServiceAccountFromFile ( dev console service account
	// -- application default credentials )

	// golang.org/x/oauth2/google/DefaultTokenSource(ctx, scope...)
	// GOOGLE_APPLICATION_CREDENTIALS env file
	// .config/gcloud/application_default_credentials.json
	// appengine or compute engine get it from env

	// file has client_id, client_secret, refresh_token - for creds
	// and private key/key id for service accounts
	// Still goes trough Oauth2 flow.
	//gcreds, err := oauth.NewApplicationDefault(context.Background(), "test_scope")
	//gcreds.GetRequestMetadata(context.Background(), "url")

	// also NewOauthAccess - seems to allow arbitrary type/value
	// could be IID token !!!!

	//opts = append(opts, grpc.WithPerRPCCredentials(jwtCreds))

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Println("Failed to dial", err)
		return nil, nil, err
	}

	return conn, NewAggregatedDiscoveryServiceClient(conn), nil
}