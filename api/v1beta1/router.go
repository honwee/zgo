package v1beta1

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"

	"zgo/internal/controller"
	ulog "zgo/pkg/logs"
)

type gm struct {
	r    *mux.Router
	addr string
}

type Interface interface {
	Create()
	register()
	run() error
}

func New() Interface {
	return &gm{}
}

func (g *gm) Create() {
	g.r = mux.NewRouter()
	g.register()
	err := g.run()
	if err != nil {
		fmt.Println(err)
		return
	}
	ulog.Log.Debug().Msg("路由实例初始化成功！")
}

func (g *gm) register() {
	wsRouter := g.r.PathPrefix("/v1beta1").Subrouter()

	wsRouter.HandleFunc("", controller.WsHandle)
	err := wsRouter.Walk(gorillaWalkFn)
	if err != nil {
		//ulog.Log.Err(err).Msg("wsRouter.Walk err")
		return
	}
	ulog.Log.Debug().Msg("路由注册成功！")

}

func (g *gm) run() error {
	server := &http.Server{
		Addr:    g.addr,
		Handler: g.r,
	}
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:8889", "0.0.0.0"))
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr.String())
	if err != nil {
		return err
	}
	//strArr := strings.Split(ln.Addr().String(), ":")
	//Port <- strArr[len(strArr)-1]
	fmt.Println(ln.Addr().String())
	err = server.Serve(ln)
	if nil != err {
		//ulog.Log.Err(err).Msg("Serve run err")
		return err
	}
	return nil
}

func gorillaWalkFn(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	path, err := route.GetPathTemplate()
	if err != nil {
		return err
	}
	ulog.Log.Debug().Msgf("路由加载 [path: %+v]  handle:%+T", path, route.GetHandler())
	return nil
}
