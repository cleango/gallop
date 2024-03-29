package gallop

import (
	"context"
	"github.com/cleango/gallop/logger"
	"github.com/cleango/gallop/third_plugins/inject"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	usage = `
 _______       _ _             
(_______)     | | |            
 _   ___ _____| | | ___  ____  
| | (_  (____ | | |/ _ \|  _ \ 
| |___) / ___ | | | |_| | |_| |
 \_____/\_____|\_)_)___/|  __/ 
                        |_|    `
)

type Gallop struct {
	modulars  map[string][]IRouter
	engine    *gin.Engine
	op        *Options
	usage     string
	actions   []IAction
	configs   []interface{}
	closeList []IClose
}

//Ignite 项目初始化
func Ignite() *Gallop {
	op := DefaultOptions()
	op.AddFlags(pflag.CommandLine)
	InitFlags()
	initConfig(op.ConfigPath)
	g := &Gallop{
		modulars:  make(map[string][]IRouter),
		engine:    gin.New(),
		op:        op,
		usage:     usage,
		actions:   make([]IAction, 0),
		configs:   make([]interface{}, 0),
		closeList: make([]IClose, 0),
	}
	OpenCors(g.engine)
	g.Beans(logger.NewLogFactory())
	return g
}

type BeforeFunc func(app *gin.Engine)

func (g *Gallop) Before(funcs ...BeforeFunc) *Gallop {
	for _, f := range funcs {
		f(g.engine)
	}
	return g
}
func (g *Gallop) Use(middes ...IMidHandler) *Gallop {
	for _, mid := range middes {
		g.engine.Use(MidFactory(mid))
	}
	return g
}

func (g *Gallop) Validate(validators ...IValidator) *Gallop {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for _, item := range validators {
			vv, tag := item.Validate()
			if err := v.RegisterValidation(tag, vv); err != nil {
				logger.Fatal("binding validator Error")
			}
		}
	}

	return g
}

func (g *Gallop) Modular(name string, routers ...IRouter) *Gallop {
	for _, r := range routers {
		aop.Provide(&inject.Object{Value: r})
	}
	if r, ok := g.modulars[name]; ok {
		r = append(r, routers...)
		g.modulars[name] = r
		return g
	}
	g.modulars[name] = routers
	return g
}

//Preload 预执行，便于一些不需要web的项目，如果保持程序允许用开发自行控制
func (g *Gallop) Preload() {
	if err := aop.Populate(); err != nil {
		log.Fatal(err)
	}
	log.Println(g.usage)
	for _, v := range g.actions {
		v.Exec()
	}
}
func (g *Gallop) Actions(acts ...IAction) *Gallop {
	for _, v := range acts {
		if f, ok := v.(IClose); ok {
			g.closeList = append(g.closeList, f)
		}
		aop.Provide(&inject.Object{Value: v})
		g.actions = append(g.actions, v)
	}
	return g
}
func (g *Gallop) Launch(addr ...string) {
	g.Preload()
	if len(addr) > 0 {
		g.op.AddrPort = addr[0]
	}

	g.run(g.op)
}

func (g *Gallop) Banner(banner string) {
	g.usage = banner
}

func (g *Gallop) run(op *Options) {
	gin.SetMode(viper.GetString("mode"))
	for k, routers := range g.modulars {
		group := &Group{g.engine.Group(k)}
		for _, r := range routers {
			r.Builder(group)
		}
	}
	srv := &http.Server{
		Addr:    op.AddrPort,
		Handler: g.engine,
	}
	go func() {
		// service connections
		log.Printf("Listening and serving HTTP on %s\n", op.AddrPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	for _, v := range g.closeList {
		v.Shutdown(context.Background())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
