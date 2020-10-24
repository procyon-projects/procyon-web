package web

type RequestMapping struct {
}

type RequestMappingInfo struct {
	name                  string
	methodRequestMatcher  MethodRequestMatcher
	paramsRequestMatcher  ParametersRequestMatcher
	patternRequestMatcher PatternRequestMatcher
}

func newRequestMappingInfo(name string,
	methodRequestMatcher MethodRequestMatcher,
	paramsRequestMatcher ParametersRequestMatcher,
	patternRequestMatcher PatternRequestMatcher) RequestMappingInfo {
	return RequestMappingInfo{
		name,
		methodRequestMatcher,
		paramsRequestMatcher,
		patternRequestMatcher,
	}
}

func (mappingInfo RequestMappingInfo) getMethodRequestMatcher() MethodRequestMatcher {
	return mappingInfo.methodRequestMatcher
}

func (mappingInfo RequestMappingInfo) getParametersRequestMatcher() ParametersRequestMatcher {
	return mappingInfo.paramsRequestMatcher
}

func (mappingInfo RequestMappingInfo) getPatternRequestMatcher() PatternRequestMatcher {
	return mappingInfo.patternRequestMatcher
}

func (mappingInfo RequestMappingInfo) MatchRequest(req HttpRequest) interface{} {
	method := mappingInfo.methodRequestMatcher.MatchRequest(req)
	if method == nil {
		return nil
	} else {
		params := mappingInfo.paramsRequestMatcher.MatchRequest(req)
		if params == nil {
			return nil
		} else {
			pattern := mappingInfo.patternRequestMatcher.MatchRequest(req)
			if pattern == nil {
				return nil
			}
		}
	}
	return nil
}

func (mappingInfo RequestMappingInfo) hashCode() int {
	return 31*mappingInfo.patternRequestMatcher.hashCode() +
		mappingInfo.methodRequestMatcher.hashCode() +
		mappingInfo.paramsRequestMatcher.hashCode()
}

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
	return nil
}

func (matcher ParametersRequestMatcher) hashCode() int {
	return 0
}

type PatternRequestMatcher struct {
	hash     int
	patterns []string
}

func newPatternRequestMatcher(prefix string, paths []string) PatternRequestMatcher {
	patterns := make([]string, len(paths))
	hashCode := 0
	for index, path := range paths {
		patterns[index] = prefix + path
		hashCode = 31*hashCode + hashCodeForString(patterns[index])
	}
	return PatternRequestMatcher{
		hashCode,
		patterns,
	}
}

func (matcher PatternRequestMatcher) hashCode() int {
	return matcher.hash
}

func (matcher PatternRequestMatcher) MatchRequest(req HttpRequest) interface{} {
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
