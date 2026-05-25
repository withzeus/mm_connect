package main

import (
	"log"
	"net"

	clientv1 "github.com/withzeus/id_contracts/gen/go/service/v1"
	"github.com/withzeus/mm_connect/internal/auth/jwt"
	"github.com/withzeus/mm_connect/internal/repository/mysql"
	"github.com/withzeus/mm_connect/internal/service"
	igrpc "github.com/withzeus/mm_connect/internal/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	db, err := mysql.NewDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	repo := mysql.NewClientAuthRepository(db)
	jwt := jwt.New("dummy-secret")

	s := service.NewAuthService(repo, jwt)
	svr := grpc.NewServer()

	clientv1.RegisterClientAuthServiceServer(svr, igrpc.New(s))

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gRpc listening on :9090")
	if err := svr.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
