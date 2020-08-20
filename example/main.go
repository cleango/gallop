package main

import (
	"github.com/cleango/gallop"
	"github.com/cleango/gallop/example/app/config"
	"github.com/cleango/gallop/example/app/router"
)

func main() {
	gallop.Ignite().
		Beans(config.NewConfiguration()).
		Modular("",router.NewRouter()).
		Modular("api",router.NewApi()).
		Launch()
}
