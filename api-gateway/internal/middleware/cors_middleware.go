package middleware

import (
    "net/http"
    "os"
    "strings"
)

// getAllowedOrigin decides which origin to echo back based on a comma-separated allowlist in CORS_ORIGINS.
// If the request has no Origin header, returns empty string.
func getAllowedOrigin(r *http.Request) (string, bool) {
    origin := r.Header.Get("Origin")
    if origin == "" {
        return "", false
    }
    // Allowlist from env or sensible dev defaults
    allowlist := os.Getenv("CORS_ORIGINS")
    if allowlist == "" {
        allowlist = "http://localhost:3000,http://127.0.0.1:3000"
    }
    // Wildcard support: if '*' present, allow any origin (dev convenience)
    for _, o := range strings.Split(allowlist, ",") {
        if strings.TrimSpace(o) == origin {
            return origin, true
        }
        if strings.TrimSpace(o) == "*" {
            return origin, true
        }
    }
    return "", false
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Always vary on Origin
        w.Header().Add("Vary", "Origin")
        w.Header().Add("Vary", "Access-Control-Request-Method")
        w.Header().Add("Vary", "Access-Control-Request-Headers")

        // Allowed methods and headers - only set if not already set
        if w.Header().Get("Access-Control-Allow-Methods") == "" {
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        }

        // Dynamically allow requested headers (Chrome client-hints like sec-ch-ua*, etc.)
        if w.Header().Get("Access-Control-Allow-Headers") == "" {
            if reqHdrs := r.Header.Get("Access-Control-Request-Headers"); reqHdrs != "" {
                w.Header().Set("Access-Control-Allow-Headers", reqHdrs)
            } else {
                w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Accept-Language")
            }
        }
        if w.Header().Get("Access-Control-Max-Age") == "" {
            w.Header().Set("Access-Control-Max-Age", "86400")
        }

        if origin, ok := getAllowedOrigin(r); ok {
            // Echo back the requesting origin and allow credentials for cookie/token scenarios
            // Only set if not already set to prevent duplication
            if w.Header().Get("Access-Control-Allow-Origin") == "" {
                w.Header().Set("Access-Control-Allow-Origin", origin)
            }
            if w.Header().Get("Access-Control-Allow-Credentials") == "" {
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }
        } else if r.Header.Get("Origin") != "" {
            // SECURITY: Reject requests from non-allowed origins
            // Do not set any CORS headers - browser will block the request
            w.WriteHeader(http.StatusForbidden)
            w.Write([]byte(`{"error":"Origin not allowed","code":"CORS_BLOCKED"}`))
            return
        }

        // Handle preflight requests early
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// SetJSONContentType sets the Content-Type header to application/json
func SetJSONContentType(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}
