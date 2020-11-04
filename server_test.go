package web

import (
	"net/http"
	"runtime"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(1)
}

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func benchRequest(b *testing.B, router http.Handler, r *http.Request) {
	w := new(mockResponseWriter)
	u := r.URL
	rq := u.RawQuery
	r.RequestURI = u.RequestURI()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		u.RawQuery = rq
		router.ServeHTTP(w, r)
	}
}

func procyonHandle() {}

func setUpProcyonSingle(method RequestMethod, path string, handlerFunc RequestHandlerFunc) http.Handler {
	handlerRegistry := newSimpleHandlerRegistry()
	handlerRegistry.Register(NewHandler(handlerFunc, WithMethod(method), WithPath(path)))
	server := NewProcyonWebServerForBenchmark(handlerRegistry)
	if server != nil {
		return server
	}
	return nil
}

func BenchmarkProcyon_Param(b *testing.B) {
	router := setUpProcyonSingle(RequestMethodGet, "/user/:name}", procyonHandle)

	request, _ := http.NewRequest("GET", "/user/test", nil)
	benchRequest(b, router, request)
}

const fiveBrace = "/:a/:b/:c/:d/:e}"
const fiveRoute = "/test/test/test/test/test"

func BenchmarkProcyon_Param5(b *testing.B) {
	router := setUpProcyonSingle("GET", fiveBrace, procyonHandle)

	request, _ := http.NewRequest("GET", fiveRoute, nil)
	benchRequest(b, router, request)
}

const twentyBrace = "/:a/:b/:c/:d/:e/:f/:g/:h/:i/:j/:k/:l/:m/:n/:o/:p/:q/:r/:s/:t}"
const twentyRoute = "/a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p/q/r/s/t"

func BenchmarkProcyon_Param20(b *testing.B) {
	router := setUpProcyonSingle("GET", twentyBrace, procyonHandle)

	request, _ := http.NewRequest("GET", twentyRoute, nil)
	benchRequest(b, router, request)
}
