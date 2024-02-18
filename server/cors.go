package server

import (
	"net/http"
	"strings"
)

type CORSOptions struct {
	AllowedOrigins []string
	AllowedHeaders []string
	AllowedMethods []string
}

func CORS(options CORSOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				allowed := false
				for _, o := range options.AllowedOrigins {
					if strings.EqualFold(origin, o) {
						allowed = true
						break
					}
				}

				if allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					if len(options.AllowedMethods) > 0 {
						w.Header().Set("Access-Control-Allow-Methods", strings.Join(options.AllowedMethods, ", "))
					} else {
						w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
					}

					if len(options.AllowedHeaders) > 0 {
						w.Header().Set("Access-Control-Allow-Headers", strings.Join(options.AllowedHeaders, ", "))
					} else {
						w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
