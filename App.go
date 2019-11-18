package goKLC

import (
	"fmt"
	"github.com/emirpasic/gods/maps/hashmap"
	"net/http"
)

type App struct {
}

var rb *hashmap.Map
var response string

func NewApp() *App {
	rb = hashmap.New()

	return &App{}
}

func (a *App) Run() {

	fmt.Println(rb)
	err := http.ListenAndServe(":8093", a)

	fmt.Println(err)

}

func (a *App) Route() Route {
	rg := NewRouteGroup()
	rg.prefix = ""

	return rg.Route()
}

func (a *App) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r, ok := match(rb, req)

	if !ok {

		fmt.Fprintf(rw, "404", nil)
		return
	}

	response = r.controller()

	fmt.Fprintf(rw, response)
}
