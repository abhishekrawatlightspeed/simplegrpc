package main

import (
	"context"
	"flag"
	"io"
	"log"

	"github.com/simplegrpc/stream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	address = flag.String("address", "localhost:9009", "the address to connect to")
)

func main() {
	clientConn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error listning on address %s", *address)
	}

	client := stream.NewStreamServiceClient(clientConn)
	ctx := context.Background()
	input := &stream.Request{ID: 101}
	response, err := client.FetchResponse(ctx, input)
	if err != nil {
		log.Fatalf("error fetching response from stream %v", err)
	}

	listningChan := make(chan bool)

	go func() {
		for {
			recv, err := response.Recv()
			if err != nil {
				if err == io.EOF {
					listningChan <- true
					return
				}
				log.Fatalf("error receiving response %v", err)
			}
			log.Printf("Message received %s", recv.Resp)
		}
	}()
	<-listningChan
	log.Println("finished listening to all the messages")
}
