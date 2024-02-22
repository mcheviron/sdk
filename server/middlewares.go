package server

import (
	"log/slog"
	"net/http"
)

func (s *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Info("Received request",
			slog.String("method", r.Method),
			slog.String("URI", r.URL.Path),
			slog.String("IP", r.RemoteAddr),
			slog.String("ORIGIN", r.Header.Get("Origin")),
			slog.String("HOST", r.Host),
		)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.Error("Panic recovered", slog.String("message", err.(string)))
				s.InternalServerError(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
