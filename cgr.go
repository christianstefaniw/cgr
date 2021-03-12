package cgr

import (
	"context"
	"fmt"
	"log"
	"net/http"
)



// ServeHTTP dispatches the handler registered in the matched route.
func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	defer internalError(&w)

	method := req.Method
	path := req.URL.Path

	r, err := router.routes.search(method, path)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err), http.StatusNotImplemented)
		return
	}
	paramsAsMap := r.params.paramsToMap()

	ctx := context.WithValue(req.Context(), "params", paramsAsMap)
	r.handlerFunc.ServeHTTP(w, req.WithContext(ctx))
}


func internalError(w *http.ResponseWriter) {
	if r := recover(); r != nil {
		http.Error(*w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

// Run attaches the Router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, Router)
func Run(port string, router *Router) {

	for _, warning := range router.warnings {
		fmt.Print(warning)
	}

	fmt.Println("Listing on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
