package goKLC

import (
	"fmt"
	"net/http"
	"strings"
)

const GET Method = "GET"
const POST Method = "POST"
const PUT Method = "PUT"
const PATCH Method = "PATCH"
const DELETE Method = "DELETE"

type Method string

type RouteParams map[string]interface{}

type RouteGroup struct {
	prefix     string
	name       string
	middleware *MiddlewareNode
}

type Route struct {
	address    string
	name       string
	controller ControllerFunc
	group      *RouteGroup
	method     Method
	middleware *MiddlewareNode
}

func NewRouteTree() *RouteNode {
	return &RouteNode{
		key:   "",
		route: &Route{},
		next:  nil,
		child: nil,
	}
}

func NewRouteGroup() RouteGroup {

	return RouteGroup{}
}

func NewRoute() Route {

	return Route{}
}

func (rg *RouteGroup) Route() Route {
	r := NewRoute()
	r.group = rg

	if rg.middleware != nil {
		r.middleware = rg.middleware
	}

	return r
}

func (rg RouteGroup) Group(prefix string) RouteGroup {
	newRg := NewRouteGroup()
	newRg.prefix = fmt.Sprintf("%s/%s", checkPrefix(rg.prefix), prefix)

	if rg.middleware != nil {
		newRg.middleware = rg.middleware
	}

	if len(rg.name) > 0 {
		newRg.name = rg.name
	}

	return newRg
}

func (r Route) Group(prefix string) *RouteGroup {
	rg := NewRouteGroup()

	rg.prefix = prefix

	if r.group != nil && r.group.middleware != nil {

		rg.middleware = r.group.middleware
	}

	return &rg
}

func (r Route) Get(address string, controller ControllerFunc) *Route {

	return addNewRoute(r, address, controller, GET)
}

func (r Route) Post(address string, controller ControllerFunc) *Route {

	return addNewRoute(r, address, controller, POST)
}

func (r Route) Put(address string, controller ControllerFunc) *Route {

	return addNewRoute(r, address, controller, PUT)
}

func (r Route) Patch(address string, controller ControllerFunc) *Route {

	return addNewRoute(r, address, controller, PATCH)
}

func (r Route) Delete(address string, controller ControllerFunc) *Route {

	return addNewRoute(r, address, controller, DELETE)
}

func (r *Route) Middleware(m MiddlewareInterface) *Route {
	mn := NewMiddlewareNode()
	mn.middleware = m

	if r.middleware == nil {
		r.middleware = mn
	} else {
		r.middleware.AddChild(mn)
	}

	return r
}

func (rg *RouteGroup) Middleware(m MiddlewareInterface) *RouteGroup {
	mn := NewMiddlewareNode()
	mn.middleware = m

	if rg.middleware == nil {
		rg.middleware = mn
	} else {
		rg.middleware.AddChild(mn)
	}

	return rg
}

func (r *Route) Name(name string) *Route {
	if r.group != nil && len(r.group.name) > 0 {
		name = fmt.Sprintf("%s.%s", r.group.name, name)
	}

	r.name = name

	routes.Add(r.name, r.address)

	return r
}

func (rg *RouteGroup) Name(name string) *RouteGroup {
	if len(rg.name) > 0 {
		name = fmt.Sprintf("%s.%s", rg.name, name)
	}

	rg.name = name

	return rg
}

func addNewRoute(r Route, address string, controller ControllerFunc, method Method) *Route {
	if len(r.group.prefix) > 0 {
		address = checkPrefix(r.group.prefix) + "/" + checkPrefix(address)
	}

	if r.group.middleware != nil {
		r.middleware = r.group.middleware
	}

	r.address = checkPrefix(address)
	r.controller = controller
	r.method = method

	path := getPath(r.address, method)
	routeTree.AddFromPath(path, &r)

	return &r
}

func match(request *http.Request) (*Route, bool, RouteParams) {
	url := request.URL.Path
	url = checkPrefix(url)
	path := getPath(url, Method(request.Method))

	node, params := routeTree.FindFromPath(path)

	if node == nil || node.GetRoute() == nil {
		return nil, false, nil
	}

	return node.route, true, params
}

func checkPrefix(address string) string {
	if strings.HasPrefix(address, "/") {
		address = strings.TrimPrefix(address, "/")
	}

	if strings.HasSuffix(address, "/") {
		address = strings.TrimSuffix(address, "/")
	}

	return address
}

func checkParams(path string) string {
	if strings.HasPrefix(path, "$") {
		return strings.TrimPrefix(path, "$")
	}

	return ""
}

func getPath(url string, method Method) []string {
	url = fmt.Sprintf("%s/%s", method, url)

	return strings.Split(url, "/")
}
