package goKLC

import (
	"fmt"
	"net/http"
)

type App struct {
}

var routeTree *Node
var response string

func NewApp() *App {
	routeTree = NewRouteTree()

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

func (a *App) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	route, ok, params := match(req)

	if !ok {

		fmt.Fprintf(rw, "404", nil)
		return
	}

	response = route.controller(req, params)

	fmt.Fprintf(rw, response)
}
