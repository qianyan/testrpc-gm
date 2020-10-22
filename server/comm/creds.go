package comm

import (
	"context"
	"errors"
	"github.com/tjfoc/gmsm/sm2"
	tls "github.com/tjfoc/gmtls"
	"github.com/tjfoc/gmtls/gmcredentials"
	"google.golang.org/grpc/credentials"
	"net"
	"sync"
)

var (
	ClientHandshakeNotImplError = errors.New("core/comm: Client handshakes" +
		"are not implemented with serverCreds")
	OverrrideHostnameNotSupportedError = errors.New(
		"core/comm: OverrideServerName is " +
			"not supported")
	ServerHandshakeNotImplementedError = errors.New("core/comm: server handshakes are not implemented with clientCreds")

	MissingServerConfigError = errors.New(
		"core/comm: `serverConfig` cannot be nil")
	// alpnProtoStr are the specified application level protocols for gRPC.
	alpnProtoStr = []string{"h2"}
)

type TLSConfig struct {
	config *tls.Config
	lock   sync.RWMutex
}

func NewTLSConfig(config *tls.Config) *TLSConfig {
	return &TLSConfig{
		config: config,
	}
}

func (t *TLSConfig) Config() tls.Config {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if t.config != nil {
		return *t.config.Clone()
	}

	return tls.Config{}
}

func (t *TLSConfig) AddClientRootCA(cert *sm2.Certificate) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.config.ClientCAs.AddCert(cert)
}

func NewServerTransportCredentials(
	serverConfig *TLSConfig) credentials.TransportCredentials {

	// NOTE: unlike the default grpc/credentials implementation, we do not
	// clone the tls.Config which allows us to update it dynamically
	serverConfig.config.NextProtos = alpnProtoStr
	// override TLS version and ensure it is 1.2
	serverConfig.config.MinVersion = tls.VersionTLS12
	serverConfig.config.MaxVersion = tls.VersionTLS12
	return &serverCreds{serverConfig: serverConfig}
}

// serverCreds is an implementation of grpc/credentials.TransportCredentials.
type serverCreds struct {
	serverConfig *TLSConfig
}

// ClientHandShake is not implemented for `serverCreds`.
func (sc *serverCreds) ClientHandshake(context.Context,
	string, net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return nil, nil, ClientHandshakeNotImplError
}

// ServerHandshake does the authentication handshake for servers.
func (sc *serverCreds) ServerHandshake(rawConn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	serverConfig := sc.serverConfig.Config()

	conn := tls.Server(rawConn, &serverConfig)
	if err := conn.Handshake(); err != nil {
		return nil, nil, err
	}
	return conn, gmcredentials.TLSInfo{State: conn.ConnectionState()}, nil
}

// Info provides the ProtocolInfo of this TransportCredentials.
func (sc *serverCreds) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{
		SecurityProtocol: "tls",
		SecurityVersion:  "1.2",
	}
}

// Clone makes a copy of this TransportCredentials.
func (sc *serverCreds) Clone() credentials.TransportCredentials {
	config := sc.serverConfig.Config()
	serverConfig := NewTLSConfig(&config)
	creds := NewServerTransportCredentials(serverConfig)
	return creds
}

// OverrideServerName overrides the server name used to verify the hostname
// on the returned certificates from the server.
func (sc *serverCreds) OverrideServerName(string) error {
	return OverrrideHostnameNotSupportedError
}
