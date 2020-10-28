package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net"
	echopb "testgm/server/echo"

	"github.com/Hyperledger-TWGC/tjfoc-gm/gmtls"
	"github.com/Hyperledger-TWGC/tjfoc-gm/gmtls/gmcredentials"
	"github.com/Hyperledger-TWGC/tjfoc-gm/x509"
	"google.golang.org/grpc"
)

var caCert = "testdata/ca.cert"

// 签名证书和加密证书一定要由同一个 ca 签发
// 证书的 KeyUsage 必须要对应上，分别是 Digital Signature 和 Key Encipherment
var signCert = "testdata/sign.cert"
var signKey = "testdata/sign.key"
var encryptCert = "testdata/encrypt.cert"
var encryptKey = "testdata/encrypt.key"

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

const (
	port = ":50051"
)

func main() {
	signCert, err := gmtls.LoadGMX509KeyPair(signCert, signKey)
	if err != nil {
		log.Fatal(err)
	}

	encryptCert, err := gmtls.LoadGMX509KeyPair(encryptCert, encryptKey)
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()

	cacert, err := ioutil.ReadFile(caCert)
	if err != nil {
		log.Fatal(err)
	}
	certPool.AppendCertsFromPEM(cacert)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("fail to listen: %v", err)
	}
	creds := gmcredentials.NewTLS(&gmtls.Config{
		GMSupport:    &gmtls.GMSupport{}, //必须的，国密支持开关
		ClientAuth:   gmtls.RequireAndVerifyClientCert,
		Certificates: []gmtls.Certificate{signCert, encryptCert}, // 证书数组构造时候，签名证书一定要在加密证书前面
		ClientCAs:    certPool,
	})

	srv := grpc.NewServer(grpc.Creds(creds))
	// register the GRPC test server
	echopb.RegisterEmptyServiceServer(srv, &emptyServiceServer{})

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
