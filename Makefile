.PHONY: proto clean build test docker-build docker-up docker-down docker-logs

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	export PATH=$$PATH:$$(go env GOPATH)/bin && \
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/auth.proto
	@echo "Copying auth proto files to auth service..."
	cp proto/auth.pb.go auth-service/proto/
	cp proto/auth_grpc.pb.go auth-service/proto/

# Clean generated files
clean:
	rm -f proto/*.pb.go
	rm -f auth-service/proto/*.pb.go
	rm -f course-service/proto/*.pb.go
	rm -f progress-service/proto/*.pb.go

# Build all services
build: proto
	@echo "Building auth service..."
	cd auth-service && go build -o bin/auth-service ./cmd/
	@echo "Build complete"

# Run tests
test:
	@echo "Running tests..."
	cd auth-service && go test ./...
	@echo "Tests complete"

# Run auth service locally
run-auth: build
	@echo "Running auth service..."
	cd auth-service && ./bin/auth-service

# Docker build
docker-build:
	docker-compose build --no-cache

# Docker run
docker-up:
	docker-compose up -d

# Docker down
docker-down:
	docker-compose down

# Docker logs
docker-logs:
	docker-compose logs -f