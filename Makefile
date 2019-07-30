
ledger/%.pb.go: ledger/%.proto
	protoc -I ./ledger $< --go_out=plugins=grpc:./ledger/