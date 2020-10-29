package server;

import Common.Config;
import Common.Keys;
import com.google.protobuf.ByteString;
import echo.EchoServiceGrpc.EchoServiceImplBase;

import java.io.ByteArrayInputStream;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.Security;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.logging.Logger;

import echo.EmptyServiceGrpc;
import echo.Test;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.NettyServerBuilder;
import io.grpc.stub.StreamObserver;
import io.netty.bootstrap.ServerBootstrap;
import io.netty.channel.EventLoopGroup;
import io.netty.channel.nio.NioEventLoopGroup;
import io.netty.channel.socket.nio.NioServerSocketChannel;
import io.netty.handler.logging.LogLevel;
import io.netty.handler.logging.LoggingHandler;
import io.netty.handler.ssl.*;
import io.netty.util.ResourceLeakDetector;
import io.netty.util.internal.logging.InternalLoggerFactory;
import io.netty.util.internal.logging.Log4JLoggerFactory;

import java.io.IOException;
import java.util.concurrent.TimeUnit;

public class HelloServer {
    static {
//        InternalLoggerFactory.setDefaultFactory(Log4JLoggerFactory.INSTANCE);
    }
    private static final Logger logger = Logger.getLogger(HelloServer.class.getName());
    private static int port = Config.port;
    private Server server;

    private void start() throws IOException {
        SslContext sslCtx = SslContextGMBuilder
                .forServer(Keys.ENC_CERT, Keys.ENC_KEY, Keys.SIGN_CERT, Keys.SIGN_KEY, null)
//                .forServer(selfSignedCertPEM, selfSignedKeyPEM, selfSignedCertPEM, selfSignedKeyPEM, null)
                .applicationProtocolConfig(Config.NPN_AND_ALPN)
                .build();

        server = NettyServerBuilder.forPort(port)
                .addService(new EchoServiceImplBase() {
                    @Override
                    public void echoCall(Test.Echo request, StreamObserver<Test.Echo> responseObserver) {
                        Test.Echo helloResponse = Test.Echo.newBuilder().setPayload(ByteString.copyFrom("pong".getBytes())).build();
                        responseObserver.onNext(helloResponse);
                        responseObserver.onCompleted();
                    }
                })
                .addService(new EmptyServiceGrpc.EmptyServiceImplBase() {
                    @Override
                    public void emptyCall(Test.Empty request, StreamObserver<Test.Empty> responseObserver) {
                        Test.Empty EmptyResponse = Test.Empty.newBuilder().build();
                        responseObserver.onNext(EmptyResponse);
                        responseObserver.onCompleted();
                    }
                })
                .sslContext(sslCtx).build();

        server.start();

        System.out.println("Server started, listening on " + port);
        Runtime.getRuntime().addShutdownHook(new Thread() {
            @Override
            public void run() {
                // Use stderr here since the logger may have been reset by its JVM shutdown hook.
                System.err.println("*** shutting down gRPC server since JVM is shutting down");
                try {
                    HelloServer.this.stop();
                } catch (InterruptedException e) {
                    e.printStackTrace(System.err);
                }
                System.err.println("*** server shut down");
            }
        });
    }

    private void stop() throws InterruptedException {
        if (server != null) {
            server.shutdown().awaitTermination(30, TimeUnit.SECONDS);
        }
    }

    private void blockUntilShutdown() throws InterruptedException {
        if (server != null) {
            server.awaitTermination();
        }
    }

    public static void main(String[] args) throws IOException, InterruptedException {
        final HelloServer server = new HelloServer();
        server.start();
        server.blockUntilShutdown();
    }
}
