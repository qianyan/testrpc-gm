package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/Hyperledger-TWGC/tjfoc-gm/gmtls"
	"github.com/Hyperledger-TWGC/tjfoc-gm/gmtls/gmcredentials"
	"github.com/Hyperledger-TWGC/tjfoc-gm/x509"
	"google.golang.org/grpc"
	echopb "testgm/client/echo"
)

var caCert = "testdata/ca.cert"
var clientCert = "testdata/client.cert"
var clientKey = "testdata/client.key"

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

const (
	address = "localhost:50051"
)

func main() {
	clientCert, err := gmtls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(caCert)
	if err != nil {
		log.Fatal(err)
	}

	certPool.AppendCertsFromPEM(caCert)
	creds := gmcredentials.NewTLS(&gmtls.Config{
		GMSupport:    &gmtls.GMSupport{},
		ServerName:   "test.example.com",
		Certificates: []gmtls.Certificate{clientCert},
		RootCAs:      certPool,
		ClientAuth:   gmtls.RequireAndVerifyClientCert,
	})

	// GRPC client options
	var dialOptions []grpc.DialOption
	dialOptions = append(dialOptions, grpc.WithTransportCredentials(creds))

	_, err = invokeEmptyCall(address, dialOptions)

	if err != nil {
		log.Fatalf("GRPC client failed to invoke the EmptyCall service on %s: %v",
			address, err)
	} else {
		log.Println("GRPC client successfully invoked the EmptyCall service: " + address)
	}
}
