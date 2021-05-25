package config

import (
	"fmt"
	"github.com/cleango/gallop"
	"time"
)

type Demo struct {
	Name string
}
type Configuration struct {
	Name string `value:"name"`
	A    int    `value:"a"`
	B    *struct {
		C string `value:"c"`
	} `value:"b"`
}

func (c *Demo) Shutdown(ctx gallop.CloseContext) {
	print("close")
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}

func (c *Configuration) NewDemo() *Demo {

	go func() {
		for {
			time.Sleep(5 * time.Second)
			fmt.Println(c.B.C, "======")
		}
	}()
	return &Demo{c.Name}
}

func (c *Configuration) NewDemo1() (string, *Demo) {
	return "demo", &Demo{Name: "larry"}
}
