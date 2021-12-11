package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	pb "github.com/darianfd99/grpc/proto-grpc"
)

const (
	port          = ":50051"
	timeout       = 10
	mongoPort     = 12307
	mongoHost     = "localhost"
	mongoDatabase = "db"
	mongoCollection
)

type server struct {
	pb.UnimplementedGetInfoServer
}

func SaveComment(comment string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	mongoURI := fmt.Sprintf("mongodb://%s:%d", mongoHost, mongoPort)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	databases, err := mongoClient.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(databases)

	database := mongoClient.Database(mongoDatabase)
	collection := database.Collection(mongoCollection)

	var bdoc interface{}

	err = bson.UnmarshalExtJSON([]byte(comment), true, &bdoc)
	if err != nil {
		log.Fatal(err)
	}

	insertResult, err := collection.InsertOne(ctx, bdoc)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(insertResult)

}

func (s *server) ReturnInfo(ctx context.Context, in *pb.RequestId) (*pb.ReplyInfo, error) {
	SaveComment(in.GetId())
	fmt.Printf("Received data: %v\n", in.GetId())

	return &pb.ReplyInfo{Info: "Received comment: " + in.GetId()}, nil
}

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()
	pb.RegisterGetInfoServer(srv, &server{})

	if err := srv.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
