package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"google.golang.org/grpc"
	student_service "http_grpc_bench/micro_service/grpc"
	"io"
	"log"
	"net"
	"net/http"
)

type Server struct {
	student_service.UnimplementedStudentServiceServer
}

func (s *Server) GetStudent(ctx context.Context, id *student_service.StudentID) (*student_service.Student, error) {
	return &student_service.Student{
		Name:      "hu",
		CreatedAt: 1990,
		Scores:    nil,
		Locations: nil,
		Gender:    true,
		Age:       35,
		Height:    42,
		Id:        id.GetId() + 2,
	}, nil
}

func startGrpcServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("fail to listen: %v", err)
	}

	s := grpc.NewServer()
	student_service.RegisterStudentServiceServer(
		s,
		&Server{})
	log.Printf("start grpc server listen on %d", port)

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("fail to serve grpc. %v", err)
	}
}

func startHttpServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var req student_service.StudentID
		var buf, err = io.ReadAll(request.Body)
		if err != nil {
			log.Fatalf("fail to read request body. err: %v", err)
		}
		request.Body.Close()
		err = sonic.Unmarshal(buf, &req)
		if err != nil {
			log.Fatalf("fail to unmarshal request. err: %v", err)
		}
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		var resp = &student_service.Student{
			Name:      "hu",
			CreatedAt: 1990,
			Scores:    nil,
			Locations: nil,
			Gender:    true,
			Age:       35,
			Height:    42,
			Id:        req.GetId() + 2,
		}
		buf, err = json.Marshal(resp)
		if err != nil {
			log.Fatalf("fail to marhal response. err: %v", err)
		}
		writer.Write(buf)
	})

	var addr = fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("start http server listen on %d", port)
	var err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("fail to serve http. %v", err)
	}
}

func main() {
	var grpcPort = 5678
	var httpPort = 5679

	go startGrpcServer(grpcPort)
	go startHttpServer(httpPort)

	select {}
}
