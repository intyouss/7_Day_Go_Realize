package gee

import (
	"strings"
)

type node struct {
	pattern  string
	part     string //ps: :lang
	children []*node
	isFuzzy  bool
}

func (n *node) matchChild(part string) (*node, bool) {
	for _, child := range n.children {
		if part == child.part {
			return child, true
		}
	}
	return nil, false
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if part == child.part || child.isFuzzy {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]

	child, ok := n.matchChild(part)
	if !ok {
		child = &node{part: part, isFuzzy: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
