package web

import core "github.com/procyon-projects/procyon-core"

func registerPool(typ *core.Type, newFunc func() interface{}) {
	err := core.PoolManager.Register(typ, newFunc)
	if err != nil {
		panic(err)
	}
}

func getFromPool(typ *core.Type) interface{} {
	instance, err := core.PoolManager.Get(typ)
	if err != nil {
		panic(err)
	}
	return instance
}

func putToPool(instances ...interface{}) {
	for _, instance := range instances {
		core.PoolManager.Put(instance)
	}
}

var (
	httpRequestType  = core.GetType((*HttpRequest)(nil))
	httpResponseType = core.GetType((*HttpResponse)(nil))
)

func getHttpRequestFromPool() HttpRequest {
	return getFromPool(httpRequestType).(HttpRequest)
}

func getHttpResponseFromPool() HttpResponse {
	return getFromPool(httpResponseType).(HttpResponse)
}
