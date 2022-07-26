package main

import (
	"context"
	"log"
	"logger/data"
	"logger/logs"
	"net"

	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	// get input
	input := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	res := &logs.LogResponse{Result: "logged!"}
	return res, nil
}

func (app *Config) gRPCListen() {
	listener, err := net.Listen("tcp", ":"+gRpcPort)
	if err != nil {
		log.Fatal("Failed to listen to gRPC: " + err.Error())
	}

	server := grpc.NewServer()

	logs.RegisterLogServiceServer(server, &LogServer{Models: app.Models})

	log.Print("gRPC Server stated on port: " + gRpcPort)

	if err := server.Serve(listener); err != nil {
		log.Fatal("Failed to listen to gRPC: " + err.Error())
	}
}
