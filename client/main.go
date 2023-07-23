package main

import (
	"context"
	"grpc-lesson/pb"
	"io"
	"log"
	"os"
	"time"

	"fmt"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileServiceClient(conn)
	// callListFiles(client)
	// callDownload(client)
	// CallUpload(client)
	CallUploadAndNotifyProgress(client)
}

// Unary RPC(Client側)
func callListFiles(client pb.FileServiceClient) {
	res, err := client.ListFiles(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res.GetFilenames())
}

// サーバーストリーミングRPC(Client側)
func callDownload(client pb.FileServiceClient) {
	req := &pb.DownloadRequest{Filename: "name.txt"}
	stream, err := client.DownLoad(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Response from Download(bytes): %v", res.GetData())
		log.Printf("Response from Download(string): %v", string(res.GetData()))
	}
}

// クライアントストリーミングRPC（Client側）
func CallUpload(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/go/src/workspace/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 5)

	for {
		n, err := file.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		req := &pb.UploadRequest{Data: buf[:n]}
		sendErr := stream.Send(req)
		if sendErr != nil {
			log.Fatalln(sendErr)
		}

		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("received data size: %v", res.GetSize())
}

// 双方向ストリーミングRPC(Client側)
func CallUploadAndNotifyProgress(client pb.FileServiceClient) {
	filename := "sports.txt"
	path := "/go/src/workspace/storage/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	stream, err := client.UploadAndNotifyProgress(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	// 並行処理(goroutine)
	// request
	buf := make([]byte, 5)
	go func() {
		for {
			// 指定したファイルを読み込む
			n, err := file.Read(buf)
			// データの読み込み終了時にループを抜ける
			if n == 0 || err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}

			req := &pb.UploadAndNotifyProgressRequest{Data: buf[:n]}
			// streamでリクエストを行う
			sendErr := stream.Send(req)
			if sendErr != nil {
				log.Fatalln(err)
			}
			time.Sleep(1 * time.Second)
		}

		// サーバー側にio.EOFが通知される
		err := stream.CloseSend()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// response
	ch := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			// サーバー側からクライアント側でio.EOFが通知されるので、ループを抜ける
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("Receved message: %v", res.GetMsg())
		}
		// チャネルをクローズ
		close(ch)
	}()
	// 待機していたチャネルを抜けて全体の処理が終了する
	<-ch
}
