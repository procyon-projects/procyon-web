package web

import (
	"github.com/valyala/fasthttp"
	"testing"
)

type testRoute struct {
	method        string
	path          string
	pathVariables []string
}

var (
	githubFastHttpRequests []fasthttp.Request
)

var githubAPI = []testRoute{
	// OAuth Authorizations
	{"GET", "/authorizations", []string{}},
	{"GET", "/authorizations/:id", []string{"id"}},
	{"POST", "/authorizations", []string{}},
	//{"PUT", "/authorizations/clients/:client_id"},
	//{"PATCH", "/authorizations/:id"},
	{"DELETE", "/authorizations/:id", []string{"id"}},
	{"GET", "/applications/:client_id/tokens/:access_token", []string{"client_id", "access_token"}},
	{"DELETE", "/applications/:client_id/tokens", []string{"client_id"}},
	{"DELETE", "/applications/:client_id/tokens/:access_token", []string{"client_id", "access_token"}},

	// Activity
	{"GET", "/events", []string{}},
	{"GET", "/repos/:owner/:repo/events", []string{"owner", "repo"}},
	{"GET", "/networks/:owner/:repo/events", []string{"owner", "repo"}},
	{"GET", "/orgs/:org/events", []string{"org"}},
	{"GET", "/users/:user/received_events", []string{"user"}},
	{"GET", "/users/:user/received_events/public", []string{"user"}},
	{"GET", "/users/:user/events", []string{"user"}},
	{"GET", "/users/:user/events/public", []string{"user"}},
	{"GET", "/users/:user/events/orgs/:org", []string{"user", "org"}},
	{"GET", "/feeds", []string{}},
	{"GET", "/notifications", []string{}},
	{"GET", "/repos/:owner/:repo/notifications", []string{"owner", "repo"}},
	{"PUT", "/notifications", []string{}},
	{"PUT", "/repos/:owner/:repo/notifications", []string{"owner", "repo"}},
	{"GET", "/notifications/threads/:id", []string{"id"}},
	//{"PATCH", "/notifications/threads/:id"},
	{"GET", "/notifications/threads/:id/subscription", []string{"id"}},
	{"PUT", "/notifications/threads/:id/subscription", []string{"id"}},
	{"DELETE", "/notifications/threads/:id/subscription", []string{"id"}},
	{"GET", "/repos/:owner/:repo/stargazers", []string{"owner", "repo"}},
	{"GET", "/users/:user/starred", []string{"user"}},
	{"GET", "/user/starred", []string{}},
	{"GET", "/user/starred/:owner/:repo", []string{"owner", "repo"}},
	{"PUT", "/user/starred/:owner/:repo", []string{"owner", "repo"}},
	{"DELETE", "/user/starred/:owner/:repo", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/subscribers", []string{"owner", "repo"}},
	{"GET", "/users/:user/subscriptions", []string{"user"}},
	{"GET", "/user/subscriptions", []string{}},
	{"GET", "/repos/:owner/:repo/subscription", []string{"owner", "repo"}},
	{"PUT", "/repos/:owner/:repo/subscription", []string{"owner", "repo"}},
	{"DELETE", "/repos/:owner/:repo/subscription", []string{"owner", "repo"}},
	{"GET", "/user/subscriptions/:owner/:repo", []string{"owner", "repo"}},
	{"PUT", "/user/subscriptions/:owner/:repo", []string{"owner", "repo"}},
	{"DELETE", "/user/subscriptions/:owner/:repo", []string{"owner", "repo"}},

	// Gists
	{"GET", "/users/:user/gists", []string{"user"}},
	{"GET", "/gists", []string{}},
	//{"GET", "/gists/public"},
	//{"GET", "/gists/starred"},
	{"GET", "/gists/:id", []string{"id"}},
	{"POST", "/gists", []string{}},
	//{"PATCH", "/gists/:id"},
	{"PUT", "/gists/:id/star", []string{"id"}},
	{"DELETE", "/gists/:id/star", []string{"id"}},
	{"GET", "/gists/:id/star", []string{"id"}},
	{"POST", "/gists/:id/forks", []string{"id"}},
	{"DELETE", "/gists/:id", []string{"id"}},

	// Git Data
	{"GET", "/repos/:owner/:repo/git/blobs/:sha", []string{"owner", "repo", "sha"}},
	{"POST", "/repos/:owner/:repo/git/blobs", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/git/commits/:sha", []string{"owner", "repo", "sha"}},
	{"POST", "/repos/:owner/:repo/git/commits", []string{"owner", "repo"}},
	//{"GET", "/repos/:owner/:repo/git/refs/*ref"},
	{"GET", "/repos/:owner/:repo/git/refs", []string{"owner", "repo"}},
	{"POST", "/repos/:owner/:repo/git/refs", []string{"owner", "repo"}},
	//{"PATCH", "/repos/:owner/:repo/git/refs/*ref"},
	//{"DELETE", "/repos/:owner/:repo/git/refs/*ref"},
	{"GET", "/repos/:owner/:repo/git/tags/:sha", []string{"owner", "repo", "sha"}},
	{"POST", "/repos/:owner/:repo/git/tags", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/git/trees/:sha", []string{"owner", "repo", "sha"}},
	{"POST", "/repos/:owner/:repo/git/trees", []string{"owner", "repo"}},

	// Issues
	{"GET", "/issues", []string{}},
	{"GET", "/user/issues", []string{}},
	{"GET", "/orgs/:org/issues", []string{"org"}},
	{"GET", "/repos/:owner/:repo/issues", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/issues/:number", []string{"owner", "repo", "number"}},
	{"POST", "/repos/:owner/:repo/issues", []string{"owner", "repo"}},
	//{"PATCH", "/repos/:owner/:repo/issues/:number"},
	{"GET", "/repos/:owner/:repo/assignees", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/assignees/:assignee", []string{"owner", "repo", "assignee"}},
	{"GET", "/repos/:owner/:repo/issues/:number/comments", []string{"owner", "repo", "number"}},
	//{"GET", "/repos/:owner/:repo/issues/comments"},
	//{"GET", "/repos/:owner/:repo/issues/comments/:id"},
	{"POST", "/repos/:owner/:repo/issues/:number/comments", []string{"owner", "repo", "number"}},
	//{"PATCH", "/repos/:owner/:repo/issues/comments/:id"},
	//{"DELETE", "/repos/:owner/:repo/issues/comments/:id"},
	{"GET", "/repos/:owner/:repo/issues/:number/events", []string{"owner", "repo", "number"}},
	//{"GET", "/repos/:owner/:repo/issues/events"},
	//{"GET", "/repos/:owner/:repo/issues/events/:id"},
	{"GET", "/repos/:owner/:repo/labels", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/labels/:name", []string{"owner", "repo", "name"}},
	{"POST", "/repos/:owner/:repo/labels", []string{"owner", "repo"}},
	//{"PATCH", "/repos/:owner/:repo/labels/:name"},
	{"DELETE", "/repos/:owner/:repo/labels/:name", []string{"owner", "repo", "name"}},
	{"GET", "/repos/:owner/:repo/issues/:number/labels", []string{"owner", "repo", "number"}},
	{"POST", "/repos/:owner/:repo/issues/:number/labels", []string{"owner", "repo", "number"}},
	{"DELETE", "/repos/:owner/:repo/issues/:number/labels/:name", []string{"owner", "repo", "number", "name"}},
	{"PUT", "/repos/:owner/:repo/issues/:number/labels", []string{"owner", "repo", "number"}},
	{"DELETE", "/repos/:owner/:repo/issues/:number/labels", []string{"owner", "repo", "number"}},
	{"GET", "/repos/:owner/:repo/milestones/:number/labels", []string{"owner", "repo", "number"}},
	{"GET", "/repos/:owner/:repo/milestones", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/milestones/:number", []string{"owner", "repo", "number"}},
	{"POST", "/repos/:owner/:repo/milestones", []string{"owner", "repo"}},
	//{"PATCH", "/repos/:owner/:repo/milestones/:number"},
	{"DELETE", "/repos/:owner/:repo/milestones/:number", []string{"owner", "repo", "number"}},

	// Miscellaneous
	{"GET", "/emojis", []string{}},
	{"GET", "/gitignore/templates", []string{}},
	{"GET", "/gitignore/templates/:name", []string{"name"}},
	{"POST", "/markdown", []string{}},
	{"POST", "/markdown/raw", []string{}},
	{"GET", "/meta", []string{}},
	{"GET", "/rate_limit", []string{}},

	// Organizations
	{"GET", "/users/:user/orgs", []string{"user"}},
	{"GET", "/user/orgs", []string{}},
	{"GET", "/orgs/:org", []string{"org"}},
	//{"PATCH", "/orgs/:org"},
	{"GET", "/orgs/:org/members", []string{"org"}},
	{"GET", "/orgs/:org/members/:user", []string{"org", "user"}},
	{"DELETE", "/orgs/:org/members/:user", []string{"org", "user"}},
	{"GET", "/orgs/:org/public_members", []string{"org"}},
	{"GET", "/orgs/:org/public_members/:user", []string{"org", "user"}},
	{"PUT", "/orgs/:org/public_members/:user", []string{"org", "user"}},
	{"DELETE", "/orgs/:org/public_members/:user", []string{"org", "user"}},
	{"GET", "/orgs/:org/teams", []string{"org"}},
	{"GET", "/teams/:id", []string{"id"}},
	{"POST", "/orgs/:org/teams", []string{"org"}},
	//{"PATCH", "/teams/:id"},
	{"DELETE", "/teams/:id", []string{"id"}},
	{"GET", "/teams/:id/members", []string{"id"}},
	{"GET", "/teams/:id/members/:user", []string{"id", "user"}},
	{"PUT", "/teams/:id/members/:user", []string{"id", "user"}},
	{"DELETE", "/teams/:id/members/:user", []string{"id", "user"}},
	{"GET", "/teams/:id/repos", []string{"id"}},
	{"GET", "/teams/:id/repos/:owner/:repo", []string{"id", "owner", "repo"}},
	{"PUT", "/teams/:id/repos/:owner/:repo", []string{"id", "owner", "repo"}},
	{"DELETE", "/teams/:id/repos/:owner/:repo", []string{"id", "owner", "repo"}},
	{"GET", "/user/teams", []string{}},

	// Pull Requests
	{"GET", "/repos/:owner/:repo/pulls", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/pulls/:number", []string{"owner", "repo", "number"}},
	{"POST", "/repos/:owner/:repo/pulls", []string{"owner", "repo"}},
	//{"PATCH", "/repos/:owner/:repo/pulls/:number"},
	{"GET", "/repos/:owner/:repo/pulls/:number/commits", []string{"owner", "repo", "number"}},
	{"GET", "/repos/:owner/:repo/pulls/:number/files", []string{"owner", "repo", "number"}},
	{"GET", "/repos/:owner/:repo/pulls/:number/merge", []string{"owner", "repo", "number"}},
	{"PUT", "/repos/:owner/:repo/pulls/:number/merge", []string{"owner", "repo", "number"}},
	{"GET", "/repos/:owner/:repo/pulls/:number/comments", []string{"owner", "repo", "number"}},
	//{"GET", "/repos/:owner/:repo/pulls/comments"},
	//{"GET", "/repos/:owner/:repo/pulls/comments/:number"},
	{"PUT", "/repos/:owner/:repo/pulls/:number/comments", []string{"owner", "repo", "number"}},
	//{"PATCH", "/repos/:owner/:repo/pulls/comments/:number"},
	//{"DELETE", "/repos/:owner/:repo/pulls/comments/:number"},

	// Repositories
	{"GET", "/user/repos", []string{}},
	{"GET", "/users/:user/repos", []string{"user"}},
	{"GET", "/orgs/:org/repos", []string{"org"}},
	{"GET", "/repositories", []string{}},
	{"POST", "/user/repos", []string{}},
	{"POST", "/orgs/:org/repos", []string{"org"}},
	{"GET", "/repos/:owner/:repo", []string{"owner", "repo"}},
	//{"PATCH", "/repos/:owner/:repo"},
	{"GET", "/repos/:owner/:repo/contributors", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/languages", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/teams", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/tags", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/branches", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/branches/:branch", []string{"owner", "repo", "branch"}},
	{"DELETE", "/repos/:owner/:repo", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/collaborators", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/collaborators/:user", []string{"owner", "repo", "user"}},
	{"PUT", "/repos/:owner/:repo/collaborators/:user", []string{"owner", "repo", "user"}},
	{"DELETE", "/repos/:owner/:repo/collaborators/:user", []string{"owner", "repo", "user"}},
	{"GET", "/repos/:owner/:repo/comments", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/commits/:sha/comments", []string{"owner", "repo", "sha"}},
	{"POST", "/repos/:owner/:repo/commits/:sha/comments", []string{"owner", "repo", "sha"}},
	{"GET", "/repos/:owner/:repo/comments/:id", []string{"owner", "repo", "id"}},
	//{"PATCH", "/repos/:owner/:repo/comments/:id"},
	{"DELETE", "/repos/:owner/:repo/comments/:id", []string{"owner", "repo", "id"}},
	{"GET", "/repos/:owner/:repo/commits", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/commits/:sha", []string{"owner", "repo", "sha"}},
	{"GET", "/repos/:owner/:repo/readme", []string{"owner", "repo"}},
	//{"GET", "/repos/:owner/:repo/contents/*path"},
	//{"PUT", "/repos/:owner/:repo/contents/*path"},
	//{"DELETE", "/repos/:owner/:repo/contents/*path"},
	//{"GET", "/repos/:owner/:repo/:archive_format/:ref"},
	{"GET", "/repos/:owner/:repo/keys", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/keys/:id", []string{"owner", "repo", "id"}},
	{"POST", "/repos/:owner/:repo/keys", []string{"owner", "repo"}},
	//{"PATCH", "/repos/:owner/:repo/keys/:id"},
	{"DELETE", "/repos/:owner/:repo/keys/:id", []string{"owner", "repo", "id"}},
	{"GET", "/repos/:owner/:repo/downloads", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/downloads/:id", []string{"owner", "repo", "id"}},
	{"DELETE", "/repos/:owner/:repo/downloads/:id", []string{"owner", "repo", "id"}},
	{"GET", "/repos/:owner/:repo/forks", []string{"owner", "repo"}},
	{"POST", "/repos/:owner/:repo/forks", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/hooks", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/hooks/:id", []string{"owner", "repo", "id"}},
	{"POST", "/repos/:owner/:repo/hooks", []string{"owner", "repo"}},
	//{"PATCH", "/repos/:owner/:repo/hooks/:id"},
	{"POST", "/repos/:owner/:repo/hooks/:id/tests", []string{"owner", "repo", "id"}},
	{"DELETE", "/repos/:owner/:repo/hooks/:id", []string{"owner", "repo", "id"}},
	{"POST", "/repos/:owner/:repo/merges", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/releases", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/releases/:id", []string{"owner", "repo", "id"}},
	{"POST", "/repos/:owner/:repo/releases", []string{"owner", "repo"}},
	//{"PATCH", "/repos/:owner/:repo/releases/:id"},
	{"DELETE", "/repos/:owner/:repo/releases/:id", []string{"owner", "repo", "id"}},
	{"GET", "/repos/:owner/:repo/releases/:id/assets", []string{"owner", "repo", "id"}},
	{"GET", "/repos/:owner/:repo/stats/contributors", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/stats/commit_activity", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/stats/code_frequency", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/stats/participation", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/stats/punch_card", []string{"owner", "repo"}},
	{"GET", "/repos/:owner/:repo/statuses/:ref", []string{"owner", "repo", "ref"}},
	{"POST", "/repos/:owner/:repo/statuses/:ref", []string{"owner", "repo", "ref"}},

	// Search
	{"GET", "/search/repositories", []string{}},
	{"GET", "/search/code", []string{}},
	{"GET", "/search/issues", []string{}},
	{"GET", "/search/users", []string{}},
	{"GET", "/legacy/issues/search/:owner/:repository/:state/:keyword", []string{"owner", "repository", "state", "keyword"}},
	{"GET", "/legacy/repos/search/:keyword", []string{"keyword"}},
	{"GET", "/legacy/user/search/:keyword", []string{"keyword"}},
	{"GET", "/legacy/user/email/:email", []string{"email"}},

	// Users
	{"GET", "/users/:user", []string{"user"}},
	{"GET", "/user", []string{}},
	//{"PATCH", "/user"},
	{"GET", "/users", []string{}},
	{"GET", "/user/emails", []string{}},
	{"POST", "/user/emails", []string{}},
	{"DELETE", "/user/emails", []string{}},
	{"GET", "/users/:user/followers", []string{"user"}},
	{"GET", "/user/followers", []string{}},
	{"GET", "/users/:user/following", []string{"user"}},
	{"GET", "/user/following", []string{}},
	{"GET", "/user/following/:user", []string{"user"}},
	{"GET", "/users/:user/following/:target_user", []string{"user", "target_user"}},
	{"PUT", "/user/following/:user", []string{"user"}},
	{"DELETE", "/user/following/:user", []string{"user"}},
	{"GET", "/users/:user/keys", []string{"user"}},
	{"GET", "/user/keys", []string{}},
	{"GET", "/user/keys/:id", []string{"id"}},
	{"POST", "/user/keys", []string{}},
	//{"PATCH", "/user/keys/:id"},
	{"DELETE", "/user/keys/:id", []string{"id"}},
}

func init() {
	githubFastHttpRequests = make([]fasthttp.Request, 0)
	for _, githubRoute := range githubAPI {
		req := fasthttp.AcquireRequest()
		req.SetRequestURI(githubRoute.path)
		req.Header.SetMethod(githubRoute.method)
		req.URI().SetPath(githubRoute.path)
		githubFastHttpRequests = append(githubFastHttpRequests, *req)
	}
}

func TestRouter(t *testing.T) {
	router := newRouterTree()
	for _, route := range githubAPI {
		handlerChain := NewHandlerChain(func(context *WebRequestContext) {
			context.SetModel(route.method + ":" + route.path)
		}, nil, nil)
		router.AddRoute(route.path, RequestMethod(route.method), handlerChain)
	}

	webRequestContext := newWebRequestContext().(*WebRequestContext)
	fastHttpRequestContext := &fasthttp.RequestCtx{}

	for index, route := range githubAPI {
		fastHttpRequestContext.Request = githubFastHttpRequests[index]
		webRequestContext.reset()
		webRequestContext.fastHttpRequestContext = fastHttpRequestContext

		if route.path == "/repos/:owner/:repo/milestones" && route.method == "GET" {
			router.Get(webRequestContext)

		}
	}
}
