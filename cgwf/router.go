package cgwf

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
)

var AppendSlash bool

type Router struct {
	routes []RouteEntry
}
type RouteEntry struct {
	Path        *regexp.Regexp
	Method      string
	HandlerFunc http.HandlerFunc
}

func (routeEntry *RouteEntry) match(r *http.Request) map[string]string {
	match := routeEntry.Path.FindStringSubmatch(r.URL.Path)

	if match == nil {
		if AppendSlash {
			match = routeEntry.Path.FindStringSubmatch(r.URL.Path + "/")
			if match == nil {
				return nil
			}
		} else {
			return nil
		}
	}

	params := make(map[string]string)
	groupNames := routeEntry.Path.SubexpNames()
	for i, group := range match {
		params[groupNames[i]] = group
	}
	return params
}

func pathToRegex(path string) *regexp.Regexp {
	var newPath string

	for i := 0; i < utf8.RuneCountInString(path); i++ {
		if path[i] == ':' {
			var param string
			nearestSlash := strings.IndexRune(path[i+1:], '/')
			var section string
			if nearestSlash == -1 {
				section = path[i+1:]
			} else {
				section = path[i+1 : nearestSlash+i+1]
			}

			for k, paramRune := range section {
				i++
				param += string(paramRune)
				if k == (utf8.RuneCountInString(section) - 1) {
					newPath += `(?P<` + param + `>\w+)`

					break
				}

			}
		} else {
			newPath += string(path[i])
		}
	}
	return regexp.MustCompile("^" + newPath + "$")
}

func (router *Router) Route(method, path string, handlerFunc http.HandlerFunc) {
	exactPath := pathToRegex(path)

	route := RouteEntry{
		Method:      method,
		Path:        exactPath,
		HandlerFunc: handlerFunc,
	}

	router.routes = append(router.routes, route)
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer internalError(&w)

	for _, route := range router.routes {
		params := route.match(r)
		if params == nil {
			// match not found
			continue
		}
		ctx := context.WithValue(r.Context(), "params", params)
		route.HandlerFunc.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	http.NotFound(w, r)
}

func internalError(w *http.ResponseWriter) {
	if r := recover(); r != nil {
		http.Error(*w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

func NewRouter() *Router {
	return &Router{}
}
