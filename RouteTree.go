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
	node.next = &RouteNode{
		key:     key,
		dynamic: checkParams(key),
		route:   route,
		next:    nil,
		child:   nil,
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
	temp = node.child
	var params RouteParams = RouteParams{}

	for i, key := range path {
		if len(key) == 0 {
			key = "/"
		}

		temp = temp.FindNext(key)

		if len(temp.dynamic) > 0 {
			params[temp.dynamic] = key
		}

		if temp == nil {
			return nil, nil
		}

		if i == len(path)-1 {

			return temp, params
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

	for i, key := range path {
		if len(key) == 0 {
			key = "/"
		}

		if temp.child != nil {
			temp = temp.child
			tempNext = temp.FindNext(key)

			if tempNext != nil {
				temp = tempNext
			} else {
				if i != len(path)-1 {
					temp = temp.AddNext(key, nil)
				} else {
					temp = temp.AddNext(key, route)
				}
			}
		} else {
			if i != len(path)-1 {
				temp = temp.AddChild(key, nil)
			} else {
				temp = temp.AddChild(key, route)
			}
		}
	}
}

func (node *RouteNode) GetRoute() *Route {

	return node.route
}
