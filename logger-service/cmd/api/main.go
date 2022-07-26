package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// register rpc server
	err = rpc.Register(new(RPCServer))
	// start serve rpc
	go app.rpcListen()
	// start serve grpc
	go app.gRPCListen()

	log.Printf("Starting broker service on port %s", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port: ", rpcPort)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		rpcConn, err := listener.Accept()

		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	log.Printf("Starting connecting to Mongo on url %s", mongoURL)

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting", err)
		return nil, err
	}

	return c, nil
}
