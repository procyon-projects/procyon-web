package web

import (
	"regexp"
	"strings"
)

type RouterNodeType uint8

const PathSegmentNode RouterNodeType = 0
const PathVariableNode RouterNodeType = 1

type RouterPathNode struct {
	nodeType          RouterNodeType
	path              string
	fullPath          string
	handler           RequestHandlerFunc
	childNodes        []*RouterPathNode
	indices           string
	pathVariableNames []string
	pathVariableRegex []*regexp.Regexp
}

func (node *RouterPathNode) AddChildNode(path string, fullPath string, handler RequestHandlerFunc) {

	fullPathIndex := 0
	startIndex := 0
	endIndex := len(path) - 1
	var pathVariableNames []string
	var pathVariableRegex []*regexp.Regexp

visit:
	for {

		if startIndex >= endIndex {
			break
		}

		matchIndex := 0
		for matchIndex < len(node.path) && matchIndex < len(path) && path[matchIndex] == node.path[matchIndex] {
			matchIndex++
		}

		if node.nodeType == PathVariableNode {
			character := path[0]
			if character != '{' {
				if len(node.childNodes) == 0 {
					tempEndIndex := 0
					for ; tempEndIndex < len(path) && path[tempEndIndex] != '{'; tempEndIndex++ {
					}
					pathNode := &RouterPathNode{
						nodeType: PathSegmentNode,
						path:     path[matchIndex:tempEndIndex],
						fullPath: sanitizePath(fullPath[:fullPathIndex+matchIndex+tempEndIndex]),
					}
					node.childNodes = append(node.childNodes, pathNode)
					node = pathNode
					startIndex += len(path[:tempEndIndex])
					path = path[tempEndIndex:]
					continue

				} else {
					node = node.childNodes[0]
					continue
				}
			} else {
				processed, variableName, regex := node.handlePathVariable(0, len(path), path, fullPath)
				pathVariableNames = append(pathVariableNames, variableName)
				pathVariableRegex = append(pathVariableRegex, regex)
				startIndex += processed
				path = path[processed:]
				fullPathIndex += processed
				node = node.childNodes[0]
				continue
			}
		} else if node.nodeType == PathSegmentNode {
			if matchIndex < len(node.path) {

				character := path[0]
				if character != '{' {
					newNode := &RouterPathNode{
						nodeType:          PathSegmentNode,
						path:              node.path[matchIndex:],
						childNodes:        node.childNodes,
						indices:           node.indices,
						handler:           node.handler,
						fullPath:          node.fullPath,
						pathVariableNames: node.pathVariableNames,
						pathVariableRegex: node.pathVariableRegex,
					}
					node.childNodes = []*RouterPathNode{newNode}
					node.indices = string([]byte{node.path[matchIndex]})
					node.path = sanitizePath(path[:matchIndex])
					node.handler = nil
					node.pathVariableNames = nil
					node.pathVariableRegex = nil
					node.fullPath = sanitizePath(fullPath[:fullPathIndex+matchIndex])
					startIndex += matchIndex
					fullPathIndex += len(node.path)
				}
			}

			if matchIndex < len(path) {
				path = path[matchIndex:]

				character := path[0]

				searchCharacter := character
				if character == '{' {
					searchCharacter = '*'
				}

				for i := 0; i < len(node.indices); i++ {
					if searchCharacter == node.indices[i] {
						if searchCharacter == '*' {
							fullPathIndex += len(node.path)
							startIndex += len(node.path)
							processed, variableName, regex := node.handlePathVariable(0, len(path), path, fullPath)
							pathVariableNames = append(pathVariableNames, variableName)
							pathVariableRegex = append(pathVariableRegex, regex)
							fullPathIndex += processed
							startIndex += processed
							path = path[processed:]
						} else {
							fullPathIndex += len(node.path)
							startIndex += len(node.path)
						}
						node = node.childNodes[i]
						continue visit
					}
				}

				if character != '{' {
					tempEndIndex := 0
					for ; tempEndIndex < len(path) && path[tempEndIndex] != '{'; tempEndIndex++ {
					}
					pathSegmentNode := &RouterPathNode{
						nodeType: PathSegmentNode,
						path:     sanitizePath(path[:tempEndIndex]),
						fullPath: fullPath,
					}
					startIndex += len(path[:tempEndIndex])
					startIndex += matchIndex
					path = path[tempEndIndex:]

					node.indices += string([]byte{character})
					node.childNodes = append(node.childNodes, pathSegmentNode)
					node = pathSegmentNode
					continue
				} else {
					processed, variableName, regex := node.handlePathVariable(0, len(path), path, fullPath)
					pathVariableNames = append(pathVariableNames, variableName)
					pathVariableRegex = append(pathVariableRegex, regex)

					variableNode := &RouterPathNode{
						nodeType: PathVariableNode,
						path:     sanitizePath(path[:processed]),
						fullPath: sanitizePath(fullPath[:fullPathIndex+processed]),
					}
					startIndex += processed
					path = path[processed:]
					fullPathIndex += processed

					node.indices += string([]byte{'*'})
					node.childNodes = append(node.childNodes, variableNode)
					node = variableNode
					continue
				}
			}
		}
	}
	node.handler = handler
	node.fullPath = sanitizePathSegment(fullPath, false)
	node.pathVariableNames = pathVariableNames
	node.pathVariableRegex = pathVariableRegex
}

func (node *RouterPathNode) handlePathVariable(startIndex int, endIndex int, path string, fullPath string) (int, string, *regexp.Regexp) {
	// path variables

	offset := 0
	variableName := ""
	variableEndIndex := -1
	colonIndex := -1
	var regex *regexp.Regexp

	tempEndIndex := startIndex + 1
	if startIndex != -1 {
		for tempEndIndex < endIndex && path[tempEndIndex] != '/' {
			if path[tempEndIndex] == '}' {
				variableEndIndex = tempEndIndex
			}
			tempEndIndex++
		}
		if variableEndIndex == -1 && path[tempEndIndex] == '}' && tempEndIndex == endIndex {
			variableEndIndex = tempEndIndex
		}
	}

	if startIndex != -1 && variableEndIndex == -1 {
		panic("Close the bracket : " + path[startIndex:] + " in path " + fullPath)
	} else if startIndex != -1 && variableEndIndex != -1 {
		if variableEndIndex-startIndex < 2 {
			panic("Give a name your path variable" + path[startIndex:endIndex] + " in path " + fullPath)
		}

		pathVariable := path[startIndex : variableEndIndex+1]
		offset = variableEndIndex - startIndex + 1
		colonIndex = strings.Index(pathVariable, ":")

		if colonIndex == -1 {

			variableName = pathVariable[1 : len(pathVariable)-1]
		} else {
			variablePattern := pathVariable[colonIndex+1 : len(pathVariable)-1]
			variableName = pathVariable[1:colonIndex]
			result, err := regexp.Compile("^" + variablePattern + "$")

			if err != nil {
				panic("Invalid regex " + path[startIndex:endIndex] + " in path " + fullPath)
			}
			regex = result
		}

		if len(variableName) == 0 {
			panic("Give a name your path variable" + path[startIndex:endIndex] + " in path " + fullPath)
		}

		return offset, variableName, regex
	}
	return offset, "", regex
}

type RouterMethodNode struct {
	method           RequestMethod
	root             *RouterPathNode
	registeredRoutes []string
}

func (methodNode *RouterMethodNode) AddRoute(path string, handler RequestHandlerFunc) {
	if methodNode.root == nil {
		rootNode := &RouterPathNode{
			path:     "/",
			fullPath: "/",
		}
		rootNode.AddChildNode(path, path, handler)
		methodNode.root = rootNode
	} else {
		conflictCheck(methodNode.registeredRoutes, path)
		methodNode.root.AddChildNode(path, path, handler)
	}
	methodNode.registeredRoutes = append(methodNode.registeredRoutes, path)
}

func splitRoute(c rune) bool {
	return c == '/'
}

func conflictCheck(routes []string, route string) {
	for _, registeredRoute := range routes {
		conflictCheckRoute(registeredRoute, route)
	}
}

func conflictCheckRoute(registeredRoute string, route string) {
	splitRegisteredRoutes := strings.FieldsFunc(registeredRoute, splitRoute)
	splitRoutes := strings.FieldsFunc(route, splitRoute)

	fullRegisteredRoute := ""
	fullRoute := ""
	if len(splitRegisteredRoutes) == len(splitRoutes) {
		index := 0
		for ; index < len(splitRegisteredRoutes); index++ {
			routeToken := splitRoutes[index]
			if routeToken == "" {
				panic("Path segment cannot be null")
			}
			registeredRouteToken := splitRegisteredRoutes[index]

			if isPathVariable(registeredRouteToken) && isPathVariable(routeToken) {
				fullRegisteredRoute += "/*"
				fullRoute += "/*"
			} else if isPathVariable(registeredRouteToken) && !isPathVariable(routeToken) {
				fullRegisteredRoute += "/" + routeToken
				fullRoute += "/" + routeToken
			} else if !isPathVariable(registeredRouteToken) && isPathVariable(routeToken) {
				fullRegisteredRoute += "/" + registeredRouteToken
				fullRoute += "/" + registeredRouteToken
			} else {
				fullRegisteredRoute += "/" + registeredRouteToken
				fullRoute += "/" + routeToken
			}
		}
		if fullRoute == fullRegisteredRoute {
			panic("Conflict between " + route + " and " + registeredRoute)
		}
	}
}

func isPathVariable(str string) bool {
	if str[0] == '{' {
		return true
	}
	return false
}

func sanitizePath(path string) string {
	return sanitizePathSegment(path, true)
}

func sanitizePathSegment(path string, wildcard bool) string {
	offset := 0
	result := ""
	colon := false
	for index := 0; index < len(path); index++ {
		if colon && path[index] != '/' {
			offset++
			continue
		} else {
			colon = false
		}
		if path[index] != ':' {
			continue
		}
		result = result + path[offset:index]
		offset = index
		colon = true
	}
	result = result + path[offset:]
	if !wildcard {
		return result
	}
	return sanitizePathVariable(result)
}

func sanitizePathVariable(path string) string {
	if path == "/" || path == "" {
		return path
	}
	result := ""
	tokens := strings.FieldsFunc(path, splitRoute)
	for _, token := range tokens {
		if token[0] == '{' {
			result += "/*"
		} else {
			result += "/" + token
		}
	}
	if path[len(path)-1] == '/' && result[len(result)-1] != '/' {
		result += "/"
	}
	if path[0] != '/' && result[0] == '/' {
		result = result[1:]
	}
	return result
}
