package router

import (
	"github.com/gallop"
	"github.com/gallop/example/app/controller"
)

type Router struct {
	Hello *controller.HelloController `inject:""`
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Builder(group *gallop.Group) {
	group.GET("/", r.Hello.Hello)
	group.GET("/json", r.Hello.Json)
}
