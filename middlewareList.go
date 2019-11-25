package goKLC

type MiddlewareNode struct {
	middleware *Middleware
	child      *MiddlewareNode
	parent     *MiddlewareNode
}

func NewMiddlewareNode() *MiddlewareNode {
	return &MiddlewareNode{}
}

func (node *MiddlewareNode) AddChild(childNode *MiddlewareNode) {
	if node == nil {
		node = childNode
	} else if node.child != nil {
		node.child.AddChild(childNode)
	} else {
		childNode.parent = node
		childNode.child = node.child
		node.child = childNode
	}
}

func (node *MiddlewareNode) AddParent(parentNode *MiddlewareNode) {
	parentNode.child = node
	parentNode.parent = node.parent
	node.parent = parentNode
}

func (node *MiddlewareNode) Child() *MiddlewareNode {

	return node.child
}

func (node *MiddlewareNode) Parent() *MiddlewareNode {

	return node.parent
}

func (node *MiddlewareNode) Handle(request *Request) *MiddlewareNode {
	if node.middleware != nil {
		node.middleware.Handle(request)
	}

	if node.child != nil {
		return node.child.Handle(request)
	}

	return node
}

func (node *MiddlewareNode) Terminate() {
	if node.middleware != nil {
		node.middleware.Terminate()
	}

	if node.parent != nil {
		node.parent.Terminate()
	}
}
