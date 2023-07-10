package gee

import "net/http"

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func (e *Engine) ServeHttp() {}

func (e *Engine) Run(port string) {

}

func (e *Engine) GET(path string, handler HandlerFunc) {

}

func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}
