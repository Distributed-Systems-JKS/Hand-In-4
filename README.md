# Hand-In-4
How to start our system

Our program works for three nodes, no more or less



--- Egne noter ---
Run server with port 5051
go run server/server.go -port=5051

Run client with port 5051
go run client/client.go -addr=localhost:5051

Update .proto:
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    grpc/proto.proto
