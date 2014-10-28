package session

import (
	"log"
	"net/http"
)

type Middleware struct {
	session bool
}

// Middleware is a struct that has a ServeHTTP method
func NewMiddleware() *Middleware {
	return &Middleware{true}
}

// The middleware handler
func (l *Middleware) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	// Call the next middleware handler
	next(w, req)
}
