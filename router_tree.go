package web

type RouterTree struct {
	methodNodes []*RouterMethodNode
}

func newRouterTree() *RouterTree {
	return &RouterTree{
		make([]*RouterMethodNode, 0),
	}
}

func (tree *RouterTree) GetHandlerMethod(path string, method RequestMethod) {
	methodNode := tree.GetMethodNode(method)
	if methodNode == nil || methodNode.root == nil {
		return
	}
	routePath := tree.search(methodNode.root, 0, path, 0)
	if routePath != nil {

	}
}

func (tree *RouterTree) search(node *RouterPathNode, index int, path string, processed int) interface{} {
	var wildcardNode *RouterPathNode

search:
	if index >= len(path) {
		if node.handler != nil {
			return node
		}
		return nil
	}

	if node.nodeType == PathSegmentNode {

		for index < len(path) && (index-processed) < len(node.path) && path[index] == node.path[index-processed] {
			index++
		}

		if len(node.path) != index-processed {
			return nil
		}

		if index >= len(path) {
			if node.handler != nil {
				return node
			}
			return nil
		}

		wildcardNode = nil
		for nodeIndex := 0; nodeIndex < len(node.indices); nodeIndex++ {
			if path[index] == node.indices[nodeIndex] {
				processed += len(node.path)
				result := tree.search(node.childNodes[nodeIndex], index, path, processed)
				if result != nil {
					return result
				}
				break
			} else if '*' == node.indices[nodeIndex] {
				wildcardNode = node.childNodes[nodeIndex]
			}
		}

		if wildcardNode != nil {
			processed += len(node.path)
			result := tree.search(wildcardNode, index, path, processed)
			if result != nil {
				return result
			}
		}

	} else if node.nodeType == PathVariableNode {
		tempIndex := index
		for ; index < len(path) && path[index] != '/'; index++ {
		}

		processed += index - tempIndex
		if len(node.childNodes) != 0 {
			result := tree.search(node.childNodes[0], index, path, processed)
			if result != nil {
				return result
			}
		} else {
			goto search
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

func (tree *RouterTree) AddHandlerMethod(path string, method RequestMethod, handler RequestHandlerFunc) {
	methodNode := tree.GetMethodNode(method)
	methodNode.AddRoute(path, handler)
}
