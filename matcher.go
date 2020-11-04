package web

import "net/http"

type RequestMatcher interface {
	MatchRequest(requestContext RequestContext, req *http.Request) bool
}

type MethodRequestMatcher struct {
	hash      int
	methods   []RequestMethod
	methodMap map[string]bool
}

func newMethodRequestMatcher(methods []RequestMethod) MethodRequestMatcher {
	methodMap := make(map[string]bool)
	hashCode := 0
	for _, method := range methods {
		methodMap[string(method)] = true
		hashCode = 31*hashCode + hashCodeForString(string(method))
	}
	return MethodRequestMatcher{
		hashCode,
		methods,
		methodMap,
	}
}

func (matcher MethodRequestMatcher) MatchRequest(requestContext RequestContext, req *http.Request) bool {
	if _, ok := matcher.methodMap[req.Method]; ok {
		return true
	}
	return false
}

func (matcher MethodRequestMatcher) hashCode() int {
	return matcher.hash
}

type PatternRequestMatcher struct {
	pathMatcher PathMatcher
	hash        int
	patterns    []string
}

func newPatternRequestMatcher(pathMatcher PathMatcher, prefix string, paths []string) PatternRequestMatcher {
	patterns := make([]string, len(paths))
	hashCode := 0
	for index, path := range paths {
		patterns[index] = prefix + path
		pathMatcher.(SimplePathMatcher).RegisterPatternCache(patterns[index])
		hashCode = 31*hashCode + hashCodeForString(patterns[index])
	}
	return PatternRequestMatcher{
		pathMatcher,
		hashCode,
		patterns,
	}
}

func (matcher PatternRequestMatcher) hashCode() int {
	return matcher.hash
}

func (matcher PatternRequestMatcher) MatchRequest(requestContext RequestContext, req *http.Request) bool {
	for _, pattern := range matcher.patterns {
		result := matcher.pathMatcher.MatchPath(requestContext, req.URL.Path, pattern)
		if result {
			return true
		}
	}
	return false
}

func hashCodeForString(str string) int {
	if len(str) > 0 {
		hash := 0
		for _, character := range str {
			hash = 31*hash + int(character)
		}
		return hash
	}
	return 1
}
