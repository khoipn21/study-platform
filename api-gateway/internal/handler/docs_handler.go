package handler

import (
	"net/http"
)

type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

// GetAPISpec returns the OpenAPI specification for the API Gateway
func (h *DocsHandler) GetAPISpec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// OpenAPI 3.0 specification
	spec := `{
  "openapi": "3.0.0",
  "info": {
    "title": "Study Platform API",
    "description": "API Gateway for the Study Platform microservices architecture",
    "version": "1.0.0",
    "contact": {
      "name": "Study Platform Team",
      "email": "support@studyplatform.com"
    }
  },
  "servers": [
    {
      "url": "http://localhost:8080/api/v1",
      "description": "Development server"
    }
  ],
  "paths": {
    "/health": {
      "get": {
        "summary": "Health check endpoint",
        "description": "Returns the health status of the API Gateway",
        "responses": {
          "200": {
            "description": "API Gateway is healthy",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "status": {"type": "string", "example": "healthy"},
                    "service": {"type": "string", "example": "api-gateway"},
                    "version": {"type": "string", "example": "1.0.0"}
                  }
                }
              }
            }
          }
        }
      }
    },
    "/auth/register": {
      "post": {
        "summary": "User registration",
        "description": "Register a new user account",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["username", "email", "password"],
                "properties": {
                  "username": {"type": "string", "example": "johndoe"},
                  "email": {"type": "string", "format": "email", "example": "john@example.com"},
                  "password": {"type": "string", "minLength": 8, "example": "password123"}
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "User registered successfully",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "success": {"type": "boolean", "example": true},
                    "message": {"type": "string", "example": "User registered successfully"},
                    "data": {
                      "type": "object",
                      "properties": {
                        "token": {"type": "string", "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
                        "user": {
                          "type": "object",
                          "properties": {
                            "id": {"type": "string", "example": "123e4567-e89b-12d3-a456-426614174000"},
                            "username": {"type": "string", "example": "johndoe"},
                            "email": {"type": "string", "example": "john@example.com"},
                            "role": {"type": "string", "example": "student"}
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          },
          "400": {
            "description": "Invalid input"
          },
          "409": {
            "description": "User already exists"
          }
        }
      }
    },
    "/auth/login": {
      "post": {
        "summary": "User login",
        "description": "Authenticate a user and return a JWT token",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["email", "password"],
                "properties": {
                  "email": {"type": "string", "format": "email", "example": "john@example.com"},
                  "password": {"type": "string", "example": "password123"}
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Login successful",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "success": {"type": "boolean", "example": true},
                    "message": {"type": "string", "example": "Login successful"},
                    "data": {
                      "type": "object",
                      "properties": {
                        "token": {"type": "string", "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."},
                        "user": {
                          "type": "object",
                          "properties": {
                            "id": {"type": "string", "example": "123e4567-e89b-12d3-a456-426614174000"},
                            "username": {"type": "string", "example": "johndoe"},
                            "email": {"type": "string", "example": "john@example.com"},
                            "role": {"type": "string", "example": "student"}
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          },
          "401": {
            "description": "Invalid credentials"
          }
        }
      }
    },
    "/courses": {
      "get": {
        "summary": "List courses",
        "description": "Get a paginated list of courses",
        "parameters": [
          {
            "name": "page",
            "in": "query",
            "description": "Page number",
            "schema": {"type": "integer", "default": 1}
          },
          {
            "name": "page_size",
            "in": "query",
            "description": "Number of items per page",
            "schema": {"type": "integer", "default": 10, "maximum": 100}
          }
        ],
        "responses": {
          "200": {
            "description": "List of courses",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "success": {"type": "boolean", "example": true},
                    "message": {"type": "string", "example": "Courses retrieved successfully"},
                    "data": {
                      "type": "object",
                      "properties": {
                        "courses": {
                          "type": "array",
                          "items": {
                            "type": "object",
                            "properties": {
                              "id": {"type": "string", "example": "123e4567-e89b-12d3-a456-426614174000"},
                              "title": {"type": "string", "example": "Introduction to Go Programming"},
                              "description": {"type": "string", "example": "Learn the fundamentals of Go programming language"},
                              "instructor_id": {"type": "string", "example": "123e4567-e89b-12d3-a456-426614174001"},
                              "category": {"type": "string", "example": "Programming"},
                              "level": {"type": "string", "example": "beginner"},
                              "price": {"type": "number", "example": 99.99},
                              "currency": {"type": "string", "example": "USD"},
                              "status": {"type": "string", "example": "published"}
                            }
                          }
                        },
                        "total": {"type": "integer", "example": 100},
                        "page": {"type": "integer", "example": 1},
                        "page_size": {"type": "integer", "example": 10},
                        "total_pages": {"type": "integer", "example": 10}
                      }
                    }
                  }
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Create course",
        "description": "Create a new course (requires authentication)",
        "security": [{"bearerAuth": []}],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["title", "description", "category", "level", "price"],
                "properties": {
                  "title": {"type": "string", "example": "Introduction to Go Programming"},
                  "description": {"type": "string", "example": "Learn the fundamentals of Go programming language"},
                  "category": {"type": "string", "example": "Programming"},
                  "level": {"type": "string", "enum": ["beginner", "intermediate", "advanced"], "example": "beginner"},
                  "price": {"type": "number", "minimum": 0, "example": 99.99},
                  "currency": {"type": "string", "example": "USD"},
                  "tags": {"type": "array", "items": {"type": "string"}, "example": ["go", "programming", "backend"]},
                  "thumbnail_url": {"type": "string", "format": "uri", "example": "https://example.com/thumbnail.jpg"}
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Course created successfully"
          },
          "401": {
            "description": "Unauthorized"
          },
          "400": {
            "description": "Invalid input"
          }
        }
      }
    },
    "/enrollments": {
      "get": {
        "summary": "List user enrollments",
        "description": "Get a list of courses the user is enrolled in",
        "security": [{"bearerAuth": []}],
        "responses": {
          "200": {
            "description": "List of enrollments"
          },
          "401": {
            "description": "Unauthorized"
          }
        }
      },
      "post": {
        "summary": "Create enrollment",
        "description": "Enroll in a course",
        "security": [{"bearerAuth": []}],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": ["course_id"],
                "properties": {
                  "course_id": {"type": "string", "example": "123e4567-e89b-12d3-a456-426614174000"}
                }
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Enrollment created successfully"
          },
          "401": {
            "description": "Unauthorized"
          },
          "409": {
            "description": "Already enrolled"
          }
        }
      }
    }
  },
  "components": {
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    }
  }
}`
	
	w.Write([]byte(spec))
}

// GetSwaggerUI returns the Swagger UI HTML page
func (h *DocsHandler) GetSwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Study Platform API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3.52.5/swagger-ui-bundle.js"></script>
    <script>
        SwaggerUIBundle({
            url: '/api/v1/docs/openapi.json',
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIBundle.presets.standalone
            ],
            plugins: [
                SwaggerUIBundle.plugins.DownloadUrl
            ],
            layout: "StandaloneLayout"
        })
    </script>
</body>
</html>`
	
	w.Write([]byte(html))
}