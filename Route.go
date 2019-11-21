package goKLC

import (
	"fmt"
	"net/http"
	"strings"
)

type RouteParams map[string]string

type RouteGroup struct {
	prefix string
}

type Route struct {
	address    string
	name       string
	controller ControllerFunc
	group      *RouteGroup
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

	path := strings.Split(r.address, "/")
	routeTree.AddFromPath(path, &r)
}

func match(request *http.Request) (*Route, bool, RouteParams) {
	path := request.URL.Path
	path = checkPrefix(path)

	node, params := routeTree.FindFromPath(strings.Split(path, "/"))

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
