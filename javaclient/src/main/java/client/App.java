package client;

import echo.EmptyServiceGrpc;
import echo.Test;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;

import java.util.logging.Logger;

public class App {
    private static final Logger logger = Logger.getLogger(App.class.getName());

    public static void main(String[] args) throws InterruptedException {
        ManagedChannel channel = ManagedChannelBuilder.forAddress("localhost", 50051)
            .useTransportSecurity()
            .usePlaintext()
            .build();

        EmptyServiceGrpc.EmptyServiceBlockingStub stub =
            EmptyServiceGrpc.newBlockingStub(channel);

        Test.Empty empty = stub.emptyCall(Test.Empty.newBuilder().build());
        logger.info("empty call successfully.");

        channel.shutdown();
    }
}
