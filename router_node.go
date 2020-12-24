package web

type RouterNodeType byte

const (
	PathSegmentNode  RouterNodeType = 0
	PathVariableNode RouterNodeType = 1
	PathWildcardNode RouterNodeType = 2
)

type RouterPathNode struct {
	path                []byte
	length              uint
	nodeType            RouterNodeType
	handlerChain        *HandlerChain
	childNodes          []*RouterPathNode
	indices             []byte
	pathVariableNames   []string
	wildCardNode        *RouterPathNode
	pathVariableNode    *RouterPathNode
	childNode           *RouterPathNode
	hasWildcard         bool
	hasPathVariableNode bool
	childIndex          byte
	childStartIndex     byte
	childEndIndex       byte
}

func (node *RouterPathNode) addChildNode(child *RouterPathNode) {
	character := child.path[0]

	if len(node.childNodes) == 0 {
		node.childNodes = append(node.childNodes, nil)
		node.childStartIndex = character
		node.indices = []byte{0}
	} else {
		var difference byte
		if character < node.childStartIndex {
			difference = node.childStartIndex - character
		} else if character >= node.childEndIndex {
			difference = character - node.childEndIndex + 1
		}

		if character < node.childStartIndex {
			newIndices := make([]byte, difference+byte(len(node.indices)))
			copy(newIndices[difference:], node.indices)
			node.childStartIndex = character
			node.indices = newIndices
		} else if character >= node.childEndIndex {
			newIndices := make([]byte, difference+byte(len(node.indices)))
			copy(newIndices, node.indices)
			node.indices = newIndices
		}

	}

	node.childEndIndex = node.childStartIndex + byte(len(node.indices))
	index := node.indices[character-node.childStartIndex]

	if index == 0 {
		node.indices[character-node.childStartIndex] = byte(len(node.childNodes))
		node.childNodes = append(node.childNodes, child)
	} else {
		node.childNodes[index] = child
	}

	tempIndex := 47 - int(node.childStartIndex)
	if tempIndex >= 0 && len(node.indices) > tempIndex {
		node.childIndex = node.indices[tempIndex]
		node.childNode = node.childNodes[node.childIndex]
	}
}

func (node *RouterPathNode) handlePathSegment(path []byte, chain *HandlerChain) {

search:
	for {

		if len(path) == 0 {
			node.handlerChain = chain
			break
		}

		pathVariableIndex := -1
		tempIndex := 0
		for tempIndex < len(path) {
			if path[tempIndex] == ':' || path[tempIndex] == '*' {
				pathVariableIndex = tempIndex
				break
			}
			tempIndex++
		}

		if pathVariableIndex == -1 {
			if len(node.path) == 0 {
				node.path = path
				node.length = uint(len(path))
				node.handlerChain = chain
				break
			}

			child := &RouterPathNode{
				path:         path,
				length:       uint(len(path)),
				handlerChain: chain,
			}

			node.addChildNode(child)
			break
		}

		if pathVariableIndex == 0 {
			tempIndex := 1
			for tempIndex < len(path) {
				if path[tempIndex] == '/' {
					break
				}
				tempIndex++
			}

			pathVariableName := path[1:tempIndex]
			child := &RouterPathNode{
				path:   []byte("*"),
				length: 1,
			}

			if path[0] == ':' {
				child.nodeType = PathVariableNode
				pathVariableName := string(pathVariableName)
				if len(pathVariableName) == 0 {
					panic("Path variable cannot be empty " + string(path))
				}
				chain.updatePathVariableMetadata(len(chain.pathVariables), pathVariableName)
				chain.pathVariables = append(chain.pathVariables, pathVariableName)

				node.pathVariableNode = child
				node.hasPathVariableNode = true
				node = child
				path = path[tempIndex:]
				continue search
			}

			child.nodeType = PathWildcardNode
			child.handlerChain = chain
			node.wildCardNode = child
			node.hasWildcard = true
			break
		}

		if len(node.path) == 0 {
			node.path = path[:pathVariableIndex]
			node.length = uint(len(path[:pathVariableIndex]))
			path = path[pathVariableIndex:]
			continue search
		}

		child := &RouterPathNode{
			path:   path[:pathVariableIndex],
			length: uint(len(path[:pathVariableIndex])),
		}

		if child.path[0] == '/' {
			child.handlerChain = node.handlerChain
		}

		node.addChildNode(child)
		node = child
		path = path[pathVariableIndex:]
	}

}

func (node *RouterPathNode) updateMetadata(pathVariableIndex int, pathVariableName string) {

}
