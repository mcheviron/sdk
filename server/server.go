package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mcheviron/sdk/logger"
	"github.com/mcheviron/sdk/serde"
)

type Server struct {
	logger.Logger
}

func NewDevelopment() (*Server, error) {
	logger, err := logger.New("dev")
	if err != nil {
		return nil, fmt.Errorf("logger: %w", err)
	}
	return &Server{
		Logger: logger,
	}, nil
}

func NewProduction() (*Server, error) {
	logger, err := logger.New("prod")
	if err != nil {
		return nil, fmt.Errorf("logger: %w", err)
	}
	return &Server{
		Logger: logger,
	}, nil
}

func (s *Server) response(w http.ResponseWriter, statusCode int, v any) {
	if err := serde.Encode(w, statusCode, v); err != nil {
		s.Logger.Error("encoding json response", slog.String("error", err.Error()), slog.Any("data", v))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	}
}

func (s *Server) Ok(w http.ResponseWriter, v any) {
	s.response(w, http.StatusOK, v)
}

func (s *Server) Created(w http.ResponseWriter, v any) {
	s.response(w, http.StatusCreated, v)
}

func (s *Server) NoContent(w http.ResponseWriter) {
	s.response(w, http.StatusNoContent, nil)
}

func (s *Server) BadRequest(w http.ResponseWriter, err error) {
	s.response(w, http.StatusBadRequest, err)
}

func (s *Server) Unauthorized(w http.ResponseWriter) {
	s.response(w, http.StatusUnauthorized, nil)
}

func (s *Server) Forbidden(w http.ResponseWriter) {
	s.response(w, http.StatusForbidden, nil)
}

func (s *Server) NotFound(w http.ResponseWriter, err error) {
	s.response(w, http.StatusNotFound, err)
}

func (s *Server) InternalServerError(w http.ResponseWriter) {
	s.response(w, http.StatusInternalServerError, nil)
}
