package goKLC

import (
	"fmt"
	"net/http"
	"sync"
)

type App struct {
}

var routeTree *RouteNode
var middlewareList *MiddlewareNode
var mux = &sync.RWMutex{}

func NewApp() *App {
	routeTree = NewRouteTree()
	middlewareList = NewMiddlewareNode()

	return &App{}
}

func (a *App) Run() {

	err := http.ListenAndServe(":8093", a)

	fmt.Println(err)

}

func (a *App) Route() Route {
	rg := NewRouteGroup()
	rg.prefix = ""

	return rg.Route()
}

func (a *App) Middleware(m *Middleware) {
	if middlewareList == nil {
		middlewareList.middleware = m
	} else {
		mn := NewMiddlewareNode()
		mn.middleware = m

		middlewareList.AddChild(mn)
	}
}

func (a *App) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	mux.Lock()

	request := NewRequest(req)
	route, ok, params := match(req)

	if !ok {

		fmt.Fprintf(rw, "404", nil)
		return
	}

	middlewareList.Handle(request)
	response := route.controller(request, params)

	rw.WriteHeader(response.status)
	rw.Write([]byte(response.content))

	mux.Unlock()
}