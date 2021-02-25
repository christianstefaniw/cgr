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

// route configurations
type routeConf struct {
	appendSlash bool
	skipClean   bool
}

// set an http protocol to the route
func (route *route) Method(m string) *route {
	route.method = strings.ToUpper(m)
	return route
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
func (router *router) Route(path string) *route {

	router.check(path)

	// TODO select cleaning for singular route
	if !router.skipClean{
		path = cleanPath(path)
	}

	regexPath := pathToRegex(path)

	r := &route{
		path:      regexPath,
		rawPath:   path,
		method:    "GET",
		routeConf: router.routeConf,
	}

	if utf8.RuneCountInString(path) == 1 {
		r.letter = '/'
		err := router.routes.insert(r)
		if err != nil {
			panic(err)
		}

	} else {
		r.letter = rune(path[1])
		err := router.routes.insert(r)
		if err != nil {
			panic(err)
		}
	}

	return r
}

// Returns a pointer to a new route configuration with the default configurations
func NewRouteConf() *routeConf {
	conf := &routeConf{}
	conf.setDefaultRouteConf()
	return conf
}

// Set custom configurations for a route
func (route *route) SetConf(conf *routeConf) *route {
	route.routeConf = *conf
	return route
}

/*
example.com/path is treated the same as example.com/path/

Default is true
*/
func (conf *routeConf) AppendSlash(value bool) *routeConf {
	conf.appendSlash = value
	return conf
}

func GetVars(r *http.Request) map[string]string {
	ctx := r.Context()
	p := ctx.Value("params").(map[string]string)
	return p
}
