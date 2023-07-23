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
	// TODO: 認証処理から実装
	// grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	// grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
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

// 双方向ストリーミングRPC(Server側)
func (*server) UploadAndNotifyProgress(stream pb.FileService_UploadAndNotifyProgressServer) error {
	fmt.Println("UploadAndNotifyProgress wad Invoked")

	size := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatal(err)
		}

		data := req.GetData()
		log.Printf("received data: %v", data)
		size += len(data)

		res := &pb.UploadAndNotifyProgressResponse{
			Msg: fmt.Sprintf("received %vbytes", size),
		}
		err = stream.Send(res)
		if err != nil {
			return err
		}
	}
}

func myLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Printf("request data: %+v", req)

		// handlerはクライアントからコールされたメソッド
		resp, err = handler(ctx, req)
		if err != nil {
			return nil, err
		}
		log.Printf("response data: %+v", resp)

		return resp, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Fail to lesson: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(myLogging()))
	pb.RegisterFileServiceServer(s, &server{})

	fmt.Println("server is running...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to Serve: %v", err)
	}
}
