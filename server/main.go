package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/simplegrpc/stream"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 9009, "The server port")
)

type Connection struct {
	ID     string
	Stream stream.Cast_CreateStreamServer
	Active bool
	err    chan error
}

type Server struct {
	stream.UnimplementedCastServer
	Connections []*Connection
}

func (s *Server) CreateStream(userConn *stream.Connect, stream stream.Cast_CreateStreamServer) error {
	conn := &Connection{
		Stream: stream,
		ID:     userConn.User.ID,
		Active: true,
		err:    make(chan error),
	}
	s.Connections = append(s.Connections, conn)
	return <-conn.err
}

func (s *Server) CastMessage(ctx context.Context, msg *stream.Message) (*stream.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, connection := range s.Connections {
		log.Printf("rececing connection for connection with #%d", connection.ID)
		wait.Add(1)

		go func(msg *stream.Message, conn *Connection) {
			defer wait.Done()

			if conn.Active {
				log.Printf("Sending message %s to user %s", msg.Message, msg.User.ID)
				if err := conn.Stream.Send(msg); err != nil {
					log.Printf("Error with stream %v. Error: %v", conn.Stream, err)
					conn.Active = false
					conn.err <- err
				}
			}
		}(msg, connection)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &stream.Close{}, nil
}

func main() {
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("error listning to tcp on port %d %v", port, err)
	}
	grpcServer := grpc.NewServer()

	stream.RegisterCastServer(grpcServer, &Server{})

	log.Println("starting server......")

	if err := grpcServer.Serve(tcpListener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
