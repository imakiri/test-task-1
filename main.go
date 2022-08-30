package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/imakiri/test-task-1/internal"
	proto "github.com/imakiri/test-task-1/internal/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

var (
	portGRPC      = flag.String("port-grpc", "25565", "gRPC server port")
	portREST      = flag.String("port-rest", "8080", "REST server endpoint")
	portSwaggerUI = flag.String("port-swagger_ui", "80", "SwaggerUI server endpoint")
)

func NewGRPCService() (*grpc.Server, error) {
	var server = grpc.NewServer()
	var service, err = internal.NewService()
	if err != nil {
		log.Fatalln(err)
	}

	proto.RegisterServiceServer(server, service)
	return server, nil
}

func MustLaunchGRPC() <-chan error {
	var ch = make(chan error)
	var service, err1 = NewGRPCService()
	if err1 != nil {
		log.Fatalln(err1)
	}

	var lis, err2 = net.Listen("tcp", ":"+*portGRPC)
	if err2 != nil {
		log.Fatalln(err1)
	}

	go func() {
		fmt.Printf("Launching gRPC service at %s\n", lis.Addr().String())
		ch <- service.Serve(lis)
	}()

	return ch
}

func NewRESTService(ui bool) (*http.Server, error) {
	var mux = runtime.NewServeMux()
	var opts = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	var err = proto.RegisterServiceHandlerFromEndpoint(context.Background(), mux, ":"+*portGRPC, opts)
	if err != nil {
		return nil, err
	}

	var server = new(http.Server)
	if ui {
		var _mux = new(http.ServeMux)
		_mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./dist"))))
		_mux.Handle("/v1/", mux)
		server.Handler = _mux
	} else {
		server.Handler = mux
	}

	return server, nil
}

func MustSwaggerUI(lis net.Listener) <-chan error {
	var ch = make(chan error)
	var server = new(http.Server)
	var mux = new(http.ServeMux)
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./dist"))))
	server.Handler = mux

	if lis == nil {
		var err error
		lis, err = net.Listen("tcp", ":"+*portSwaggerUI)
		if err != nil {
			log.Fatalln(err)
		}
	}

	go func() {
		fmt.Printf("Launching SwaggerUI service at %s\n", lis.Addr().String())
		ch <- server.Serve(lis)
	}()

	return ch
}

func MustLaunchREST() <-chan error {
	var ch = make(chan error)
	var ui bool
	if *portSwaggerUI == *portREST {
		ui = true
	} else {
		ui = false
	}

	var service, err1 = NewRESTService(ui)
	if err1 != nil {
		log.Fatalln(err1)
	}

	var lis, err2 = net.Listen("tcp", ":"+*portREST)
	if err2 != nil {
		log.Fatalln(err1)
	}

	go func() {
		fmt.Printf("Launching REST service at %s\n", lis.Addr().String())
		ch <- service.Serve(lis)
	}()

	if ui {
		var serviceUI = MustSwaggerUI(lis)
		go func() {
			select {
			case err := <-serviceUI:
				ch <- err
			}
		}()
	}

	return ch
}

func main() {
	flag.Parse()
	var chGRPC = MustLaunchGRPC()
	var chREST = MustLaunchREST()
	select {
	case err := <-chREST:
		log.Fatalln(err)
	case err := <-chGRPC:
		log.Fatalln(err)
	}
}
