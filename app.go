package goKLC

import (
	"fmt"
	"net/http"
	"sync"
)

type App struct {
}

var app *App
var routeTree *RouteNode
var middlewareList *MiddlewareNode
var mux = &sync.RWMutex{}

func GetApp() *App {

	if app == nil {
		routeTree = NewRouteTree()
		middlewareList = NewMiddlewareNode()

		app = &App{}
	}

	return app
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

func (a *App) Middleware(m MiddlewareInterface) {
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
	route, ok, params := match(req)
	request := NewRequest(req, params)

	if !ok {

		fmt.Fprintf(rw, "404", nil)
		return
	}

	response, middleware := middlewareList.Handle(request)

	if response == nil {
		response = route.controller(request)
	}

	middleware.Terminate(response)

	rw.WriteHeader(response.status)
	rw.Write([]byte(response.content))

	request = nil
	response = nil
	mux.Unlock()
}
