package gallop

import (
	"context"
	"github.com/cleango/gallop/logger"
	"github.com/cleango/gallop/third_plugins/inject"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Gallop struct {
	modulars map[string][]IRouter
	engine   *gin.Engine
	op       *Options
}

//Ignite 项目初始化
func Ignite() *Gallop {
	op := DefaultOptions()
	op.AddFlags(pflag.CommandLine)
	InitFlags()
	initConfig(op.ConfigPath)
	g:= &Gallop{
		modulars: make(map[string][]IRouter),
		engine:   gin.New(),
		op:       op,
	}
	g.Beans(logger.NewLogFactory())
	return g
}
func (g *Gallop) Use(middes ...IMidHandler) *Gallop {
	for _, mid := range middes {
		g.engine.Use(MidFactory(mid))
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

func (g *Gallop) Launch(addr ...string) {

	if err := aop.Populate(); err != nil {
		log.Fatal(err)
	}
	if len(addr) > 0 {
		g.op.AddrPort = addr[0]
	}
	g.run(g.op)
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
