package cgr

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// ServeHTTP dispatches the handler registered in the matched route.
func (router *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer internalError(&w)

	for _, r := range router.routes {
		var params map[string]string
		var found bool
		if router.skipClean{
			found, params = r.match(req.URL.Path)
		} else {
			found, params = r.match(cleanPath(req.URL.Path))
		}

		if !found {
			// match not found
			continue
		}
		ctx := context.WithValue(req.Context(), "params", params)
		r.handlerFunc.ServeHTTP(w, req.WithContext(ctx))
		return
	}

	http.NotFound(w, req)
}


func internalError(w *http.ResponseWriter) {
	if r := recover(); r != nil {
		http.Error(*w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
func Run(port string, router *router) {

	for _, warning := range router.warnings {
		fmt.Print(warning)
	}

	fmt.Println("Listing on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
