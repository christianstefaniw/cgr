package cgr

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"
)

type router struct {

	// Pointers to routes will be stored as values related to their first letter
	routes   *tree
	warnings []string
	routeConf
}

type params map[string]string

// Check if the route.path matches the requested URL Path (r.URL.Path)
func (route *route) match(r *http.Request) (*params, error) {
	match := route.path.FindStringSubmatch(r.URL.Path)
	p := params{}
	if match == nil {
		fmt.Println(string(r.URL.Path[utf8.RuneCountInString(r.URL.Path)-1]))
		if route.appendSlash && r.URL.Path[utf8.RuneCountInString(r.URL.Path)-1] != pathDelimiter {
			match = route.path.FindStringSubmatch(r.URL.Path + string(pathDelimiter))
			if match == nil {
				return &p, nil
			}
		} else {
			return &p, nil
		}
	}
	if r.Method != route.method {
		return &p, errors.New("method is not allowed")
	}

	groupNames := route.path.SubexpNames()
	for i, group := range match {
		p[groupNames[i]] = group
	}

	// params includes the path at empty string value
	return &p, nil
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
