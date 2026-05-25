APP_NAME=mm_connect

GRPC_CMD=./cmd/grpc
GATEWAY_CMD=./cmd/gateway

PROTO_DIR=proto
GEN_DIR=gen

.PHONY: proto build run-grpc run-gateway clean test

# =========================
# BUILD
# =========================

build:
	go build -o bin/grpc $(GRPC_CMD)
	go build -o bin/gateway $(GATEWAY_CMD)

# =========================
# RUN (LOCAL DEV)
# =========================

run-grpc:
	go run $(GRPC_CMD)

run-gateway:
	go run $(GATEWAY_CMD)

# =========================
# PROTO GENERATION
# =========================

proto:
	protoc \
    -I proto \
    -I third_party/googleapis \
    --go_out proto --go_opt paths=source_relative \
    --go-grpc_out proto --go-grpc_opt paths=source_relative \
    --grpc-gateway_out proto --grpc-gateway_opt paths=source_relative \
    proto/auth/v1/*.proto proto/auth/v1/client/*.proto

# =========================
# CLEAN
# =========================

clean:
	rm -rf bin

# =========================
# DB MIGRATIONS
# =========================

DB_URL=mysql://root@tcp(localhost:3306)/mm_connect

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-force:
	migrate -path migrations -database "$(DB_URL)" force $(version)

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down 

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-version:
	migrate -path migrations -database "$(DB_URL)" version

migrate-drop:
	migrate -path migrations -database "$(DB_URL)" drop

migrate-drop-force:
	migrate -path migrations -database "$(DB_URL)" drop -f