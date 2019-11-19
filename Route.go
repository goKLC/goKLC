package goKLC

import (
	"fmt"
	"net/http"
	"strings"
)

type RouteGroup struct {
	prefix string
}

type Route struct {
	address    string
	name       string
	controller func() string
	group      *RouteGroup
}

func NewRouteTree() *Node {
	return &Node{
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

func (r Route) Get(address string, controller func() string) {
	if len(r.group.prefix) > 0 {
		address = checkPrefix(r.group.prefix) + address
	}

	r.address = checkPrefix(address)
	r.controller = controller

	path := strings.Split(r.address, "/")
	routeTree.AddFromPath(path, &r)
}

func match(request *http.Request) (*Route, bool) {
	path := request.URL.Path
	path = checkPrefix(path)

	node := routeTree.FindFromPath(strings.Split(path, "/"))

	if node == nil || node.GetRoute() == nil {
		return nil, false
	}

	return node.route, true
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
