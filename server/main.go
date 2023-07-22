package main

import (
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

func (*server) ListFiles(ctx context.Context, req *pb.ListFilesRequest) (res *pb.ListFilesResponse, err error) {
	fmt.Println("ListFiles was invoked")

	dir := "/go/src/workspace/storage"

	paths, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	var filenames []string
	for _, path := range paths {
		if !path.IsDir() {
			filenames = append(filenames, path.Name())
		}
	}

	res = &pb.ListFilesResponse{
		Filenames: filenames,
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Fail to lesson: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFileServiceServer(s, &server{})

	fmt.Println("server is running...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}
}
