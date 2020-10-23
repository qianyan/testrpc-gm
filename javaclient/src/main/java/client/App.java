package client;

import echo.EmptyServiceGrpc;
import echo.Test;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.netty.NegotiationType;
import io.grpc.netty.NettyChannelBuilder;
import io.netty.handler.ssl.*;

import javax.net.ssl.SSLException;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.logging.Logger;

import static java.lang.String.format;

public class App {
    private static final Logger logger = Logger.getLogger(App.class.getName());
    private static final String TLS_PROTOCOL = "TLSv1.2";
    private static final String pem = "-----BEGIN CERTIFICATE-----\n" +
            "MIICFDCCAbugAwIBAgIQTH+Jw6wgrqvFn8nN2Z4iNjAKBggqgRzPVQGDdTBcMQsw\n" +
            "CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy\n" +
            "YW5jaXNjbzEPMA0GA1UEChMGc2VydmVyMQ8wDQYDVQQDEwZzZXJ2ZXIwHhcNMjAx\n" +
            "MDEzMTQ1NDQ0WhcNMzAxMDExMTQ1NDQ0WjBcMQswCQYDVQQGEwJVUzETMBEGA1UE\n" +
            "CBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZyYW5jaXNjbzEPMA0GA1UEChMG\n" +
            "c2VydmVyMQ8wDQYDVQQDEwZzZXJ2ZXIwWTATBgcqhkjOPQIBBggqgRzPVQGCLQNC\n" +
            "AAS2CgsRr8CP/ErjeBiJx9ppfbAfZbIQI9dHUm0AQsbVWlO6jNDgxTi47Wmf5gti\n" +
            "lYeUqIBScI/BaWkQAn+1jIwho18wXTAOBgNVHQ8BAf8EBAMCAaYwDwYDVR0lBAgw\n" +
            "BgYEVR0lADAPBgNVHRMBAf8EBTADAQH/MA0GA1UdDgQGBAQBAgMEMBoGA1UdEQQT\n" +
            "MBGCCWxvY2FsaG9zdIcEfwAAATAKBggqgRzPVQGDdQNHADBEAiBqCgFi2yXg0a9y\n" +
            "DvcAZzzLBLve48PAjZfYTi24YA6ovAIgfDXO5BIASJE/aY/0Mkdg6YabI7RJhEcX\n" +
            "/4Mt25/Fsmc=\n" +
            "-----END CERTIFICATE-----";
    private static final String GRPC_EXP_VERSION = "grpc-exp";
    private static final String HTTP2_VERSION = "h2";
    private static final List<String> NEXT_PROTOCOL_VERSIONS =
            Collections.unmodifiableList(Arrays.asList(GRPC_EXP_VERSION, HTTP2_VERSION));

    private static final ApplicationProtocolConfig NPN_AND_ALPN = new ApplicationProtocolConfig(
            ApplicationProtocolConfig.Protocol.NPN_AND_ALPN,
            ApplicationProtocolConfig.SelectorFailureBehavior.NO_ADVERTISE,
            ApplicationProtocolConfig.SelectedListenerFailureBehavior.ACCEPT,
            NEXT_PROTOCOL_VERSIONS);

    private static List<String> ciphers = Collections.unmodifiableList(Arrays
            .asList(
                    "TLS_ECDHE_WITH_SM4_SM3",
                    "TLS_ECC_WITH_SM4_SM3"
            ));

    public static void main(String[] args) throws InterruptedException {
        ManagedChannel channel = null;
        try {

            SslProvider sslprovider = SslProvider.OPENSSL;
            NegotiationType ntype = NegotiationType.TLS;


            SslContext sslContext = SslContextGMBuilder.forClient()
                    .protocols(TLS_PROTOCOL)
                    .ciphers(ciphers, SupportedCipherSuiteFilter.INSTANCE)
                    .applicationProtocolConfig(NPN_AND_ALPN)
                    .trustManager(pem)
//                    .keyManager(ENC_CERT, ENC_KEY, SIGN_CERT, SIGN_KEY, null)
                    .build();

            NettyChannelBuilder channelBuilder = NettyChannelBuilder
                    .forAddress("localhost", 50051)
                    .sslContext(sslContext)
                    .negotiationType(ntype);


            channel = channelBuilder.build();

            EmptyServiceGrpc.EmptyServiceBlockingStub stub =
                    EmptyServiceGrpc.newBlockingStub(channel);

            Test.Empty empty = stub.emptyCall(Test.Empty.newBuilder().build());
            logger.info("empty call successfully.");

        } catch (SSLException e) {
            e.printStackTrace();
        } finally {
            if (channel != null)
                channel.shutdown();
        }
    }
}
