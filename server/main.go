package main

import (
	"bytes"
	"context"
	"fmt"
	"grpc-lesson/pb"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

// Unary RPC(Server側)
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

// サーバーストリーミングRPC(Server側)
func (*server) DownLoad(req *pb.DownloadRequest, stream pb.FileService_DownLoadServer) error {
	fmt.Println("Download was Invoked...")

	filename := req.GetFilename()

	path := "/go/src/workspace/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		res := &pb.DownloadResponse{Data: buf[:n]}
		sendErr := stream.Send(res)
		if sendErr != nil {
			return sendErr
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

// クライアンストリーミングRPC(Server側)
// リクエストから送られてきたデータを順にbufに書き込み、io.EOFが送信されるとbufのサイズを返す
func (*server) Upload(stream pb.FileService_UploadServer) error {
	fmt.Println("Upload was invoked")

	var buf bytes.Buffer
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := &pb.UploadResponse{Size: int32(buf.Len())}
			return stream.SendAndClose(res)
		}
		if err != nil {
			return err
		}

		data := req.GetData()
		log.Printf("received data(bytes): %v", data)
		log.Printf("received data(string): %v", string(data))
		buf.Write(data)
	}
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
