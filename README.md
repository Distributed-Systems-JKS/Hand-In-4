# Hand-In-4
How to start our system



--- Egne noter ---

Update .proto:
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    grpc/proto.proto
