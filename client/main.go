package main

import (
	"context"
	x509 "github.com/tjfoc/gmsm/sm2"
	credentials "github.com/tjfoc/gmtls/gmcredentials"
	"google.golang.org/grpc"
	"log"
	echopb "testgm/client/echo"
	"time"
)

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

func invokeEmptyCall(address string, dialOptions []grpc.DialOption) (*echopb.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	//create GRPC client conn
	clientConn, err := grpc.DialContext(ctx, address, dialOptions...)
	if err != nil {
		return nil, err
	}
	defer clientConn.Close()

	//create GRPC client
	client := echopb.NewEmptyServiceClient(clientConn)

	//invoke service
	empty, err := client.EmptyCall(context.Background(), new(echopb.Empty))
	if err != nil {
		return nil, err
	}

	return empty, nil
}

func main() {
	certPool := x509.NewCertPool()

	if !certPool.AppendCertsFromPEM([]byte(selfSignedCertPEM)) {
		log.Fatal("Failed to append certificate to client credentials")
	}

	creds := credentials.NewClientTLSFromCert(certPool, "")

	// GRPC client options
	var dialOptions []grpc.DialOption
	dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))

	address := "127.0.0.1:57373"
	_, err := invokeEmptyCall(address, dialOptions)

	if err != nil {
		log.Fatalf("GRPC client failed to invoke the EmptyCall service on %s: %v",
			address, err)
	} else {
		log.Println("GRPC client successfully invoked the EmptyCall service: " + address)
	}
}
