package cgr

import (
	"errors"
	"net/http"
	"unicode/utf8"
)

type tree struct {
	method map[string]*node
}

type result struct {
	*params
	*route
}

type node struct {
	route    *route
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
		},
	}
}

func (t *tree) insert(r *route) error {
	methodNode := t.method[r.method]

	if _, ok := methodNode.children[string(r.letter)]; !ok {
		t.initLetterBranch(r.letter, methodNode)
	}

	for _, letterNode := range methodNode.children {
		if r.letter == letterNode.letter {
			letterNode.children[r.rawPath] = routeNode(r)
		}
	}
	return nil
}

func routeNode(r *route) *node {
	return &node{
		letter:   ' ',
		route:    r,
		children: nil,
	}
}

func (t *tree) initLetterBranch(letter rune, methodNode *node) {
	methodNode.children[string(letter)] = &node{
		route:    nil,
		letter:   letter,
		children: make(map[string]*node),
	}
}

func (t *tree) search(method string, path string) (*result, error) {
	methodNode := t.method[method]
	var letter rune

	if len(methodNode.children) == 0 {
		return nil, errors.New("there are no " + method + " routes")
	}

	if path == "/" {
		r := methodNode.children["/"].children["/"].route
		return &result{params: r.checkClean(path), route: r}, nil
	} else {
		letter = rune(path[1])
	}

	for _, n := range methodNode.children[string(letter)].children {
		match := n.route.path.FindStringSubmatch(path)
		if match != nil {
			return &result{params: n.route.checkClean(path), route: n.route}, nil
		} else {
			if n.route.appendSlash {
				if n.route.checkAppendSlash(path) {
					return &result{params: n.route.checkClean(path), route: n.route}, nil
				}
			}
		}
	}

	return nil, errors.New("path not found")
}

func (route *route) checkClean(path string) *params{
	if route.skipClean{
		return route.params(path)
	} else {
		return route.params(cleanPath(path))
	}
}


func (route *route) checkAppendSlash(path string) bool {
	if path[utf8.RuneCountInString(path)-1] != pathDelimiter {
		match := route.path.FindStringSubmatch(path + string(pathDelimiter))
		if match == nil {
			return false
		}
	}
	return true
}


func (router *router) ViewRouteTree() []string{
	var strTree []string
	for k, n := range router.routes.method{
		for p, o := range n.children{
			for j := range o.children{
				strTree = append(strTree, "method:" + k + " -> letter:" + p + " -> route:" + j + "\n")
			}
		}
	}
	return strTree
}