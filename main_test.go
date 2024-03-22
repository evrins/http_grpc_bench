package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/bytedance/sonic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	student_service "http_grpc_bench/micro_service/grpc"
	"io"
	"net/http"
	"testing"
	"time"
)

const studentId = 10

func TestGrpc(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		"127.0.0.1:5678",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("fail to dial. err: %v", err)
	}

	client := student_service.NewStudentServiceClient(
		conn,
	)

	request := student_service.StudentID{Id: studentId}
	response, err := client.GetStudent(context.Background(), &request)
	if err != nil {
		t.Fatalf("get student failed. err: %v", err)
	}

	t.Log(response.Id)
}

func TestHttp(t *testing.T) {
	client := http.Client{}

	sid := student_service.StudentID{Id: studentId}
	var buf = bytes.NewBuffer(nil)
	var err = json.NewEncoder(buf).Encode(&sid)
	if err != nil {
		t.Fatalf("fail to encode student id. err: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:5679", buf)
	if err != nil {
		t.Fatalf("fail to new request. err: %v", err)
	}

	resp, err := client.Do(request)
	if err != nil {
		t.Fatalf("fail to make http request. err: %v", err)
	}
	defer resp.Body.Close()
	var s student_service.Student
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		t.Fatalf("fail to decode response body. err: %v", err)
	}
	t.Log(s.Id)
}

func BenchmarkGrpc(b *testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		"127.0.0.1:5678",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		b.Fatalf("fail to dial. err: %v", err)
	}

	client := student_service.NewStudentServiceClient(conn)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		request := student_service.StudentID{Id: studentId}
		client.GetStudent(context.Background(), &request)
	}
}

func BenchmarkHttp(b *testing.B) {
	client := http.Client{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sid := student_service.StudentID{Id: studentId}
		var buf, err = sonic.Marshal(&sid)
		if err != nil {
			b.Fatalf("fail to encode student id. err: %v", err)
		}

		request, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:5679", bytes.NewReader(buf))
		if err != nil {
			b.Fatalf("fail to new request. err: %v", err)
		}

		resp, err := client.Do(request)
		if err != nil {
			b.Fatalf("fail to make http request. err: %v", err)
		}
		buf, err = io.ReadAll(resp.Body)
		if err != nil {
			b.Fatalf("fail to read response body. err: %v", err)
		}
		resp.Body.Close()
		var s student_service.Student
		err = sonic.Unmarshal(buf, &s)
		if err != nil {
			b.Fatalf("fail to decode response body. err: %v", err)
		}
	}
}
