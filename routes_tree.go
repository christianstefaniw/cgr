package cgr

import (
	"errors"
	"net/http"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

type tree struct {
	method map[string]*node
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

	if _, ok := methodNode.children[string(r.letter)]; !ok{
		t.initLetterBranch(r.letter, methodNode)
	}

	for _, letterNode := range methodNode.children {
		if r.letter == letterNode.letter {
			letterNode.children[r.rawPath] = routeNode(r)
		}
	}
	return nil
}

func routeNode(r *route) *node{
	return &node{
		letter:   ' ',
		route:    r,
		children: nil,
	}
}

func (t *tree) initLetterBranch(letter rune, methodNode *node){
	methodNode.children[string(letter)] = &node{
		route:    nil,
		letter: letter,
		children: make(map[string]*node),
	}
}

func (t *tree) search(method string, path string) (*route, error) {
	methodNode := t.method[method]
	var letter rune

	if len(methodNode.children) == 0 {
		return nil, errors.New("there are no " + method + " routes")
	}

	if path == "/" {
		return methodNode.children["/"].children["/"].route, nil
	} else {
		letter = rune(path[1])
	}


	for _, n := range methodNode.children[string(alphabet[t.binarySearchLetterNodePos(uint8(letter))])].children {
		match := n.route.path.FindStringSubmatch(path)
		if match != nil {
			return n.route, nil
		}
	}

	return nil, errors.New("path not found")
}


func (t *tree) binarySearchLetterNodePos(letter uint8) int {
	start := 0
	end := len(alphabet)

	for start <= end {

		midIndex := (int(start) + int(end))/2
		midLetter := alphabet[midIndex]

		if midLetter == letter {
			return midIndex
		} else if midLetter < letter {
			start = midIndex + 1
		} else {
			end = midIndex - 1
		}
	}
	return -1
}
