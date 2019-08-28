
ledger/%.pb.go: ledger/%.proto
	protoc -I ./ledger $< --go_out=plugins=grpc:./ledger/

run-server:
	go run cmd/server/main.go
