package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/Hyperledger-TWGC/tjfoc-gm/gmtls"
	"github.com/Hyperledger-TWGC/tjfoc-gm/x509"
	"github.com/gorilla/mux"
)

var caCert = "testdata/ca.cert"

// 签名证书和加密证书一定要由同一个 ca 签发
// 证书的 KeyUsage 必须要对应上，分别是 Digital Signature 和 Key Encipherment
var signCert = "testdata/sign.cert"
var signKey = "testdata/sign.key"
var encryptCert = "testdata/encrypt.cert"
var encryptKey = "testdata/encrypt.key"

func main() {
	certPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(caCert)
	if err != nil {
		log.Fatal(err)
	}
	certPool.AppendCertsFromPEM(caCert)

	signCert, err := gmtls.LoadGMX509KeyPair(signCert, signKey)
	if err != nil {
		log.Fatal(err)
	}

	encryptCert, err := gmtls.LoadGMX509KeyPair(encryptCert, encryptKey)
	if err != nil {
		log.Fatal(err)
	}

	addr := net.JoinHostPort("127.0.0.1", "7054")

	tlsConfig := &gmtls.Config{
		//必须的，国密支持开关
		GMSupport: &gmtls.GMSupport{},
		// 证书数组构造时候，签名证书一定要在加密证书前面
		Certificates: []gmtls.Certificate{signCert, encryptCert},
		ClientCAs:    certPool,
	}

	listener, err := gmtls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}

	sMux := mux.NewRouter()
	sMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err:= w.Write([]byte("hello"))
		if err != nil {
			log.Fatalf("Failed to write response: %v", err)
		}
	})

	err = http.Serve(listener, sMux)
	err = listener.Close()
	if err != nil {
		log.Fatalf("Stop: failed to close listener: %v", err)
	}

}
