package gallop

import (
	"fmt"
	"github.com/cleango/gallop/third_plugins/inject"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"reflect"
)


func injectValue(elem reflect.Value, prefix string) {
	//注入Bean
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)
		conf := field.Tag.Get("value")
		if conf != "" {
			if prefix != "" {
				conf = fmt.Sprintf("%s.%s", prefix, conf)
			}
			if elem.Field(i).Kind() == reflect.Ptr {
				vv := reflect.New(elem.Field(i).Type().Elem())
				injectValue(vv.Elem(), conf)
				elem.Field(i).Set(reflect.ValueOf(vv.Interface()))
			} else if elem.Field(i).Kind() == reflect.Struct {
				injectValue(elem.Field(i), conf)
			} else {
				cc := viper.Get(conf)
				if cc != nil {
					val := reflect.ValueOf(cc)
					if !val.IsZero() {
						elem.Field(i).Set(val)
					}
				}
			}
		}

	}
}
func (g *Gallop) Beans(configs ...interface{}) *Gallop {

	for _, cc := range configs {
		g.configs = append(g.configs, cc)
		v := reflect.ValueOf(cc)
		if v.Kind() != reflect.Ptr {
			log.Fatal(" config is not ptr")
		}
		//获取值
		elem := v.Elem()
		injectValue(elem, "")
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
			if len(out) > 0 {
				aop.Provide(&obj)
			}
		}
		aop.Provide(&inject.Object{Value: cc})
	}
	g.onConfigChange()
	return g
}

func (g *Gallop) onConfigChange() {
	viper.OnConfigChange(func(in fsnotify.Event) {
		for _, cc := range g.configs {
			v := reflect.ValueOf(cc)
			//获取值
			elem := v.Elem()
			injectValue(elem, "")
		}
	})
}
