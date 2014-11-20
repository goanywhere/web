/**
 * ----------------------------------------------------------------------
 *       ______      ___                         __
 *      / ____/___  /   |  ____  __  ___      __/ /_  ___  ________
 *     / / __/ __ \/ /| | / __ \/ / / / | /| / / __ \/ _ \/ ___/ _ \
 *    / /_/ / /_/ / ___ |/ / / / /_/ /| |/ |/ / / / /  __/ /  /  __/
 *    \____/\____/_/  |_/_/ /_/\__. / |__/|__/_/ /_/\___/_/   \___/
 *                            /____/
 * ----------------------------------------------------------------------
 *  Copyright © 2014 GoAnywhere Ltd. All Rights Reserved.
 * ----------------------------------------------------------------------*/

package web

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/goanywhere/web/env"
	"github.com/gorilla/mux"
)

type (
	Application struct {
		router      *mux.Router
		middlewares []Middleware
	}

	HandlerFunc func(*Context)

	// Conventional method to implement custom middlewares.
	Middleware func(http.Handler) http.Handler

	// Shortcut to create map.
	H map[string]interface{}
)

// New creates an application instance & setup its default settings..
func New() *Application {
	app := &Application{mux.NewRouter(), nil}
	return app
}

// ---------------------------------------------------------------------------
//  Custom handler func with Context Supports
// ---------------------------------------------------------------------------
func (self HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self(NewContext(w, r))
}

// ---------------------------------------------------------------------------
//  HTTP Requests Handlers
// ---------------------------------------------------------------------------
// Supported Handler Types
//	* http.Handler
//	* http.HandlerFunc	=> func(w http.ResponseWriter, r *http.Request)
//	* web.HandlerFunc	=> func(ctx *Context)
func (self *Application) handle(method, pattern string, h interface{}) {
	var handler http.Handler

	switch h.(type) {
	// Standard net/http.Handler/HandlerFunc
	case http.Handler:
		handler = h.(http.Handler)
	case func(w http.ResponseWriter, r *http.Request):
		handler = http.HandlerFunc(h.(func(w http.ResponseWriter, r *http.Request)))
	case func(ctx *Context):
		handler = HandlerFunc(h.(func(ctx *Context)))
	default:
		panic(fmt.Sprintf("Unknown handler type (%v) passed in.", handler))
	}
	// finds the full function name (with package)
	name := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	self.router.Handle(pattern, handler).Methods(method).Name(name)
}

// Address fetches the address predefined in `os.environ` by combineing
// `os.Getenv("host")` & os.Getenv("port").
func (self *Application) Address() string {
	return fmt.Sprintf("%s:%s", env.Get("host"), env.Get("port"))
}

// GET is a shortcut for app.HandleFunc(pattern, handler).Methods("GET"),
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Get(pattern string, handler interface{}) {
	self.handle("GET", pattern, handler)
}

// POST is a shortcut for app.HandleFunc(pattern, handler).Methods("POST")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Post(pattern string, handler interface{}) {
	self.handle("POST", pattern, handler)
}

// PUT is a shortcut for app.HandleFunc(pattern, handler).Methods("PUT")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Put(pattern string, handler interface{}) {
	self.handle("PUT", pattern, handler)
}

// DELETE is a shortcut for app.HandleFunc(pattern, handler).Methods("DELETE")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Delete(pattern string, handler interface{}) {
	self.handle("DELETE", pattern, handler)
}

// PATCH is a shortcut for app.HandleFunc(pattern, handler).Methods("PATCH")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Patch(pattern string, handler http.HandlerFunc) {
	self.handle("PATCH", pattern, handler)
}

// HEAD is a shortcut for app.HandleFunc(pattern, handler).Methods("HEAD")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Head(pattern string, handler http.HandlerFunc) {
	self.handle("HEAD", pattern, handler)
}

// OPTIONS is a shortcut for app.HandleFunc(pattern, handler).Methods("OPTIONS")
// it also fetch the full function name of the handler (with package) to name the route.
func (self *Application) Options(pattern string, handler http.HandlerFunc) {
	self.handle("OPTIONS", pattern, handler)
}

// Group creates a new application group under the given path.
func (self *Application) Group(path string) *Application {
	return &Application{self.router.PathPrefix(path).Subrouter(), nil}
}

// ---------------------------------------------------------------------------
//  HTTP Server with Middleware Supports
// ---------------------------------------------------------------------------
func (self *Application) Use(middlewares ...Middleware) {
	self.middlewares = append(self.middlewares, middlewares...)
}

// ServeHTTP turn Application into http.Handler by implementing the http.Handler interface.
func (self *Application) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var app http.Handler = self.router
	// Activate middlewares in FIFO order.
	if len(self.middlewares) > 0 {
		for index := len(self.middlewares) - 1; index >= 0; index-- {
			app = self.middlewares[index](app)
		}
	}
	app.ServeHTTP(writer, request)
}

// Serve starts serving the requests at the pre-defined address from settings.
// TODO command line arguments.
func (self *Application) Serve() {
	Info("Application server started [%s]", self.Address())
	if err := http.ListenAndServe(self.Address(), self); err != nil {
		panic(err)
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Application Defaults
	env.Set("debug", "true")
	env.Set("host", "0.0.0.0")
	env.Set("port", "5000")
	env.Set("templates", "templates")
}
