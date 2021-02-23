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
	handlerFunc http.HandlerFunc
	err         error
	method      string
	routeConf
}

type routeConf struct {
	appendSlash bool
	skipClean bool
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
		if path[i] == ':' {
			var param string
			nearestSlash := strings.IndexRune(path[i+1:], '/')
			nearestColon := strings.IndexRune(path[i+1:], ':')
			var section string

			if nearestSlash == -1 {
				section = path[i+1:]
			} else {
				section = path[i+1 : nearestSlash+i+1]
			}

			if (nearestSlash > nearestColon && nearestColon != -1) || (nearestSlash == -1 && nearestColon != -1){
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

	r := route{
		path:      regexPath,
		routeConf: router.routeConf,
	}

	if utf8.RuneCountInString(path) == 1{
		router.routes['.'] = append(router.routes['.'], &r)
	} else {
		router.routes[rune(path[1])] = append(router.routes[rune(path[1])], &r)
	}


	return &r
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
func (conf *routeConf) AppendSlash(value bool){
	conf.appendSlash = value
}

/*
Remove . and .. from url path

Default is false
 */
func (conf *routeConf) SkipClean(value bool){
	conf.skipClean = value
}


func GetVars(r *http.Request) map[string]string {
	ctx := r.Context()
	p := ctx.Value("params").(map[string]string)
	return p
}
