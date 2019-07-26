package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"

	"github.com/eltonjr/protocol-buffer-exercise/ledger"
)

const dbFile = "protoDB.db"

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Usage: ledger [add|list]")
		os.Exit(0)
	}

	var err error

	switch cmd := flag.Args()[0]; cmd {
	case "add":
		if flag.NArg() < 3 {
			err = fmt.Errorf("add command needs a value and a description")
		} else {
			value, err := strconv.ParseInt(flag.Args()[1], 10, 64)
			if err != nil {
				err = fmt.Errorf("value being added must be a number: %v", err)
			} else {
				err = add(value, strings.Join(flag.Args()[2:], " "))
			}
		}
	case "list":
		err = list()
	default:
		err = fmt.Errorf("invalid command '%s'", cmd)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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
