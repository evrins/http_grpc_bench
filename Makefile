generate:
	protoc \
	--go_out=./micro_service/grpc \
	--go-grpc_out=./micro_service/grpc \
	--proto_path=./micro_service/grpc \
	student.proto

benchmark: benchmarkHttp benchmarkGrpc

benchmarkHttp:
	go test -bench=^BenchmarkHttp -run=^$ -count=1 -benchmem -benchtime=5s

benchmarkGrpc:
	go test -bench=^BenchmarkGrpc -run=^$ -count=1 -benchmem -benchtime=5s
