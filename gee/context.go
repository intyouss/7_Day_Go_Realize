package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]any

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Params map[string]string
	Path   string
	Method string
	// response info
	StatusCode int
	// middleware
	middlewares []HandlerFunc
	index       int
	engine      *Engine
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.middlewares)
	for ; c.index < s; c.index++ {
		c.middlewares[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.middlewares)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	value := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) SetStatus(statusCode int) {
	c.StatusCode = statusCode
	c.Writer.WriteHeader(statusCode)
}

func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) JSON(statusCode int, obj any) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(statusCode)
	jsonData := json.NewEncoder(c.Writer)
	if err := jsonData.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) String(statusCode int, format string, values ...any) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(statusCode)
	_, err := c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) Data(statusCode int, data []byte) {
	c.SetStatus(statusCode)
	_, err := c.Writer.Write(data)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) HTML(statusCode int, name string, data any) {
	c.SetHeader("Content-Type", "text/html")
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(http.StatusInternalServerError, err.Error())
	} else {
		c.SetStatus(statusCode)
	}
}
