package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func ParsePattern(pattern string) []string {
	r := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range r {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRouter(method, pattern string, handler HandlerFunc) {
	parts := ParsePattern(pattern)
	key := method + "-" + pattern
	root, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	root.insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRouter(method string, pattern string) (*node, map[string]string) {
	searchParts := ParsePattern(pattern)
	params := make(map[string]string, 0)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	node := root.search(searchParts, 0)
	if node != nil {
		parts := ParsePattern(node.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(parts) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
}

func (r *router) handler(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.Writer.WriteHeader(http.StatusNotFound)
	}
}
