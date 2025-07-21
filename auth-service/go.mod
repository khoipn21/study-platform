module github.com/study-platform/auth-service

go 1.23.0

toolchain go1.24.5

require (
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/google/uuid v1.6.0
	github.com/study-platform/pkg v0.0.0
	golang.org/x/crypto v0.36.0
	golang.org/x/oauth2 v0.28.0
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
)

replace github.com/study-platform/pkg => ../pkg
