package middleware

import (
	"net/http"
)

type Middlewares struct {
	// handlers []http.Handler
}

func (m *Middlewares) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
