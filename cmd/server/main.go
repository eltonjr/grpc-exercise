package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"

	"github.com/eltonjr/protocol-buffer-exercise/ledger"
)

const dbFile = "protoDB.db"

type taskServer struct{}

func (ts taskServer) List(context.Context, *ledger.Void) (*ledger.DebtList, error) {
	return nil, nil
}

func main() {
	srv := grpc.NewServer()
	var tasks taskServer
	ledger.RegisterDebtsServer(srv, tasks)

	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal("could not listen on port :8888")
	}
	log.Fatal(srv.Serve(l))
}

func add(value int64, s string) error {
	debt := &ledger.Debt{
		Desc:  s,
		Value: value,
	}

	b, err := proto.Marshal(debt)
	if err != nil {
		return fmt.Errorf("could not marshal debt: %v", err)
	}

	f, err := os.OpenFile(dbFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("could not open file %s: %v", dbFile, err)
	}

	if err = gob.NewEncoder(f).Encode(int64(len(b))); err != nil {
		return fmt.Errorf("could not write message's length: %v", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("could not write debt to file: %v", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("could not close file %s: %v", dbFile, err)
	}

	return nil
}

func list() error {
	b, err := ioutil.ReadFile(dbFile)
	if err != nil {
		return fmt.Errorf("could not open file %s: %v", dbFile, err)
	}

	for {
		if len(b) == 0 {
			return nil
		}
		if len(b) < 4 {
			return fmt.Errorf("wrong length length %d", len(b))
		}

		var length int64
		if err = gob.NewDecoder(bytes.NewReader(b[:4])).Decode(&length); err != nil {
			return fmt.Errorf("could not decode length: %v", err)
		}
		b = b[4:]

		var debt ledger.Debt
		if err = proto.Unmarshal(b[:length], &debt); err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("could not decode proto debt: %v", err)
		}
		b = b[length:]

		fmt.Printf("Debt: %d | %s\n", debt.Value, debt.Desc)
	}

	return nil
}
