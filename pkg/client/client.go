package main

import (
	api "binaryTree/api/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const address = "localhost:8080"

// Запуск сервера
func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewBinaryTreeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	//бесконечный цикл запрос-ожидание ответа
	for {
		log.Printf("Do you want to build binary tree?\n " +
			"Type 'yes' to build\n " +
			"Type 'showtree' to show tree\n " +
			"Type type 'exit' to exit\n")
		var request string
		fmt.Scanln(&request)
		if request == "exit" {
			break
		} else if request == "yes" {
			log.Printf("Type number of nodes to build binary tree from 2 to 100000")
			fmt.Scanln(&request)
		}
		//отправка реквеста
		r, err := c.GenerateRequest(ctx, &api.GenRequest{Request: request})
		if err != nil {
			log.Fatalf("could not generate request: %v", err)
		}
		//ответ от сервера
		log.Println(r.GetResult())
	}
}
