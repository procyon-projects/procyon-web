package web

type RouterTree struct {
	methodNodes []*RouterMethodNode
}

func newRouterTree() *RouterTree {
	return &RouterTree{
		make([]*RouterMethodNode, 0),
	}
}

func (tree *RouterTree) GetHandlerMethod(path string, method RequestMethod) *HandlerMethod {
	methodNode := tree.GetMethodNode(method)
	if methodNode == nil {
		return nil
	}
	return tree.search(methodNode.root, methodNode.root.wildCardNode, path)
}

func (tree *RouterTree) search(node *RouterPathNode, wildcardNode *RouterPathNode, path string) *HandlerMethod {

	if len(path) == 0 {
		return node.handler
	}

	if node != nil {
		path = path[len(node.path):]
		for nodeIndex := 0; nodeIndex < len(node.indices); nodeIndex++ {
			if path[0] == node.indices[nodeIndex] {
				return tree.search(node.childNodes[nodeIndex], node.childNodes[nodeIndex].wildCardNode, path)
			}
		}
	}
	if wildcardNode != nil {
		end := 0
		for end < len(path) && path[end] != '/' {
			end++
		}

		path = path[end:]

		if len(path) == 0 {
			return wildcardNode.handler
		}

		if len(wildcardNode.childNodes) != 0 {
			if wildcardNode.childNodes[0].nodeType == PathVariableNode {
				return tree.search(nil, wildcardNode.childNodes[0].wildCardNode, path)
			} else {
				return tree.search(wildcardNode.childNodes[0], wildcardNode.childNodes[0].wildCardNode, path)
			}
		}
	}
	return nil
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

func (tree *RouterTree) AddHandlerMethod(path string, method RequestMethod, handler *HandlerMethod) {
	methodNode := tree.GetMethodNode(method)
	methodNode.AddRoute(path, handler)
}
