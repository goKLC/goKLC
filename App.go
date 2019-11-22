package goKLC

import (
	"fmt"
	"net/http"
)

type App struct {
}

var routeTree *RouteNode
var middlewareList *MiddlewareNode
var response string

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
	middleware := middlewareList.Handle()
	route, ok, params := match(req)

	if !ok {

		fmt.Fprintf(rw, "404", nil)
		return
	}

	response = route.controller(req, params)
	middleware.Terminate()

	fmt.Fprintf(rw, response)
}
