package goKLC

import (
	"fmt"
	"github.com/emirpasic/gods/maps/hashmap"
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
		address = fmt.Sprintf("%s/%s", checkPrefix(r.group.prefix), address)
	}

	r.address = checkPrefix(address)
	r.controller = controller

	rb.Put(r.address, r)
}

func match(rb *hashmap.Map, request *http.Request) (Route, bool) {
	path := request.URL.Path
	path = checkPrefix(path)

	r, ok := rb.Get(path)

	switch r.(type) {
	case Route:
		return r.(Route), ok
	default:
		return Route{}, false
	}
}

func checkPrefix(address string) string {
	if !strings.HasPrefix(address, "/") {
		return "/" + address
	}

	if strings.HasSuffix(address, "/") {
		address = address[:len(address)-1]
	}

	return address
}
