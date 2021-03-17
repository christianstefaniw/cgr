package cgr

import (
	"context"
	"fmt"
	"net/http"
)

// ServeHTTP dispatches the handler registered in the matched route.
func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	defer internalError(&w)

	method := req.Method
	path := req.URL.Path

	r, err := router.routes.search(method, path)

	if err != nil {
		http.Error(w, fmt.Sprintf("Access %s: %s", path, err), http.StatusNotImplemented)
		return
	}

	paramsAsMap := r.params.paramsToMap()

	ctx := context.WithValue(req.Context(), "params", paramsAsMap)

	if r.middleware.len() == 0 {
		r.handlerFunc.ServeHTTP(w, req.WithContext(ctx))
	} else {
		r.executeMiddleware(w, req.WithContext(ctx))
	}
}

func internalError(w *http.ResponseWriter) {
	if r := recover(); r != nil {
		http.Error(*w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}
