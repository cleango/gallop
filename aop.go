package gallop

import "github.com/cleango/gallop/third_plugins/inject"

var (
	aop      inject.Graph
)

func GetBeanByName(name string) (interface{},error)  {
	return aop.GetObjectByName(name)
}