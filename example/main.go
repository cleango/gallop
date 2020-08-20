package main

import (
	"github.com/gallop"
	"github.com/gallop/example/app/config"
	"github.com/gallop/example/app/router"
)

func main() {
	gallop.Ignite().
		Beans(config.NewConfiguration()).
		Modular("",router.NewRouter()).
		Modular("api",router.NewApi()).
		Launch()
}
