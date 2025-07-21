module github.com/study-platform/api-gateway

go 1.23.0

toolchain go1.23.11

require (
	github.com/gorilla/mux v1.8.1
	github.com/study-platform/auth-service v0.0.0
	github.com/study-platform/course-service v0.0.0
	github.com/study-platform/pkg v0.0.0
	github.com/study-platform/progress-service v0.0.0
	google.golang.org/grpc v1.73.0
)

require (
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/study-platform/pkg => ../pkg

replace github.com/study-platform/auth-service => ../auth-service

replace github.com/study-platform/course-service => ../course-service

replace github.com/study-platform/progress-service => ../progress-service
