package cgr

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
)

type route struct {
	path        *regexp.Regexp
	rawPath     string
	handlerFunc http.HandlerFunc
	letter      rune
	params      *params
	method      string
	routeConf
}


// set an http protocol for the route
func (route *route) Method(m string) *route {
	route.method = strings.ToUpper(m)
	return route
}


// insert route into tree
func (route *route) Insert(router *Router) {
	if !route.skipClean {
		route.rawPath = cleanPath(route.rawPath)
	}

	route.path = pathToRegex(route.rawPath)

	if utf8.RuneCountInString(route.rawPath) == 1 {
		route.letter = '/'
		err := router.routes.insert(route)
		if err != nil {
			panic(err)
		}

	} else if route.rawPath[1] == ':'{
		route.letter = rune(route.rawPath[2])
		err := router.routes.insert(route)
		if err != nil {
			panic(err)
		}
	} else {
		route.letter = rune(route.rawPath[1])
		err := router.routes.insert(route)
		if err != nil {
			panic(err)
		}
	}
}

// attach a handler function to the route
func (route *route) Handler(handler http.HandlerFunc) *route {
	route.handlerFunc = handler
	return route
}

// algorithm to parse a path to a regular expression
func pathToRegex(path string) *regexp.Regexp {
	var newPath string

	for i := 0; i < utf8.RuneCountInString(path); i++ {
		if path[i] == paramDelimiter {
			var p string
			nearestSlash := strings.IndexRune(path[i+1:], pathDelimiter)
			nearestParam := strings.IndexRune(path[i+1:], paramDelimiter)
			var section string

			if nearestSlash == -1 {
				section = path[i+1:]
			} else {
				section = path[i+1 : nearestSlash+i+1]
			}

			if (nearestSlash > nearestParam && nearestParam != -1) || (nearestSlash == -1 && nearestParam != -1) {
				log.Fatal(`"/" must come before ":"`)
			}

			for k, paramRune := range section {
				i++
				p += string(paramRune)
				if k == (utf8.RuneCountInString(section) - 1) {
					newPath += `(?P<` + p + `>\w+)`
					break
				}

			}
		} else {
			newPath += string(path[i])
		}
	}
	return regexp.MustCompile("^" + newPath + "$")
}

// Create a new route entry
func (router *Router) Route(path string) *route {

	router.check(path)

	r := &route{
		rawPath:   path,
		routeConf: router.routeConf,
	}

	return r
}

// Returns a pointer to a new route configuration with the default configurations
func NewRouteConf() *routeConf {
	conf := &routeConf{}
	conf.setDefaultRouteConf()
	return conf
}

// Get url parameters
func GetParams(r *http.Request) map[string]string {
	ctx := r.Context()
	p := ctx.Value("params").(map[string]string)
	return p
}
