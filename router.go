package cgr

import (
	"strings"
	"unicode/utf8"
)

type router struct {
	routes   *tree
	warnings []string
	routeConf
}


type params map[string]string

func (route *route) params(path string) *params {
	var match []string
	match = route.path.FindStringSubmatch(path)

	if match == nil{
		match = route.path.FindStringSubmatch(appendSlash(path))
	}

	p := make(params)

	groupNames := route.path.SubexpNames()
	for i, group := range match {
		p[groupNames[i]] = group
	}

	return &p
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
	if path[0] != '/'{
		warning += "!!WARNING!! \n" +
			"Url pattern " + path + " must to start with a / \n \n"
	}
	router.warnings = append(router.warnings, warning)
}

// Returns a pointer to a new router with the default route configurations
func NewRouter() *router {
	r := &router{}
	r.setDefaultRouteConf()
	r.routes = newTree()
	return r
}

func (conf *routeConf) setDefaultRouteConf() {
	conf.appendSlash = true
	conf.skipClean = false
}
