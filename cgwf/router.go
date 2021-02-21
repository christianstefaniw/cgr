package cgwf

import (
	"context"
	"net/http"
	"regexp"
	"unicode/utf8"
)

type Router struct{
	routes []RouteEntry
}
type RouteEntry struct {
	Path    *regexp.Regexp
	Method  string
	HandlerFunc http.HandlerFunc
}


func appendSlash(path string) string{
	path += "/"
	return path
}
func removeSlash(path string) string{
	if path[utf8.RuneCountInString(path)-1] == '/' {
		path = path[:utf8.RuneCountInString(path)-1]
	}
	return path
}

func (routeEntry *RouteEntry) match(urlPath string) bool{
	match := routeEntry.Path.FindStringSubmatch(urlPath)
	if match == nil{
		return false
	}
	return true
}

func (routeEntry *RouteEntry) parseUrlParams(r *http.Request) map[string]string{
	groups := routeEntry.Path.FindStringSubmatch(r.URL.Path)
	params := make(map[string]string)
	groupNames := routeEntry.Path.SubexpNames()
	for i, group := range groups{
		params[groupNames[i]] = group
	}
	return params
}

func pathToRegex(path string) *regexp.Regexp{
	var newPath string

	path = appendSlash(path)

	for i := 0; i < utf8.RuneCountInString(path); i++{
		if path[i] == ':'{
			var param string
			for _, paramRune := range path[i+1:]{
				i++
				if paramRune == '/'{
					newPath += `(?P<` + param + `>\w+)/`
					break
				}
				param += string(paramRune)
			}
		} else {
			newPath += string(path[i])
		}
	}

	return regexp.MustCompile("^" + removeSlash(newPath) + "$")
}

func (router *Router) Route(method, path string, handlerFunc http.HandlerFunc){
	exactPath := pathToRegex(path)

	route := RouteEntry{
		Method: method,
		Path: exactPath,
		HandlerFunc: handlerFunc,
	}

	router.routes = append(router.routes, route)
}


func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer internalError(&w)

	for _, route := range router.routes {
		if !route.match(r.URL.Path){
			continue
		}
		params := route.parseUrlParams(r)
		ctx := context.WithValue(r.Context(), "params", params)
		route.HandlerFunc.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	http.NotFound(w, r)
}

func internalError(w *http.ResponseWriter){
	if r := recover(); r != nil{
		http.Error(*w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}

func NewRouter() *Router{
	return &Router{}
}