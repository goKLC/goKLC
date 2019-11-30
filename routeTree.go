package goKLC

var temp *RouteNode
var tempNext *RouteNode

type RouteNode struct {
	key     string
	dynamic string
	route   *Route
	next    *RouteNode
	child   *RouteNode
}

func (node *RouteNode) AddChild(key string, route *Route) *RouteNode {
	node.child = &RouteNode{
		key:     key,
		dynamic: checkParams(key),
		route:   route,
		next:    nil,
		child:   nil,
	}

	return node.child
}

func (node *RouteNode) AddNext(key string, route *Route) *RouteNode {

	if node.next != nil {
		return node.next.AddNext(key, route)
	} else {
		node.next = &RouteNode{
			key:     key,
			dynamic: checkParams(key),
			route:   route,
			next:    nil,
			child:   nil,
		}
	}

	return node.next
}

func (node *RouteNode) FindNext(key string) *RouteNode {
	if node.key == key || len(node.dynamic) > 0 {

		return node
	} else if node.next != nil {

		return node.next.FindNext(key)
	} else {

		return nil
	}
}

func (node *RouteNode) FindFromPath(path []string) (*RouteNode, RouteParams) {
	temp = node
	var params RouteParams = RouteParams{}

	for i, key := range path {
		temp = temp.FindNext(key)

		if temp == nil {
			return nil, nil
		}

		if len(temp.dynamic) > 0 {
			params[temp.dynamic] = key
		}

		if i == len(path)-1 {
			return temp.child, params
		}

		temp = temp.child

		if temp == nil {
			return nil, nil
		}
	}

	return temp, params
}

func (node *RouteNode) AddFromPath(path []string, route *Route) {
	temp = node
	var r *Route

	for i, key := range path {
		if i == len(path)-1 {
			r = route
		} else {
			r = &Route{}
		}

		tempNext = temp.FindNext(key)

		if tempNext == nil {
			tempNext = temp.AddNext(key, &Route{})
		}

		temp = tempNext
		tempNext = temp.child

		if tempNext == nil {
			tempNext = temp.AddChild("/", r)
		}

		temp = tempNext
	}
}

func (node *RouteNode) GetRoute() *Route {

	return node.route
}
