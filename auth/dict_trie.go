package auth

import (
	"strings"
	"sync"
)

var once sync.Once

type dictTrie struct {
	name     string
	children []*dictTrie
	end      bool
}

type Tree struct {
	root *dictTrie
}

func New() *Tree {
	return &Tree{root: &dictTrie{
		name:     "",
		children: make([]*dictTrie, 0),
		end:      false,
	}}
}

func (t *Tree) Insert(path string) {
	node := t.root
	for _, p := range strings.Split(path, "/") {
		if p == "" {
			continue
		}
		children := node.children
		for _, child := range children {
			if child.name == p {
				node = child
				break
			}
		}

	}
}
