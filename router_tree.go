package web

type RouterTree struct {
	methodNodes []*RouterMethodNode
}

func newRouterTree() *RouterTree {
	return &RouterTree{
		make([]*RouterMethodNode, 0),
	}
}

func (tree *RouterTree) GetMethodNode(method RequestMethod) *RouterMethodNode {
	for _, methodNode := range tree.methodNodes {
		if methodNode.method == method {
			return methodNode
		}
	}
	methodNode := new(RouterMethodNode)
	methodNode.method = method
	tree.methodNodes = append(tree.methodNodes, methodNode)
	return methodNode
}

func (tree *RouterTree) AddHandlerMethod(path string, method RequestMethod, handler RequestHandlerFunc) {
	methodNode := tree.GetMethodNode(method)
	methodNode.AddRoute(path, handler)
}
