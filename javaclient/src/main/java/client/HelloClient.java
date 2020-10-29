package client;

import Common.Config;
import Common.Keys;
import com.google.protobuf.ByteString;
import echo.EchoServiceGrpc;
import echo.EmptyServiceGrpc;
import echo.Test;
import io.grpc.ManagedChannel;
import io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.NegotiationType;
import io.grpc.netty.NettyChannelBuilder;
import io.netty.handler.codec.http2.Http2SecurityUtil;
import io.netty.handler.ssl.*;

import javax.net.ssl.SSLException;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.logging.Logger;

import static Common.Keys.ENC_KEY;

public class HelloClient {
    private static final Logger logger = Logger.getLogger(HelloClient.class.getName());
    private static int port = Config.port;
    private static final String pem = Keys.goPem;
    private static final ApplicationProtocolConfig NPN_AND_ALPN = Config.NPN_AND_ALPN;
    private static List<String> ciphers = Config.CIPHERS;

    public static void main(String[] args) throws InterruptedException {
        ManagedChannel channel = null;
        try {
            SslContextBuilder clientContextBuilder = GrpcSslContexts.configure(SslContextBuilder.forClient(), SslProvider.OPENSSL);

            SslContextGMBuilder gmContextBuilder = SslContextGMBuilder
                    .forClient()
                    .clientAuth(ClientAuth.REQUIRE)
//                    .protocols(Config.TLS_PROTOCOL)
                    .ciphers(ciphers, SupportedCipherSuiteFilter.INSTANCE)
                    .applicationProtocolConfig(NPN_AND_ALPN)
                    .trustManager(Keys.TRUST_CERT)
                    .keyManager(Keys.ENC_CERT, Keys.ENC_KEY, Keys.SIGN_CERT, Keys.SIGN_KEY, null);

            SslContext sslContext = gmContextBuilder.build();
            channel = NettyChannelBuilder
                    .forAddress("127.0.0.1", port)
                    .sslContext(sslContext)
                    .negotiationType(NegotiationType.TLS)
                    .build();

            {
                EmptyServiceGrpc.EmptyServiceBlockingStub stub = EmptyServiceGrpc.newBlockingStub(channel);

                Test.Empty empty = stub.emptyCall(Test.Empty.newBuilder().build());
                logger.info("empty call successfully.");
            }
            {
                EchoServiceGrpc.EchoServiceBlockingStub stub = EchoServiceGrpc.newBlockingStub(channel);
                Test.Echo echo = stub.echoCall(Test.Echo.newBuilder().setPayload(ByteString.copyFrom("ping".getBytes())).build());
                logger.info("echo call successfully. " + echo.getPayload().toString("UTF-8"));
            }

        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            if (channel != null)
                channel.shutdown();
        }
    }
}
