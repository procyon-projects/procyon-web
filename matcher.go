package web

type RequestMapping struct {
}

type RequestMappingInfo struct {
	name                  string
	methodRequestMatcher  MethodRequestMatcher
	paramsRequestMatcher  ParametersRequestMatcher
	patternRequestMatcher PatternRequestMatcher
}

func NewRequestMappingInfo(name string,
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

func (condition RequestMappingInfo) getMethodRequestMatcher() MethodRequestMatcher {
	return condition.methodRequestMatcher
}

func (condition RequestMappingInfo) getParametersRequestMatcher() ParametersRequestMatcher {
	return condition.paramsRequestMatcher
}

func (condition RequestMappingInfo) getPatternRequestMatcher() PatternRequestMatcher {
	return condition.patternRequestMatcher
}

func (condition RequestMappingInfo) MatchRequest(req HttpRequest) interface{} {
	method := condition.methodRequestMatcher.MatchRequest(req)
	if method == nil {
		return nil
	} else {
		params := condition.paramsRequestMatcher.MatchRequest(req)
		if params == nil {
			return nil
		} else {
			pattern := condition.patternRequestMatcher.MatchRequest(req)
			if pattern == nil {
				return nil
			}
		}
	}
	return nil
}

type RequestMatcher interface {
	MatchRequest(req HttpRequest) interface{}
}

type MethodRequestMatcher struct {
	methods []RequestMethod
}

func NewMethodRequestMatcher(methods []RequestMethod) MethodRequestMatcher {
	return MethodRequestMatcher{
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

type ParametersRequestMatcher struct {
}

func NewParametersRequestMatcher() ParametersRequestMatcher {
	return ParametersRequestMatcher{}
}

func (matcher ParametersRequestMatcher) MatchRequest(req HttpRequest) interface{} {
	return nil
}

type PatternRequestMatcher struct {
}

func NewPatternRequestMatcher() PatternRequestMatcher {
	return PatternRequestMatcher{}
}

func (matcher PatternRequestMatcher) MatchRequest(req HttpRequest) interface{} {
	return nil
}
