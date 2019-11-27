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
	prefix string
}

type Route struct {
	address    string
	name       string
	controller ControllerFunc
	group      *RouteGroup
	method     Method
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

func (rg RouteGroup) Route() Route {
	r := NewRoute()
	r.group = &rg

	return r
}

func (rg RouteGroup) Group(prefix string) RouteGroup {
	newRg := NewRouteGroup()
	newRg.prefix = fmt.Sprintf("%s/%s", checkPrefix(rg.prefix), prefix)

	return newRg
}

func (r Route) Group(prefix string) RouteGroup {
	rg := NewRouteGroup()

	rg.prefix = prefix

	return rg
}

func (r Route) Get(address string, controller ControllerFunc) {
	if len(r.group.prefix) > 0 {
		address = checkPrefix(r.group.prefix) + address
	}

	r.address = checkPrefix(address)
	r.controller = controller
	r.method = GET

	path := getPath(r.address, GET)
	routeTree.AddFromPath(path, &r)
}

func (r Route) Post(address string, controller ControllerFunc) {
	if len(r.group.prefix) > 0 {
		address = checkPrefix(r.group.prefix) + address
	}

	r.address = checkPrefix(address)
	r.controller = controller
	r.method = POST

	path := getPath(r.address, POST)
	routeTree.AddFromPath(path, &r)
}

func (r Route) Put(address string, controller ControllerFunc) {
	if len(r.group.prefix) > 0 {
		address = checkPrefix(r.group.prefix) + address
	}

	r.address = checkPrefix(address)
	r.controller = controller
	r.method = PUT

	path := getPath(r.address, PUT)
	routeTree.AddFromPath(path, &r)
}

func (r Route) Patch(address string, controller ControllerFunc) {
	if len(r.group.prefix) > 0 {
		address = checkPrefix(r.group.prefix) + address
	}

	r.address = checkPrefix(address)
	r.controller = controller
	r.method = PATCH

	path := getPath(r.address, PATCH)
	routeTree.AddFromPath(path, &r)
}

func (r Route) Delete(address string, controller ControllerFunc) {
	if len(r.group.prefix) > 0 {
		address = checkPrefix(r.group.prefix) + address
	}

	r.address = checkPrefix(address)
	r.controller = controller
	r.method = DELETE

	path := getPath(r.address, DELETE)
	routeTree.AddFromPath(path, &r)
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

	if !strings.HasSuffix(address, "/") {
		address = address + "/"
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
