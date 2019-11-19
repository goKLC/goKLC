package goKLC

var temp *Node
var tempNext *Node

type Node struct {
	key   string
	route *Route
	next  *Node
	child *Node
}

func (node *Node) AddChild(key string, route *Route) *Node {
	node.child = &Node{
		key:   key,
		route: route,
		next:  nil,
		child: nil,
	}

	return node.child
}

func (node *Node) AddNext(key string, route *Route) *Node {
	node.next = &Node{
		key:   key,
		route: route,
		next:  nil,
		child: nil,
	}

	return node.next
}

func (node *Node) FindNext(key string) *Node {
	if node.key == key {

		return node
	} else if node.next != nil {

		return node.next.FindNext(key)
	} else {

		return nil
	}
}

func (node *Node) FindFromPath(path []string) *Node {
	temp = node.child

	for i, key := range path {
		if len(key) == 0 {
			key = "/"
		}
		temp = temp.FindNext(key)

		if temp == nil {
			return nil
		}

		if i == len(path)-1 {

			return temp
		}

		temp = temp.child

		if temp == nil {
			return nil
		}
	}

	return temp
}

func (node *Node) AddFromPath(path []string, route *Route) {
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

func (node *Node) GetRoute() *Route {

	return node.route
}
