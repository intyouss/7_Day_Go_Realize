package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	engine      *Engine
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func Default() *Engine {
	e := New()
	e.Use(Logger(), Recovery())
	return e
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:      prefix,
		middlewares: make([]HandlerFunc, 0),
		engine:      engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRouter(method, path string, handler HandlerFunc) {
	pattern := group.prefix + path
	log.Printf("Router %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
}

func (group *RouterGroup) createStaticHandler(simplePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, simplePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Fail(http.StatusNotFound, err)
			return
		}
		c.SetStatus(http.StatusOK)
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

func (group *RouterGroup) Static(simplePath string, root string) {
	handler := group.createStaticHandler(simplePath, http.Dir(root))
	urlPattern := path.Join(simplePath, "/*filepath")
	group.GET(urlPattern, handler)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) SetTemplatePath(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

func (e *Engine) Run(addr string) (err error) {
	err = http.ListenAndServe(addr, e)
	return err
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, r)
	c.middlewares = middlewares
	c.engine = e
	e.router.handler(c)
}
