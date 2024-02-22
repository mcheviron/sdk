package router

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"
)

type Router struct {
	routes      map[string]http.Handler
	middlewares []func(http.Handler) http.Handler
}

func New() *Router {
	return &Router{
		routes:      make(map[string]http.Handler),
		middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

func (r *Router) method(method, route string, handler http.HandlerFunc) {
	if strings.HasSuffix(route, "/") {
		route += "{$}"
	}

	r.routes[fmt.Sprintf("%s %s", method, route)] = r.withMiddlewares(handler)
}

// Get registers a GET route with the specified route pattern and handler function.
func (r *Router) Get(route string, handler http.HandlerFunc) {
	r.method(http.MethodGet, route, handler)
}

// Post registers a new route with the HTTP method POST.
func (r *Router) Post(route string, handler http.HandlerFunc) {
	r.method(http.MethodPost, route, handler)
}

// Put registers a PUT route with the specified route pattern and handler function.
func (r *Router) Put(route string, handler http.HandlerFunc) {
	r.method(http.MethodPut, route, handler)
}

// Patch registers a PATCH route with the specified route pattern and handler function.
func (r *Router) Patch(route string, handler http.HandlerFunc) {
	r.method(http.MethodPatch, route, handler)
}

// Delete registers a DELETE route with the specified route pattern and handler function.
// The handler function should have the signature http.HandlerFunc.
func (r *Router) Delete(route string, handler http.HandlerFunc) {
	r.method(http.MethodDelete, route, handler)
}

// Use adds the provided middlewares to the router's middleware stack.
// The middlewares are executed in reverse order, meaning the last middleware added will be the first to be executed.
// This method allows for chaining multiple middlewares together to handle requests.
func (r *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	slices.Reverse(middlewares)
	r.middlewares = slices.Concat(middlewares, r.middlewares)
}

func (r *Router) withMiddlewares(handler http.Handler) http.Handler {
	for _, middleware := range r.middlewares {
		handler = middleware(handler)
	}
	return handler
}

// Serve starts the HTTP server and listens for incoming requests on the specified address.
// It registers all the routes defined in the Router and uses the provided address to bind the server.
// The server is configured with a read timeout of 45 seconds.
// Serve returns an error if the server fails to start or encounters an error while serving requests.
func (r *Router) Serve(addr string) error {
	mux := http.NewServeMux()

	for route, handler := range r.routes {
		mux.Handle(route, handler)
	}

	server := &http.Server{
		ReadTimeout: 45 * time.Second,
		Addr:        addr,
		Handler:     mux,
	}

	return server.ListenAndServe()
}
