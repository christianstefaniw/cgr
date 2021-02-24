package cgr

import (
	"fmt"
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
	err         error
	method      string
	routeConf
}

type routeConf struct {
	appendSlash bool
	skipClean   bool
}

func (route *route) Method(m string) *route {
	route.method = strings.ToUpper(m)
	return route
}

func (route *route) Handler(handler http.HandlerFunc) *route {
	route.handlerFunc = handler
	return route
}

// algorithm to parse a path to a regular expression
func pathToRegex(path string) *regexp.Regexp {
	var newPath string

	for i := 0; i < utf8.RuneCountInString(path); i++ {
		if path[i] == paramDelimiter {
			var param string
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

// Create a new route entry
func (router *router) Route(path string) *route {

	router.check(path)

	regexPath := pathToRegex(path)

	r := &route{
		path:      regexPath,
		rawPath:   path,
		method:    "GET",
		routeConf: router.routeConf,
	}

	if utf8.RuneCountInString(path) == 1 {
		router.routes['.'] = append(router.routes['.'], r)
		r.letter = '.'
	} else {
		router.routes[rune(path[1])] = append(router.routes[rune(path[1])], r)
		r.letter = rune(path[1])
	}

	t := newTree()
	_ = t.insert(r)
	fmt.Println(t.search(r.method, r.rawPath))

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

/*
Remove . and .. from url path

Default is false
*/
func (conf *routeConf) SkipClean(value bool) *routeConf {
	conf.skipClean = value
	return conf
}

func GetVars(r *http.Request) map[string]string {
	ctx := r.Context()
	p := ctx.Value("params").(map[string]string)
	return p
}
