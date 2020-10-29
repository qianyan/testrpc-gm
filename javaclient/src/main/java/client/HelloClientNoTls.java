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
import io.netty.handler.ssl.*;

import java.util.List;
import java.util.logging.Logger;

public class HelloClientNoTls {
    private static final Logger logger = Logger.getLogger(HelloClientNoTls.class.getName());
    private static int port = Config.portNoTls;

    public static void main(String[] args) throws InterruptedException {
        ManagedChannel channel = null;
        try {
            channel = NettyChannelBuilder
                    .forAddress("127.0.0.1", port)
                    .negotiationType(NegotiationType.PLAINTEXT)
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
