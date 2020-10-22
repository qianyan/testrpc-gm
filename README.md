Test gRPC with TLS
===


## Building and Running Service

In order to build, Go to ``Go`` module root directory location (server) and execute the following
 shell command,
```
go build -i -v -o bin/server
```

In order to run, Go to ``Go`` module root directory location (server) and execute the following
shell command,

```
./bin/server
```

## Building and Running Client

In order to build, Go to ``Go`` module root directory location (client) and execute the following
 shell command,
```
go build -i -v -o bin/client
```

In order to run, Go to ``Go`` module root directory location (client) and execute the following
shell command,

```
./bin/client
```

## Additional Information

### Generate Server and Client side code
Pre-generated stub file is included in the go project. If you need to generate the stub files please use the below
 command from the root directory
```
protoc --proto_path=echo --go_out=plugins=grpc:echo test.proto
```