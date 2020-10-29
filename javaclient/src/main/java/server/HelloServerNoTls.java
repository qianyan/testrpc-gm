package server;

import Common.Config;
import Common.Keys;
import com.google.protobuf.ByteString;
import echo.EchoServiceGrpc.EchoServiceImplBase;
import echo.EmptyServiceGrpc;
import echo.Test;
import io.grpc.Server;
import io.grpc.netty.NettyServerBuilder;
import io.grpc.stub.StreamObserver;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.SslContextGMBuilder;

import java.io.IOException;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;

public class HelloServerNoTls {
    static {
//        InternalLoggerFactory.setDefaultFactory(Log4JLoggerFactory.INSTANCE);
    }

    private static final Logger logger = Logger.getLogger(HelloServerNoTls.class.getName());
    private static int port = Config.portNoTls;
    private Server server;

    private void startNoTls() throws IOException {

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
                .build();

        server.start();

        System.out.println("Server started, listening on " + port);
        Runtime.getRuntime().addShutdownHook(new Thread() {
            @Override
            public void run() {
                // Use stderr here since the logger may have been reset by its JVM shutdown hook.
                System.err.println("*** shutting down gRPC server since JVM is shutting down");
                try {
                    HelloServerNoTls.this.stop();
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
        final HelloServerNoTls server = new HelloServerNoTls();
        server.startNoTls();
        server.blockUntilShutdown();
    }
}
