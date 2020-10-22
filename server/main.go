package main

import (
	"context"
	"io"
	"log"
	"net"
	"testgm/server/comm"
	echopb "testgm/server/echo"
	"time"
)

var selfSignedKeyPEM = `-----BEGIN EC PRIVATE KEY-----
MIGTAgEAMBMGByqGSM49AgEGCCqBHM9VAYItBHkwdwIBAQQgqlBIAUaTj6uG+ETN
WVSm5x75hLNSVUf5PnJbixQBHGCgCgYIKoEcz1UBgi2hRANCAAS2CgsRr8CP/Erj
eBiJx9ppfbAfZbIQI9dHUm0AQsbVWlO6jNDgxTi47Wmf5gtilYeUqIBScI/BaWkQ
An+1jIwh
-----END EC PRIVATE KEY-----
`
var selfSignedCertPEM = `-----BEGIN CERTIFICATE-----
MIICFDCCAbugAwIBAgIQTH+Jw6wgrqvFn8nN2Z4iNjAKBggqgRzPVQGDdTBcMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEPMA0GA1UEChMGc2VydmVyMQ8wDQYDVQQDEwZzZXJ2ZXIwHhcNMjAx
MDEzMTQ1NDQ0WhcNMzAxMDExMTQ1NDQ0WjBcMQswCQYDVQQGEwJVUzETMBEGA1UE
CBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEPMA0GA1UEChMG
c2VydmVyMQ8wDQYDVQQDEwZzZXJ2ZXIwWTATBgcqhkjOPQIBBggqgRzPVQGCLQNC
AAS2CgsRr8CP/ErjeBiJx9ppfbAfZbIQI9dHUm0AQsbVWlO6jNDgxTi47Wmf5gti
lYeUqIBScI/BaWkQAn+1jIwho18wXTAOBgNVHQ8BAf8EBAMCAaYwDwYDVR0lBAgw
BgYEVR0lADAPBgNVHRMBAf8EBTADAQH/MA0GA1UdDgQGBAQBAgMEMBoGA1UdEQQT
MBGCCWxvY2FsaG9zdIcEfwAAATAKBggqgRzPVQGDdQNHADBEAiBqCgFi2yXg0a9y
DvcAZzzLBLve48PAjZfYTi24YA6ovAIgfDXO5BIASJE/aY/0Mkdg6YabI7RJhEcX
/4Mt25/Fsmc=
-----END CERTIFICATE-----
`

type emptyServiceServer struct{}

func (ess *emptyServiceServer) EmptyCall(context.Context, *echopb.Empty) (*echopb.Empty, error) {
	return new(echopb.Empty), nil
}

func (esss *emptyServiceServer) EmptyStream(stream echopb.EmptyService_EmptyStreamServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&echopb.Empty{}); err != nil {
			return err
		}

	}
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	log.Println("Listener address: " + lis.Addr().String())
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}

	srv, err := comm.NewGRPCServerFromListener(lis, comm.ServerConfig{
		ConnectionTimeout: 250 * time.Millisecond,
		SecOpts: &comm.SecureOptions{
			UseTLS:      true,
			Certificate: []byte(selfSignedCertPEM),
			Key:         []byte(selfSignedKeyPEM)}})
	// check for error
	if err != nil {
		log.Fatalf("Failed to return new GRPC server: %v", err)
	}

	// register the GRPC test server
	echopb.RegisterEmptyServiceServer(srv.Server(), &emptyServiceServer{})

	if err := srv.Start(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
