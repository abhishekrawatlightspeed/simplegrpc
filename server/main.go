package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/simplegrpc/stream"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 9009, "The server port")
)

type Server struct {
	stream.UnsafeStreamServiceServer
}

func (s Server) FetchResponse(in *stream.Request, srv stream.StreamService_FetchResponseServer) error {

	log.Printf("fetching response for id - %d", in.ID)

	//use wait group to allow process to be concurrent
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(count int64) {
			defer wg.Done()

			//time sleep to simulate server process time
			time.Sleep(time.Duration(count) * time.Second)
			resp := stream.Response{Resp: fmt.Sprintf("Request #%d For Id:%d", count, in.ID)}
			if err := srv.Send(&resp); err != nil {
				log.Printf("error sending message %v", err)
			}
			log.Printf("finishing request  : %d", count)
		}(int64(i))
	}

	wg.Wait()
	return nil
}

func main() {
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("error listning to tcp on port %d %v", port, err)
	}
	grpcServer := grpc.NewServer()

	stream.RegisterStreamServiceServer(grpcServer, &Server{})

	log.Println("start server")

	if err := grpcServer.Serve(tcpListener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	if err := grpcServer.Serve(tcpListener); err != nil {
		log.Fatalf("error occurred when grpc server is listening to tcp listerner %v", err)
	}
}
