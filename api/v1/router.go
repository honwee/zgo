package v1

import (
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
)

type gm struct {
	r    *mux.Router
	addr string
}

type Interface interface {
	Start()
	register()
	run() error
}

func New() Interface {
	return &gm{}
}

func (g *gm) Start() {
	g.r = mux.NewRouter()
	g.register()
	err := g.run()
	if err != nil {
		fmt.Println(err)
		return
	}

	//ulog.Log.Debug().Msg("路由实例初始化成功！")
}

func (g *gm) register() {
	wsRouter := g.r.PathPrefix("/v1").Subrouter()
	//wsRouter.Handle("/", &plugin.PluginServiceHandle{})
	//wsRouter.HandleFunc("/ws", plugin.WsHandle)
	////离线安装
	//wsRouter.HandleFunc("/getList", plugin.GetListHandle)
	//wsRouter.HandleFunc("/install", plugin.InstallHandle)

	//if "debug" == config.Model {
	wsRouter.HandleFunc("/debug/pprof/", pprof.Index)
	wsRouter.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	wsRouter.HandleFunc("/debug/pprof/profile", pprof.Profile)
	wsRouter.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	wsRouter.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	wsRouter.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	wsRouter.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	wsRouter.Handle("/debug/pprof/block", pprof.Handler("block"))
	wsRouter.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	wsRouter.Handle("/debug/pprof/trace", pprof.Handler("trace"))
	wsRouter.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	//}

	err := wsRouter.Walk(gorillaWalkFn)
	if err != nil {
		//ulog.Log.Err(err).Msg("wsRouter.Walk err")
		return
	}
	//ulog.Log.Debug().Msg("路由注册成功！")

}

func (g *gm) run() error {
	server := &http.Server{
		Addr:    g.addr,
		Handler: g.r,
	}
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr.String())
	if err != nil {
		return err
	}
	//strArr := strings.Split(ln.Addr().String(), ":")
	//Port <- strArr[len(strArr)-1]
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
	fmt.Println(path)
	//ulog.Log.Debug().Msgf("路由加载 [path: %+v]  handle:%+T", path, route.GetHandler())
	return nil
}
