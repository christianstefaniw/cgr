package cgr

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Route struct {
	path        *regexp.Regexp
	rawPath     string
	handlerFunc http.HandlerFunc
	letter      rune
	params      *params
	method      string
	router      *Router
	middleware  *middlewareLinkedList
	routeConf
}

// assign middleware to the route
func (route *Route) Assign(middleware *middleware) *Route {
	route.middleware.insert(middleware)
	return route
}

// set an http protocol for the Route
func (route *Route) Method(m string) *Route {
	route.method = strings.ToUpper(m)
	return route
}

// insert Route into tree
func (route *Route) Insert() {
	if !route.skipClean {
		route.rawPath = cleanPath(route.rawPath)
	}

	route.path = pathToRegex(route.rawPath)

	if utf8.RuneCountInString(route.rawPath) == 1 {
		route.letter = '/'
		err := route.router.routes.insert(route)
		if err != nil {
			panic(err)
		}

	} else if route.rawPath[1] == ':' {
		route.letter = ' '
		err := route.router.routes.insert(route)
		if err != nil {
			panic(err)
		}
	} else {
		route.letter = rune(route.rawPath[1])
		err := route.router.routes.insert(route)
		if err != nil {
			panic(err)
		}
	}
}

// attach a handler function to the Route
func (route *Route) Handler(handler http.HandlerFunc) *Route {
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

// Create a new Route entry
func (router *Router) Route(path string) *Route {

	router.check(path)

	r := &Route{
		rawPath:    path,
		routeConf:  router.routeConf,
		router:     router,
		middleware: new(middlewareLinkedList),
	}

	return r
}

// Returns a pointer to a new Route configuration with the default configurations
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

func (route *Route) executeMiddleware(w http.ResponseWriter, r *http.Request) {
	currNode := route.middleware.head
	for currNode != nil {
		currNode.mware.run(w, r)
		currNode = currNode.next
	}
	route.handlerFunc.ServeHTTP(w, r)
}
