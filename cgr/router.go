package cgr

import (
	"strings"
	"unicode/utf8"
)

type router struct {

	// TODO improve how routes are stored
	routes   []*route
	warnings []string
	routeConf
}

// Check if the route.path matches the requested URL Path (r.URL.Path)
func (route *route) match(path string) (bool, map[string]string) {
	match := route.path.FindStringSubmatch(path)
	params := make(map[string]string)
	if match == nil {
		if route.appendSlash && path[utf8.RuneCountInString(path)-1] != '/' {
			match = route.path.FindStringSubmatch(path + "/")
			if match == nil {
				return false, params
			}
		} else {
			return false, params
		}
	}

	groupNames := route.path.SubexpNames()
	for i, group := range match {
		params[groupNames[i]] = group
	}


	// params includes the path at empty string value
	return true, params
}

// Check for bad patterns
func (router *router) check(path string) {
	var warning string
	if strings.Contains("(?P<", path) ||
		strings.Index(path, "^") == 0 ||
		strings.Index(path, "$") == utf8.RuneCountInString(path) {
		warning =
			"!!WARNING!!\n" +
				"Your url pattern " + path +
				" has a route that contains '(?P<', begins with a '^', or ends with a '$'. \n \n"
	}
	router.warnings = append(router.warnings, warning)
}

// Returns a pointer to a new router with the default route configurations
func NewRouter() *router {
	r := &router{}
	r.setDefaultRouteConf()
	return r
}

func (conf *routeConf) setDefaultRouteConf() {
	conf.appendSlash = true
	conf.skipClean = false
}
