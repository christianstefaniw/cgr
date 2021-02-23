package cgr

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"unicode/utf8"
)

// ServeHTTP dispatches the handler registered in the matched route.
func (router *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	defer internalError(&w)
	var letter rune



	if utf8.RuneCountInString(req.URL.Path) == 1{
		letter = '.'
	} else {
		letter = rune(req.URL.Path[1])
	}


	for _, r := range router.routes[letter] {
		var p *params
		var found bool
		var err error
		if router.skipClean {
			found, p, err = r.match(req)
		} else {
			req.URL.Path = cleanPath(req.URL.Path)
			found, p, err = r.match(req)
		}
		if !found {
			// match not found
			continue
		}
		if err != nil {
			methodNotAllowed(&w)
			return
		}

		paramsAsMap := paramsToMap(p)

		ctx := context.WithValue(req.Context(), "params", paramsAsMap)

		r.handlerFunc.ServeHTTP(w, req.WithContext(ctx))
		return
	}

	http.NotFound(w, req)
}

func paramsToMap(p *params) map[string]string{
	paramsAsMap := make(map[string]string)
	for i, k := range *p {
		paramsAsMap[i] = k
	}
	return paramsAsMap
}

func internalError(w *http.ResponseWriter) {
	if r := recover(); r != nil {
		http.Error(*w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

func methodNotAllowed(w *http.ResponseWriter) {
	http.Error(*w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
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
