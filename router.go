package cgr

import (
	"strings"
	"unicode/utf8"
)

type Router struct {
	routes   *tree
	warnings []string
	routeConf
}

type param struct {
	key string
	value string
}

type params []*param

func (route *route) getParams(path string) *params {
	var match []string
	var p params

	match = route.path.FindStringSubmatch(path)

	if match == nil{
		match = route.path.FindStringSubmatch(appendSlash(path))
	}

	groupNames := route.path.SubexpNames()
	for i, group := range match {
		p = append(p, &param{key: groupNames[i], value: group})
	}

	return &p
}

// Check for bad patterns
func (router *Router) check(path string) {
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

// insert route from Router instance
func (router *Router) Insert(route *route){
	route.Insert(router)
}


// returns a pointer to a new Router with the default route configurations
func NewRouter() *Router {
	r := &Router{}
	r.setDefaultRouteConf()
	r.routes = newTree()
	return r
}
