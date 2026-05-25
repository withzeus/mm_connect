package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	serviceauthv1 "github.com/withzeus/id_contracts/gen/go/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := serviceauthv1.RegisterClientAuthServiceHandlerFromEndpoint(
		ctx, mux, "localhost:9090", opts,
	)

	if err != nil {
		log.Fatalf("Failed to register service handler from endpoint: %v", err)
	}

	log.Println("REST gateway running on :8088")
	log.Fatal(http.ListenAndServe(":8088", mux))
}
