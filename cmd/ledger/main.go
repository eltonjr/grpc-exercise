package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	grpc "google.golang.org/grpc"

	"github.com/eltonjr/protocol-buffer-exercise/ledger"
)

const dbFile = "protoDB.db"

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Usage: ledger [add|list]")
		os.Exit(0)
	}

	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not dial to grpc server at 8888, %v", err)
	}
	client := ledger.NewDebtsClient(conn)

	switch cmd := flag.Args()[0]; cmd {
	case "add":
		if flag.NArg() < 3 {
			err = fmt.Errorf("add command needs a value and a description")
		} else {
			value, err := strconv.ParseInt(flag.Args()[1], 10, 64)
			if err != nil {
				err = fmt.Errorf("value being added must be a number: %v", err)
			} else {
				err = add(context.Background(), client, value, strings.Join(flag.Args()[2:], " "))
			}
		}
	case "list":
		err = list(context.Background(), client)
	default:
		err = fmt.Errorf("invalid command '%s'", cmd)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func add(ctx context.Context, client ledger.DebtsClient, value int64, s string) error {
	debt := &ledger.Debt{
		Desc:  s,
		Value: value,
	}

	_, err := client.Add(ctx, debt)
	if err != nil {
		return fmt.Errorf("could not add debt in server: %v", err)
	}

	fmt.Println("debt added successfully")
	return nil
}

func list(ctx context.Context, client ledger.DebtsClient) error {
	debts, err := client.List(ctx, &ledger.Void{})
	if err != nil {
		return fmt.Errorf("could not list from grpc server: %v", err)
	}

	for _, debt := range debts.Debts {
		fmt.Printf("Debt: %d | %s\n", debt.Value, debt.Desc)
	}
	return nil
}
