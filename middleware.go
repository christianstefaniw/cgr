package cgr

type middlewareHandle func()

type middleware struct {
	handler middlewareHandle
}

func NewMiddleware(handler middlewareHandle) *middleware {
	return &middleware{
		handler: handler,
	}
}

func (middleware *middleware) run() {
	middleware.handler()
}
