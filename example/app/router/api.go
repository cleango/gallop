package router

import (
	"gallop"
	"gallop/example/app/controller"
)

type Api struct {
	Hello *controller.HelloController `inject:""`
}

func NewApi() *Api {
	return &Api{}
}

func (a *Api) Builder(group *gallop.Group) {
	group.GET("/",a.Hello.Hello)
}
