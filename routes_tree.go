package cgr

import (
	"errors"
	"net/http"
	"unicode/utf8"
)

type tree struct {
	method map[string]*node
}

type node struct {
	route    *Route
	letter   rune
	children map[string]*node
}

func newTree() *tree {
	return &tree{
		method: map[string]*node{
			http.MethodGet: {
				children: make(map[string]*node),
				route:    nil,
				letter:   ' ',
			},
			http.MethodPost: {
				children: make(map[string]*node),
				route:    nil,
				letter:   ' ',
			},
			http.MethodPut: {
				children: make(map[string]*node),
				route:    nil,
				letter:   ' ',
			},
			http.MethodDelete: {
				children: make(map[string]*node),
				route:    nil,
				letter:   ' ',
			},
			http.MethodPatch: {
				children: make(map[string]*node),
				route:    nil,
				letter:   ' ',
			},
		},
	}
}

func (t *tree) insert(r *Route) error {
	methodNode := t.method[r.method]

	// if there is no node belonging to the path's first letter, create it
	if _, ok := methodNode.children[string(r.letter)]; !ok {
		t.initLetterNode(r.letter, methodNode)
	}


	for _, letterNode := range methodNode.children {
		// insert new node under the node belonging to the path's first letter
		if r.letter == letterNode.letter {
			letterNode.children[r.rawPath] = routeNode(r)
		}
	}
	return nil
}

func routeNode(r *Route) *node {
	return &node{
		letter:   ' ',
		route:    r,
		children: nil,
	}
}

func (t *tree) initLetterNode(letter rune, methodNode *node) {
	methodNode.children[string(letter)] = &node{
		route:    nil,
		letter:   letter,
		children: make(map[string]*node),
	}
}

func (t *tree) search(method string, path string) (*Route, error) {
	methodNode := t.method[method]
	var letter rune
	var r *Route

	if len(methodNode.children) == 0 {
		return nil, errors.New("there are no " + method + " routes")
	}

	if path == "/" {
		r = methodNode.children["/"].children["/"].route
		r.params = r.getParams(path)
		return r, nil
	} else {
		letter = rune(path[1])
	}

	if _, ok := methodNode.children[string(letter)]; !ok {

		// check if there are routes without a letter
		if methodNode.children[string(' ')] != nil{
			letter = ' '
		} else {
			return nil, errors.New("path not found")
		}
	}

	for _, n := range methodNode.children[string(letter)].children {
		match := n.route.path.FindStringSubmatch(path)
		if match != nil {
			r = n.route
			r.params = r.getParams(path)
			return r, nil
		} else {
			if n.route.checkAppendSlash(path) {
				r = n.route
				r.params = r.getParams(path)
				return r, nil
			}
		}
	}

	return nil, errors.New("path not found")
}

func (route *Route) checkAppendSlash(path string) bool {
	if route.appendSlash {
		if path[utf8.RuneCountInString(path)-1] != pathDelimiter {
			match := route.path.FindStringSubmatch(path + string(pathDelimiter))
			if match == nil {
				return false
			} else {
				return true
			}
		}
	}
	return false
}

func (router *Router) ViewRouteTree() []string {
	var strTree []string
	for k, n := range router.routes.method {
		for p, o := range n.children {
			for j := range o.children {
				strTree = append(strTree, "method:"+k+" -> letter:"+p+" -> route:"+j+"\n")
			}
		}
	}
	return strTree
}
