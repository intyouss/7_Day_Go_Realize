package gee

import (
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (e *Engine) addRoute(method, path string, handler HandlerFunc) {
	key := method + "-" + path
	e.router[key] = handler
}

func (e *Engine) Run(port string) (err error) {
	err = http.ListenAndServe(port, e)
	return
}

func (e *Engine) GET(path string, handler HandlerFunc) {
	e.addRoute("GET", path, handler)
}

func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}
