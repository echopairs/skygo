package router

import (
	"github.com/echopairs/skygo/gweb/web/common"
	"github.com/julienschmidt/httprouter"
	_ "github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
)

type Route struct {
	Name 	string
	Method 	string
	Path 	string

}

type Router struct {
	*httprouter.Router
	routes []*Route
	sync.RWMutex
}

func NewRouter() *Router {
	 hr := httprouter.New()
	router := &Router{
		Router: hr,
	}
	return router
}

var (
	defaultRouter = NewRouter()
)

func RegisterHttpRouteHandle(method, path, funcname string, handler httprouter.Handle) {
	defaultRouter.RegisterHttpRouteHandle(method, path, funcname, handler)
}

func RegisterHttpHandle(method, path, funcname string, handle http.Handler) {
	defaultRouter.RegisterHttpHandle(method, path, funcname, handle)
}

func RegisterHttpHandleFunc(method, path, funcname string, handlerFunc http.HandlerFunc) {
	defaultRouter.RegisterHttpHandleFunc(method, path, funcname, handlerFunc)
}

func (router *Router) RegisterHttpRouteHandle(method, path, funcname string, handler httprouter.Handle) {
	router.Lock()
	defer router.Unlock()
	router.Handle(method, path, handler)
	router.addRoute(method, path, funcname)

}

func (router *Router) RegisterHttpHandle(method, path, funcname string, handle http.Handler) {
	router.Lock()
	defer router.Unlock()
	router.Handler(method, path, handle)
	router.addRoute(method, path, funcname)
}

func (router *Router) RegisterHttpHandleFunc(method, path, funcname string, handlerFunc http.HandlerFunc) {
	router.Lock()
	defer router.Unlock()
	router.HandlerFunc(method, path, handlerFunc)
	router.addRoute(method, path, funcname)
}

func (router *Router) addRoute(method, path, funcName string) {
	hrouter := &Route{
		funcName,
		method,
		path,
	}
	router.routes = append(router.routes, hrouter)
}

func (router *Router) GetAllRoutes() []*Route {
	router.Lock()
	defer router.Unlock()
	return router.routes
}

func GetDefaultRouter() *Router {
	return defaultRouter
}

// Handler for get all route
func RouteIndex(w http.ResponseWriter, r *http.Request) {
	route := defaultRouter.GetAllRoutes()
	res := &common.ResBody{
		Err:0,
		Data:route,
	}
	common.WriteJson(w, res, http.StatusOK)
}
