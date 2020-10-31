package web

type RequestMatcher interface {
	MatchRequest(req HttpRequest) interface{}
}

type MethodRequestMatcher struct {
	hash    int
	methods []RequestMethod
}

func newMethodRequestMatcher(methods []RequestMethod) MethodRequestMatcher {
	hashCode := 0
	for _, method := range methods {
		hashCode = 31*hashCode + hashCodeForString(string(method))
	}
	return MethodRequestMatcher{
		hashCode,
		methods,
	}
}

func (matcher MethodRequestMatcher) MatchRequest(req HttpRequest) interface{} {
	requestMethod := GetRequestMethod(req.GetMethod())
	for _, method := range matcher.methods {
		if method == requestMethod {
			return matcher
		}
	}
	return nil
}

func (matcher MethodRequestMatcher) hashCode() int {
	return matcher.hash
}

type ParametersRequestMatcher struct {
}

func newParametersRequestMatcher() ParametersRequestMatcher {
	return ParametersRequestMatcher{}
}

func (matcher ParametersRequestMatcher) MatchRequest(req HttpRequest) interface{} {
	return ""
}

func (matcher ParametersRequestMatcher) hashCode() int {
	return 0
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

func (matcher PatternRequestMatcher) MatchRequest(req HttpRequest) interface{} {
	path := req.GetPath()
	matches := make([]string, 0)
	for _, pattern := range matcher.patterns {
		result := matcher.pathMatcher.MatchPath(path, pattern)
		if result {
			matches = append(matches, pattern)
		}
	}
	if len(matches) > 0 {
		return matches
	}
	return nil
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

type RequestMatch interface {
	GetMapping() interface{}
	GetHandlerMethod() HandlerMethod
}

type DefaultRequestMatch struct {
	mapping       *RequestMapping
	handlerMethod HandlerMethod
}

func NewDefaultRequestMatch(mapping *RequestMapping, method HandlerMethod) RequestMatch {
	return DefaultRequestMatch{
		mapping,
		method,
	}
}
func (requestMatch DefaultRequestMatch) GetMapping() interface{} {
	return requestMatch.mapping
}

func (requestMatch DefaultRequestMatch) GetHandlerMethod() HandlerMethod {
	return requestMatch.handlerMethod
}
