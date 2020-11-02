package main

import (
	"github.com/Hyperledger-TWGC/tjfoc-gm/gmtls"
	"github.com/Hyperledger-TWGC/tjfoc-gm/x509"
	"github.com/tw-bc-group/net-go-gm/http"
	"io/ioutil"
	"log"
)

var caCert = "testdata/ca.cert"

var clientCert = "testdata/client.cert"
var clientKey = "testdata/client.key"

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

	transport := new(http.Transport)
	certPool.AppendCertsFromPEM(caCert)
	transport.TLSClientConfig = &gmtls.Config{
		GMSupport:    &gmtls.GMSupport{},
		ServerName:   "test.example.com",
		RootCAs:      certPool,
		Certificates: []gmtls.Certificate{clientCert},
		ClientAuth:   gmtls.RequireAndVerifyClientCert,
	}

	httpClient := http.Client{Transport: transport}
	response, err := httpClient.Get("https://localhost:7054/")
	if err != nil {
		log.Fatalf("Failed to get: %v", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read body: %v", err)
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatalf("Failed to close the response body: %s", err.Error())
		}
	}()

	if string(body) != "hello" {
		log.Fatalf("Got hello failed.")
	} else {
		log.Println("success")
	}
}
