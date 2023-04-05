package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/simplegrpc/stream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	address  = flag.String("address", "localhost:9009", "the address to connect to")
	RandomID = func(identifier *string) string {
		randomID := sha256.Sum256([]byte(time.Now().String() + *identifier))
		return hex.EncodeToString(randomID[:])
	}
	messengerClient stream.CastClient
	waitMe          *sync.WaitGroup
)

func init() {
	waitMe = &sync.WaitGroup{}
}

func main() {
	done := make(chan bool)
	clientConn, err := grpc.Dial(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error listning on address %s", *address)
	}
	defer func(clientConn *grpc.ClientConn) {
		err := clientConn.Close()
		if err != nil {
			log.Fatalf("error closing connection %v", err)
		}
	}(clientConn)

	messengerClient = stream.NewCastClient(clientConn)
	name := flag.String("Name", fmt.Sprintf("Guest%s", time.Now().String()), "")
	flag.Parse()
	user := &stream.User{
		ID:   RandomID(name),
		Name: *name,
	}
	if err = connect(user); err != nil {
		log.Printf("error connecting to user %v. please try again", err)
	}
	sendMessages(user)

	go func() {
		waitMe.Wait()
		close(done)
	}()
	<-done
}

func sendMessages(user *stream.User) {
	waitMe.Add(1)
	go func() {
		waitMe.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			message := NewMessage(user, scanner.Text())
			_, err := messengerClient.CastMessage(context.Background(), message)
			if err != nil {
				fmt.Printf("error sending message %s", message)
				return
			}
		}
	}()
}

func connect(user *stream.User) error {
	var streamError error
	streamClient, err := messengerClient.CreateStream(context.Background(), &stream.Connect{
		User:   user,
		Active: true,
	})
	if err != nil {
		log.Fatalf("Error creating stream %v", err)
	}

	waitMe.Add(1)

	go func(streamClient stream.Cast_CreateStreamClient) {
		defer waitMe.Done()
		for {
			msg, err := streamClient.Recv()
			if err != nil {
				streamError = fmt.Errorf("error reading message: %v", err)
				break
			}

			fmt.Printf("%v : %s\n", msg.Id, msg.Message)

		}
	}(streamClient)
	return streamError
}

func NewMessage(user *stream.User, message string) *stream.Message {
	return &stream.Message{
		Id:        RandomID(&user.Name),
		User:      user,
		Message:   message,
		Timestamp: time.Now().String(),
	}
}
