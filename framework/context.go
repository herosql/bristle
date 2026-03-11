package framework

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	hasTimeout     bool
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
	}
}

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key any) any {
	return ctx.BaseContext().Value(key)
}

func (ctx *Context) GetRequest() *http.Request       { return ctx.request }
func (ctx *Context) GetWriter() http.ResponseWriter  { return ctx.responseWriter }
func (ctx *Context) WriterMux(w http.ResponseWriter) { ctx.responseWriter = w }
func (ctx *Context) SetHasTimeout(v bool)            { ctx.hasTimeout = v }
func (ctx *Context) HasTimeout() bool                { return ctx.hasTimeout }

func (ctx *Context) QueryString(key string, def string) string {
	val := ctx.request.URL.Query().Get(key)
	if val == "" {
		return def
	}
	return val
}

func (ctx *Context) QueryInt(key string, def int) int {
	valStr := ctx.request.URL.Query().Get(key)
	if valStr == "" {
		return def
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return def
	}
	return val
}

func (ctx *Context) QueryArray(key string, def []string) []string {
	values := ctx.request.URL.Query()[key]
	if len(values) == 0 {
		return def
	}
	return values
}

func (ctx *Context) QueryAll() map[string][]string {
	return ctx.request.URL.Query()
}

func (ctx *Context) FormString(key string, def string) string {
	val := ctx.request.FormValue(key)
	if val == "" {
		return def
	}
	return val
}

func (ctx *Context) FormInt(key string, def int) int {
	valStr := ctx.request.FormValue(key)
	if valStr == "" {
		return def
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return def
	}
	return val
}

func (ctx *Context) FormArray(key string, def []string) []string {
	_ = ctx.request.ParseForm()
	values := ctx.request.Form[key]
	if len(values) == 0 {
		return def
	}
	return values
}

func (ctx *Context) FormAll() map[string][]string {
	_ = ctx.request.ParseForm()
	return ctx.request.Form
}

func (ctx *Context) BindJson(obj any) error {
	defer ctx.request.Body.Close()
	return json.NewDecoder(ctx.request.Body).Decode(obj)
}

func (ctx *Context) Json(status int, obj any) error {
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	return json.NewEncoder(ctx.responseWriter).Encode(obj)
}

func (ctx *Context) Text(status int, text string) error {
	ctx.responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.responseWriter.WriteHeader(status)
	_, err := ctx.responseWriter.Write([]byte(text))
	return err
}

func (ctx *Context) HTML(status int, html string) error {
	ctx.responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.responseWriter.WriteHeader(status)
	_, err := ctx.responseWriter.Write([]byte(html))
	return err
}
