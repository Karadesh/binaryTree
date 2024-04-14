package main

import (
	api "binaryTree/api/proto"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	api.UnimplementedBinaryTreeServer
}

func (s *Server) GenerateRequest(ctx context.Context, req *api.GenRequest) (*api.GenResponse, error) {
	log.Println("Got request")
	return &api.GenResponse{Result: binaryTreeResponse(req.GetRequest())}, nil
}

type Node struct {
	Val   int
	Left  *Node
	Right *Node
}

// запуск сервера
func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterBinaryTreeServer(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// "парсинг" запроса
func binaryTreeResponse(command string) string {
	intCommand, err := strconv.Atoi(command)
	if err != nil {
		if command == "showtree" {
			return showTree()
		} else {
			return "invalid command"
		}
	}
	if intCommand < 2 {
		return "Type number that bigger than '2'"
	} else if intCommand > 100000 {
		return "Type number that lower than '100000'"
	} else {
		return makeTree(intCommand)
	}
}

// Insert функция распределения значений дерева и записи в бд
func (n *Node) Insert(val int) {
	db, err := sql.Open("postgres", "user=admin password=admin host=localhost dbname=randnumbers sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()
	if val < n.Val {
		if n.Left == nil {
			n.Left = &Node{Val: val}
			inserter := fmt.Sprintf("INSERT INTO binarytree(number, node) VALUES(%d, %d)", val, n.Val)
			result, err := db.Exec(inserter)
			if err != nil {
				fmt.Printf("failed to insert into randomnumbers: %v", err)
			}
			affectedRows, err := result.RowsAffected()
			fmt.Printf("Updated %d rows\n", affectedRows)
		} else {
			n.Left.Insert(val)
		}
	} else {
		if n.Right == nil {
			n.Right = &Node{Val: val}
			inserter := fmt.Sprintf("INSERT INTO binarytree(number, node) VALUES(%d, %d)", val, n.Val)
			result, err := db.Exec(inserter)
			if err != nil {
				fmt.Printf("failed to insert into randomnumbers: %v", err)
			}
			affectedRows, err := result.RowsAffected()
			fmt.Printf("Updated %d rows\n", affectedRows)
		} else {
			n.Right.Insert(val)
		}
	}
}

// входная точка создания дерева, если нужно сгенерировать новое.
func makeTree(nodesNumber int) string {
	db, err := sql.Open("postgres", "user=admin password=admin host=localhost dbname=randnumbers sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()
	_, err = db.Exec("TRUNCATE TABLE binarytree")
	if err != nil {
		log.Fatalf("failed to drop binaryTree: %v", err)
	}
	root := &Node{Val: rand.Intn(100)}
	for i := 0; i < nodesNumber; i++ {
		root.Insert(rand.Intn(100))
	}
	return "Done!"
}

// функция тянет дерево из бд и выводит его в виде строки
func showTree() string {
	db, err := sql.Open("postgres", "user=admin password=admin host=localhost dbname=randnumbers sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()
	var m []string
	selector := fmt.Sprintf("SELECT * FROM binarytree")
	somerow, err := db.Query(selector)
	if err != nil {
		log.Fatalf("failed to query: %v", err)
	}
	defer somerow.Close()
	for somerow.Next() {
		var someNumber int
		var someNode int
		something := somerow.Scan(&someNumber, &someNode)
		toString := fmt.Sprintf("(node:%d, value:%d)", someNode, someNumber)
		m = append(m, toString)
		fmt.Println(something)
	}
	return strings.Join(m, ",")
}
