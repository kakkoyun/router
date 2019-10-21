package route

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Router wraps httprouter.Router and adds support for prefixed sub-routers,
// per-request context injections and instrumentation.
type Router struct {
	rtr    *httprouter.Router
	instrh func(handlerName string, handler http.Handler) http.Handler
}

// New returns a new Router.
func New() *Router {
	return &Router{
		rtr: httprouter.New(),
		instrh: func(_ string, h http.Handler) http.Handler {
			return h
		},
	}
}

// WithInstrumentation returns a router with instrumentation support.
func (r *Router) WithInstrumentation(instrh func(handlerName string, handler http.Handler) http.Handler) *Router {
	return &Router{rtr: r.rtr, instrh: instrh}
}

// Get registers a new GET route.
func (r *Router) Get(path string, h http.Handler) {
	r.rtr.Handler(http.MethodGet, path, r.instrh(path, h))
}

// Options registers a new OPTIONS route.
func (r *Router) Options(path string, h http.Handler) {
	r.rtr.Handler(http.MethodOptions, path, r.instrh(path, h))
}

// Del registers a new DELETE route.
func (r *Router) Delete(path string, h http.Handler) {
	r.rtr.Handler(http.MethodDelete, path, r.instrh(path, h))
}

// Put registers a new PUT route.
func (r *Router) Put(path string, h http.Handler) {
	r.rtr.Handler(http.MethodPut, path, r.instrh(path, h))
}

// Post registers a new POST route.
func (r *Router) Post(path string, h http.Handler) {
	r.rtr.Handler(http.MethodPost, path, r.instrh(path, h))
}

// Redirect takes an absolute path and sends an internal HTTP redirect for it,
// prefixed by the router's path prefix. Note that this method does not include
// functionality for handling relative paths or full URL redirects.
func (r *Router) Redirect(w http.ResponseWriter, req *http.Request, path string, code int) {
	http.Redirect(w, req, path, code)
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.rtr.ServeHTTP(w, req)
}
