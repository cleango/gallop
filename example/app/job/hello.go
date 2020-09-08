package job

import (
	"fmt"
	"github.com/cleango/gallop/example/app/config"
	"time"
)

type Hello struct {
	isRun bool
	Demo *config.Demo  `inject:""`
}

func (h *Hello) Run() {
	defer func() {
		h.isRun=false
	}()
	fmt.Println(h.Demo.Name)
	fmt.Println("hello",h.isRun)
	h.isRun=true
	time.Sleep(5*time.Second)
}

type Hello1 struct {
	isRun bool
	Demo *config.Demo  `inject:"demo"`
}

func (h *Hello1) Run() {
	defer func() {
		h.isRun=false
	}()
	fmt.Println(h.Demo.Name)
	fmt.Println("hello",h.isRun)
	h.isRun=true
	time.Sleep(5*time.Second)
}
