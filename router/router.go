package router

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
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

func (r *Router) Method(method, route string, handler http.HandlerFunc) {
	if strings.HasSuffix(route, "/") {
		route += "{$}"
	}

	r.routes[fmt.Sprintf("%s %s", method, route)] = r.withMiddlewares(handler)
}

func (r *Router) Get(route string, handler http.HandlerFunc) {
	r.Method(http.MethodGet, route, handler)
}

func (r *Router) Post(route string, handler http.HandlerFunc) {
	r.Method(http.MethodPost, route, handler)
}

func (r *Router) Put(route string, handler http.HandlerFunc) {
	r.Method(http.MethodPut, route, handler)
}

func (r *Router) Patch(route string, handler http.HandlerFunc) {
	r.Method(http.MethodPatch, route, handler)
}

func (r *Router) Delete(route string, handler http.HandlerFunc) {
	r.Method(http.MethodDelete, route, handler)
}

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
