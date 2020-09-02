package main

import (
	"fmt"
	"github.com/cleango/gallop"
	"github.com/cleango/gallop/example/app/config"
	"github.com/cleango/gallop/example/app/router"
)

func main() {
	app:=gallop.Ignite().
		Beans(config.NewConfiguration()).
		Modular("",router.NewRouter()).
		Modular("api",router.NewApi())
	res,err:=gallop.GetBeanByName("demo")
	fmt.Println(err,res.(*config.Demo))
	app.Launch()
}
