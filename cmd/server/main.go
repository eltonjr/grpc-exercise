package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"

	"github.com/eltonjr/protocol-buffer-exercise/ledger"
)

const (
	dbFile     = "protoDB.db"
	binarySize = 8
)

type debtsServer struct{}

func main() {
	srv := grpc.NewServer()
	var server debtsServer
	ledger.RegisterDebtsServer(srv, server)

	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal("could not listen on port :8888")
	}
	fmt.Println("Listening on :8888")
	log.Fatal(srv.Serve(l))
}

func (ds debtsServer) Add(ctx context.Context, debt *ledger.Debt) (*ledger.Void, error) {
	b, err := proto.Marshal(debt)
	if err != nil {
		return nil, fmt.Errorf("could not marshal debt: %v", err)
	}

	f, err := os.OpenFile(dbFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %v", dbFile, err)
	}

	if err := binary.Write(f, binary.LittleEndian, int64(len(b))); err != nil {
		return nil, fmt.Errorf("could not encode length of message: %v", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return nil, fmt.Errorf("could not write debt to file: %v", err)
	}

	if err = f.Close(); err != nil {
		return nil, fmt.Errorf("could not close file %s: %v", dbFile, err)
	}

	return &ledger.Void{}, nil
}

func (ds debtsServer) List(context.Context, *ledger.Void) (*ledger.DebtList, error) {
	b, err := ioutil.ReadFile(dbFile)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %v", dbFile, err)
	}

	var debts ledger.DebtList
	for {
		if len(b) == 0 {
			return &debts, nil
		}
		if len(b) < binarySize {
			return nil, fmt.Errorf("wrong length length %d", len(b))
		}

		var length int64
		if err := binary.Read(bytes.NewReader(b[:binarySize]), binary.LittleEndian, &length); err != nil {
			return nil, fmt.Errorf("could not decode message length: %v", err)
		}
		b = b[binarySize:]

		var debt ledger.Debt
		if err = proto.Unmarshal(b[:length], &debt); err != nil {
			return nil, fmt.Errorf("could not decode proto debt: %v", err)
		}
		b = b[length:]

		debts.Debts = append(debts.Debts, &debt)
	}

	return &debts, nil
}
