package cloud

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var Host = "cloud.evcc.io:8080"

var (
	connections = make(map[string]*grpc.ClientConn)
)

//go:embed ca-cert.pem
var caCert []byte

func caPEM() []byte {
	copy := bytes.NewBuffer(caCert)
	return copy.Bytes()
}

func loadTLSCredentials() (*tls.Config, error) {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPEM()) {
		return nil, fmt.Errorf("failed to add CA certificate")
	}

	// create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return config, nil
}

func Connection(uri string) (*grpc.ClientConn, error) {
	var err error

	conn, ok := connections[uri]
	if !ok {
		transportOption := grpc.WithInsecure()
		if !strings.HasPrefix(uri, "localhost") {
			var tlsConfig *tls.Config
			if tlsConfig, err = loadTLSCredentials(); err != nil {
				return nil, err
			}

			transportOption = grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
		}

		if conn, err = grpc.Dial(uri, transportOption); err == nil {
			connections[uri] = conn
		}
	}

	return conn, err
}
