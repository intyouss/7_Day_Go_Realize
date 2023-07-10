package gee

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{
		router: &router{
			handlers: make(map[string]HandlerFunc),
		},
	}
}

func (e *Engine) addRoute(method, path string, handler HandlerFunc) {
	e.router.addRoute(method, path, handler)
}

func (e *Engine) Run(addr string) (err error) {
	err = http.ListenAndServe(addr, e)
	return
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.addRoute("GET", path, handler)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.router.handler(c)
}
