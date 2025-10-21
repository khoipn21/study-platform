module github.com/study-platform/api-gateway

go 1.24.0

toolchain go1.24.5

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.1
	github.com/study-platform/auth-service v0.0.0
	github.com/study-platform/course-service v0.0.0
	github.com/study-platform/pkg v0.0.0
	github.com/study-platform/progress-service v0.0.0
	github.com/swaggo/http-swagger v1.3.4
	github.com/swaggo/swag v1.16.6
	google.golang.org/grpc v1.73.0
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/go-openapi/jsonpointer v0.22.1 // indirect
	github.com/go-openapi/jsonreference v0.21.2 // indirect
	github.com/go-openapi/spec v0.22.0 // indirect
	github.com/go-openapi/swag/conv v0.25.1 // indirect
	github.com/go-openapi/swag/jsonname v0.25.1 // indirect
	github.com/go-openapi/swag/jsonutils v0.25.1 // indirect
	github.com/go-openapi/swag/loading v0.25.1 // indirect
	github.com/go-openapi/swag/stringutils v0.25.1 // indirect
	github.com/go-openapi/swag/typeutils v0.25.1 // indirect
	github.com/go-openapi/swag/yamlutils v0.25.1 // indirect
	github.com/swaggo/files v1.0.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/mod v0.29.0 // indirect
	golang.org/x/net v0.46.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	golang.org/x/tools v0.38.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/study-platform/pkg => ../pkg

replace github.com/study-platform/auth-service => ../auth-service

replace github.com/study-platform/course-service => ../course-service

replace github.com/study-platform/progress-service => ../progress-service
