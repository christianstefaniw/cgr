package cgr

import "net/http"

type middlewareHandle func(w http.ResponseWriter, r *http.Request)

type Middleware struct {
	handler middlewareHandle
}

func NewMiddleware(handler middlewareHandle) *Middleware {
	return &Middleware{
		handler: handler,
	}
}

func (middleware *Middleware) run(w http.ResponseWriter, r *http.Request) {
	middleware.handler(w, r)
}
