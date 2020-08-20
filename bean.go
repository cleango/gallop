package gallop

import (
	"github.com/cleango/gallop/third_plugins/inject"
	"github.com/spf13/viper"
	"log"
	"reflect"
)

func (g *Gallop) Beans(configs ...interface{}) *Gallop {
	for _, cc := range configs {
		v := reflect.ValueOf(cc)
		if v.Kind() != reflect.Ptr {
			log.Fatal(" config is not ptr")
		}
		//获取值
		elem := v.Elem()
		//注入Bean
		for i := 0; i < elem.NumField(); i++ {
			field := elem.Type().Field(i)
			conf := field.Tag.Get("value")
			if conf != "" {
				val:=reflect.ValueOf(viper.Get(conf))
				if !val.IsZero(){
					elem.Field(i).Set(val)
				}
			}

		}
		//注入Bean

		for i := 0; i < v.NumMethod(); i++ {
			method := v.Method(i)
			out := method.Call(nil)
			obj := inject.Object{}
			for _, v := range out {
				switch v.Kind() {
				case reflect.String:
					obj.Name = v.Interface().(string)
				case reflect.Ptr:
					obj.Value = v.Interface()
				default:
					log.Fatal("configuration type is error")
				}
			}
			g.aop.Provide(&obj)
		}
	}
	return g
}
