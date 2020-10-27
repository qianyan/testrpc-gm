package comm

import (
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Hyperledger-TWGC/tjfoc-gm/x509"
	tls "github.com/Hyperledger-TWGC/tjfoc-gm/gmtls"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

)

type GRPCServer struct {
	// Listen address for the server specified as hostname:port
	address string
	// Listener for handling network requests
	listener net.Listener
	// GRPC server
	server *grpc.Server
	// Certificate presented by the server for TLS communication
	// stored as an atomic reference
	serverCertificate atomic.Value
	// Key used by the server for TLS communication
	serverKeyPEM []byte
	// lock to protect concurrent access to append / remove
	lock *sync.Mutex
	// Set of PEM-encoded X509 certificate authorities used to populate
	// the tlsConfig.ClientCAs indexed by subject
	//clientRootCAs map[string]*x509.Certificate
	clientRootCAs map[string]*x509.Certificate
	// TLS configuration used by the grpc server
	tls *TLSConfig
}

func (gServer *GRPCServer) Server() *grpc.Server {
	return gServer.server
}

func (gServer *GRPCServer) Start() error {
	return gServer.server.Serve(gServer.listener)
}

func NewGRPCServerFromListener(listener net.Listener, serverConfig ServerConfig, certs []tls.Certificate) (*GRPCServer, error) {
	grpcServer := &GRPCServer{
		address:  listener.Addr().String(),
		listener: listener,
		lock:     &sync.Mutex{},
	}

	//set up our server options
	var serverOpts []grpc.ServerOption

	//check SecOpts
	var secureConfig SecureOptions
	if serverConfig.SecOpts != nil {
		secureConfig = *serverConfig.SecOpts
	}
	if secureConfig.UseTLS {
		//both key and cert are required
		if secureConfig.Key != nil && secureConfig.Certificate != nil {
			//load server public and private keys
			cert, err := tls.X509KeyPair(secureConfig.Certificate, secureConfig.Key)
			if err != nil {
				return nil, err
			}
			grpcServer.serverCertificate.Store(cert)

			//set up our TLS config
			if len(secureConfig.CipherSuites) == 0 {
				secureConfig.CipherSuites = DefaultTLSCipherSuites
			}
			getCert := func(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
				cert := grpcServer.serverCertificate.Load().(tls.Certificate)
				return &cert, nil
			}
			//base server certificate
			grpcServer.tls = NewTLSConfig(&tls.Config{
				GMSupport:              &tls.GMSupport{},
				VerifyPeerCertificate:  secureConfig.VerifyCertificate,
				GetCertificate:         getCert,
				Certificates:           certs,
				SessionTicketsDisabled: true,
			})

			if serverConfig.SecOpts.TimeShift > 0 {
				timeShift := serverConfig.SecOpts.TimeShift
				grpcServer.tls.config.Time = func() time.Time {
					return time.Now().Add((-1) * timeShift)
				}
			}
			grpcServer.tls.config.ClientAuth = tls.RequestClientCert
			//check if client authentication is required
			if secureConfig.RequireClientCert {
				//require TLS client auth
				grpcServer.tls.config.ClientAuth = tls.RequireAndVerifyClientCert
				//if we have client root CAs, create a certPool
				if len(secureConfig.ClientRootCAs) > 0 {
					grpcServer.clientRootCAs = make(map[string]*x509.Certificate)
					grpcServer.tls.config.ClientCAs = x509.NewCertPool()
					for _, clientRootCA := range secureConfig.ClientRootCAs {
						err = grpcServer.appendClientRootCA(clientRootCA)
						if err != nil {
							return nil, err
						}
					}
				}
			}

			// create credentials and add to server options
			creds := NewServerTransportCredentials(grpcServer.tls)
			serverOpts = append(serverOpts, grpc.Creds(creds))
		} else {
			return nil, errors.New("serverConfig.SecOpts must contain both Key and Certificate when UseTLS is true")
		}
	}
	// set max send and recv msg sizes
	serverOpts = append(serverOpts, grpc.MaxSendMsgSize(MaxSendMsgSize))
	serverOpts = append(serverOpts, grpc.MaxRecvMsgSize(MaxRecvMsgSize))
	// set the keepalive options
	serverOpts = append(serverOpts, ServerKeepaliveOptions(serverConfig.KaOpts)...)
	// set connection timeout
	if serverConfig.ConnectionTimeout <= 0 {
		serverConfig.ConnectionTimeout = DefaultConnectionTimeout
	}
	serverOpts = append(
		serverOpts,
		grpc.ConnectionTimeout(serverConfig.ConnectionTimeout))
	// set the interceptors
	if len(serverConfig.StreamInterceptors) > 0 {
		serverOpts = append(
			serverOpts,
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(serverConfig.StreamInterceptors...)),
		)
	}
	if len(serverConfig.UnaryInterceptors) > 0 {
		serverOpts = append(
			serverOpts,
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(serverConfig.UnaryInterceptors...)),
		)
	}

	grpcServer.server = grpc.NewServer(serverOpts...)

	return grpcServer, nil
}

// internal function to add a PEM-encoded clientRootCA
func (gServer *GRPCServer) appendClientRootCA(clientRoot []byte) error {

	errMsg := "Failed to append client root certificate(s): %s"
	//convert to x509
	certs, subjects, err := pemToX509Certs(clientRoot)
	if err != nil {
		return fmt.Errorf(errMsg, err.Error())
	}

	if len(certs) < 1 {
		return fmt.Errorf(errMsg, "No client root certificates found")
	}

	for i, cert := range certs {
		//first add to the ClientCAs
		gServer.tls.AddClientRootCA(cert)
		//add it to our clientRootCAs map using subject as key
		gServer.clientRootCAs[subjects[i]] = cert
	}
	return nil
}

//utility function to parse PEM-encoded certs
//func pemToX509Certs(pemCerts []byte) ([]*x509.Certificate, []string, error) {
func pemToX509Certs(pemCerts []byte) ([]*x509.Certificate, []string, error) {
	//it's possible that multiple certs are encoded
	//certs := []*x509.Certificate{}
	certs := []*x509.Certificate{}
	subjects := []string{}
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		/** TODO: check why msp does not add type to PEM header
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}
		*/

		//cert, err := x509.ParseCertificate(block.Bytes)
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, subjects, err
		} else {
			certs = append(certs, cert)
			//extract and append the subject
			subjects = append(subjects, string(cert.RawSubject))
		}
	}
	return certs, subjects, nil
}
