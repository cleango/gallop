package gallop

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IRouter interface {
	Builder(*Group)
}

type Group struct {
	*gin.RouterGroup
}

func (group *Group) Handle(httpMethod, relativePath string, handler interface{}) *Group {
	if h := Convert(handler); h != nil {
		group.RouterGroup.Handle(httpMethod, relativePath, h)
	}
	return group
}

// POST is a shortcut for router.Handle("POST", path, handle).
func (group *Group) POST(relativePath string, handlers interface{}) *Group {
	return group.Handle(http.MethodPost, relativePath, handlers)
}

// GET is a shortcut for router.Handle("GET", path, handle).
func (group *Group) GET(relativePath string, handlers interface{}) *Group {
	return group.Handle(http.MethodGet, relativePath, handlers)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle).
func (group *Group) DELETE(relativePath string, handlers interface{}) *Group {
	return group.Handle(http.MethodDelete, relativePath, handlers)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle).
func (group *Group) PATCH(relativePath string, handlers interface{}) *Group {
	return group.Handle(http.MethodPatch, relativePath, handlers)
}

// PUT is a shortcut for router.Handle("PUT", path, handle).
func (group *Group) PUT(relativePath string, handlers interface{}) *Group {
	return group.Handle(http.MethodPut, relativePath, handlers)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle).
func (group *Group) OPTIONS(relativePath string, handlers interface{}) *Group {
	return group.Handle(http.MethodOptions, relativePath, handlers)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle).
func (group *Group) HEAD(relativePath string, handlers interface{}) *Group {
	return group.Handle(http.MethodHead, relativePath, handlers)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (group *Group) Any(relativePath string, handlers interface{}) *Group {
	group.Handle(http.MethodGet, relativePath, handlers)
	group.Handle(http.MethodPost, relativePath, handlers)
	group.Handle(http.MethodPut, relativePath, handlers)
	group.Handle(http.MethodPatch, relativePath, handlers)
	group.Handle(http.MethodHead, relativePath, handlers)
	group.Handle(http.MethodOptions, relativePath, handlers)
	group.Handle(http.MethodDelete, relativePath, handlers)
	group.Handle(http.MethodConnect, relativePath, handlers)
	group.Handle(http.MethodTrace, relativePath, handlers)
	return group
}
