## Just a repo for me to play around with protocol-buffers, grpc and go

It was initially based on Francesc Campoy's Just For Func series, ep 30 and 31.

There are two components here:

- **cmd/server** is kind of a webserver that `list` and `add` debts on a ledger.
- **cmd/ledger** is a command-line tool to communicate with the server and trigger the operations.

Every communication between components is by grpc over tcp.  
The server does not have any database, it just stores the Debts in a file using protobuf.  
Both sides share the same `.proto` files from inside `ledger/`.

### Running

#### Server

    go run cmd/server/main.go

#### Client

As a command-line tool, you can just:

    go run cmd/ledger/main.go add 12 lunch
    go run cmd/ledger/main.go list
