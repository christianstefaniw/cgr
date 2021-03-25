package cgr

import "net/http"

type middlewareHandle func(w http.ResponseWriter, r *http.Request)

type middleware struct {
	handler middlewareHandle
}

func NewMiddleware(handler middlewareHandle) *middleware {
	return &middleware{
		handler: handler,
	}
}

func (middleware *middleware) run(w http.ResponseWriter, r *http.Request) {
	middleware.handler(w, r)
}
