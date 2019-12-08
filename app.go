package goKLC

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type App struct {
	key string
}

var _app *App
var _routeTree *RouteNode
var _middlewareList *MiddlewareNode
var _routeNameList routeNameList
var _configCollector *configCollector
var _config Config

func GetApp() *App {

	if _app == nil {
		_routeTree = NewRouteTree()
		_middlewareList = NewMiddlewareNode()
		_routeNameList = NewRouteNameList()
		_configCollector = newConfigCollector()
		_config = NewConfig()

		_app = &App{}
	}

	return _app
}

func (a *App) Run() {

	a.key = _config.Get("AppKey", "").(string)
	port := _config.Get("HttpPort", 8080)
	httpAddr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(httpAddr, a)

	fmt.Println(err)

}

func (a *App) Config() Config {

	return NewConfig()
}

func (a *App) Route() Route {
	rg := NewRouteGroup()
	rg.prefix = ""

	return rg.Route()
}

func GetRoute(name string) string {

	return _routeNameList.Get(name)
}

func (a *App) Middleware(m MiddlewareInterface) {
	if _middlewareList == nil {
		_middlewareList.middleware = m
	} else {
		mn := NewMiddlewareNode()
		mn.middleware = m

		_middlewareList.AddChild(mn)
	}
}

func (a *App) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	route, ok, params := match(req)
	request := NewRequest(req, params)
	var response *Response
	var middleware *MiddlewareNode

	if !ok {

		fmt.Fprintf(rw, "404", nil)
		return
	}

	response, middleware = _middlewareList.Handle(request)

	if response == nil {
		if route.middleware != nil {

			var rm *MiddlewareNode
			response, rm = route.middleware.Handle(request)

			if response == nil {
				response = route.controller(request)
			}

			rm.Terminate(response)
			rm = nil

		} else {
			response = route.controller(request)
		}
	}

	middleware.Terminate(response)

	if len(response.cookies) > 0 {
		writeCookies(rw, response)
	}

	rw.WriteHeader(response.status)
	rw.Write([]byte(response.content))

	request = nil
	response = nil
}

func writeCookies(rw http.ResponseWriter, r *Response) {
	for _, cookie := range r.cookies {
		c := http.Cookie{
			Name:   cookie.Name,
			Value:  cookie.Value,
			MaxAge: cookie.Duration,
			Path:   cookie.Path,
		}

		http.SetCookie(rw, &c)
	}
}
