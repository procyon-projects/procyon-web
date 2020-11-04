package web

import (
	"strings"
)

type PathMatcher interface {
	MatchPath(requestContext RequestContext, path string, pattern string) bool
}

type SimplePathMatcher struct {
	patternCache     map[string][]string
	pathSeparator    uint8
	pathVariableSign uint8
}

func NewSimplePathMatcher() SimplePathMatcher {
	return SimplePathMatcher{
		patternCache:     make(map[string][]string),
		pathSeparator:    '/',
		pathVariableSign: ':',
	}
}

func (pathMatcher SimplePathMatcher) MatchPath(requestContext RequestContext, path string, pattern string) bool {
	if pattern == path {
		return true
	}
	return pathMatcher.match(requestContext.(*WebRequestContext), path, pattern)
}

func (pathMatcher SimplePathMatcher) RegisterPatternCache(pattern string) {
	patternDirs := strings.FieldsFunc(pattern, func(c rune) bool { return c == rune(pathMatcher.pathSeparator) })
	pathMatcher.patternCache[pattern] = patternDirs
}

func (pathMatcher SimplePathMatcher) match(requestContext *WebRequestContext, path string, pattern string) bool {
	if path[0] == pathMatcher.pathSeparator && pattern[0] == pathMatcher.pathSeparator {
		paramIndex := 0
		for _, patternDir := range pathMatcher.patternCache[pattern] {
			if patternDir[0] == pathMatcher.pathVariableSign {
				end := 0
				for end < len(path) && path[end] != pathMatcher.pathSeparator {
					end++
				}
				requestContext.pathVariables = requestContext.pathVariables[:paramIndex+1]
				requestContext.pathVariables[paramIndex].Key = patternDir[1:]
				requestContext.pathVariables[paramIndex].Value = path[:end]
			} else if path[1:len(patternDir)+1] == patternDir {
				path = path[len(patternDir)+2:]
			} else {
				return false
			}
		}
	}
	return true
}
