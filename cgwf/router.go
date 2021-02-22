package cgwf

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Router struct {
	routes   []*RouteEntry
	warnings []string
	RouteConf
}
type RouteEntry struct {
	Path        *regexp.Regexp
	Method      string
	HandlerFunc http.HandlerFunc
	RouteConf
}
type RouteConf struct {
	AppendSlash  bool
	CheckPattern bool
}

func (routeEntry *RouteEntry) match(r *http.Request) map[string]string {
	match := routeEntry.Path.FindStringSubmatch(r.URL.Path)
	if match == nil {
		if routeEntry.AppendSlash && r.URL.Path[utf8.RuneCountInString(r.URL.Path)-1] != '/' {
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

func (router *Router) check(path string) {
	var warning string
	if strings.Contains("(?P<", path) ||
		strings.Index(path, "^") == 0 ||
		strings.Index(path, "$") == utf8.RuneCountInString(path) {
		warning =
			"!!WARNING!!\n" +
				`Your url pattern ` + path +
				` has a route that contains '(?P<', begins with a '^', or ends with a '$'.`
	}
	router.warnings = append(router.warnings, warning)
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

func (router *Router) Route(method, path string, handlerFunc http.HandlerFunc) *RouteEntry{
	if router.CheckPattern{
		router.check(path)
	}
	exactPath := pathToRegex(path)

	route := RouteEntry{
		Method:      method,
		Path:        exactPath,
		HandlerFunc: handlerFunc,
		RouteConf:   router.RouteConf,
	}

	router.routes = append(router.routes, &route)
	return &route
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
