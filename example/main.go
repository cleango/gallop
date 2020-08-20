package main

import (
	"gallop"
	"gallop/example/app/config"
	"gallop/example/app/router"
)

func main() {
	gallop.Ignite().
		Beans(config.NewConfiguration()).
		Modular("",router.NewRouter()).
		Modular("api",router.NewApi()).
		Launch()
}
