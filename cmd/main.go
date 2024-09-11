package main

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"

	simplepb "github.com/farir1408/grpc_transcodding/pkg/pb"
)

type server struct {
	simplepb.UnimplementedSampleRPCServer

	storage map[uint32]*simplepb.Item
}

func NewServer() *server {
	return &server{
		storage: make(map[uint32]*simplepb.Item, 10),
	}
}

func (s *server) SaveItem(ctx context.Context, in *simplepb.SaveItemReq) (*simplepb.SaveItemRes, error) {
	if _, ok := s.storage[in.ItemId]; ok {
		return nil, errors.New("already exist")
	}
	s.storage[in.ItemId] = &simplepb.Item{
		ItemId:  in.ItemId,
		Author:  in.Author,
		Message: in.Message,
		Title:   in.Title,
	}
	return &simplepb.SaveItemRes{Message: "item saved", Success: true}, nil
}

func (s *server) GetItem(ctx context.Context, in *simplepb.GetItemReq) (*simplepb.GetItemRes, error) {
	item, ok := s.storage[in.ItemId]
	if !ok {
		return nil, errors.New("not found")
	}
	return &simplepb.GetItemRes{
		Item: item,
	}, nil
}

func (s *server) UpdateItem(ctx context.Context, in *simplepb.UpdateItemReq) (*simplepb.UpdateItemRes, error) {
	item, ok := s.storage[in.ItemId]
	if !ok {
		return nil, errors.New("not found")
	}

	if in.Message != "" {
		item.Message = in.Message
	}

	if in.Title != "" {
		item.Title = in.Title
	}

	s.storage[in.ItemId] = item

	return &simplepb.UpdateItemRes{
		Item:   item,
		Status: true,
	}, nil
}

func (s *server) DeleteItem(ctx context.Context, in *simplepb.DeleteItemReq) (*simplepb.DeleteItemRes, error) {
	_, ok := s.storage[in.ItemId]
	if !ok {
		return nil, errors.New("not found")
	}

	delete(s.storage, in.ItemId)

	return &simplepb.DeleteItemRes{
		Status: true,
	}, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	//gRPC Gateway
	conn, err := grpc.NewClient(
		"0.0.0.0:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = simplepb.RegisterSampleRPCHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":50052",
		Handler: gwmux,
	}

	go func() {
		log.Println("Serving gRPC-Gateway on http://0.0.0.0:50052")
		log.Fatalln(gwServer.ListenAndServe())
	}()
	//---

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	simplepb.RegisterSampleRPCServer(s, NewServer())
	reflection.Register(s)
	// Serve gRPC Server
	log.Println("Serving gRPC on 0.0.0.0:50051")
	log.Fatal(s.Serve(lis))
}
