package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	flag "github.com/spf13/pflag"
	"github.com/tblyler/sheepmq/sheepmq"
	"github.com/tblyler/sheepmq/shepard"
)

func main() {
	var port uint16
	var listenAddr string

	flag.StringVarP(&listenAddr, "addr", "h", "", "The address to listen on")
	flag.Uint16VarP(&port, "port", "p", 0, "The port to bind")

	flag.Parse()

	listenAddr = fmt.Sprintf("%s:%d", listenAddr, port)

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("Failed to listen on %s: %s", listenAddr, err)
		return
	}

	defer listener.Close()

	fmt.Println("Listening on", listener.Addr())

	grpcServer := grpc.NewServer()

	sheepmqServer, err := sheepmq.NewServer()
	if err != nil {
		fmt.Println("Failed to open sheepmq server:", err)
		return
	}

	sheepmqGServer := sheepmq.NewGServer(sheepmqServer)
	if err != nil {
		fmt.Println("Failed to create sheepmq server:", err)
		return
	}

	shepard.RegisterLoqServer(grpcServer, sheepmqServer)

	grpcServer.Serve(listener)
}
