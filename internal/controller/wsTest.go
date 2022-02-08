package controller

import (
	"fmt"
	"net/http"

	"zgo/internal/service"

	ulog "zgo/pkg/logs"
)

type PluginServiceHandle struct{}

// WsHandle 插件服务handle
func WsHandle(writer http.ResponseWriter, request *http.Request) {

	pluginClient := &service.Client{}
	//打开ws连接
	if pc, ok := pluginClient.Open(writer, request); ok {
		ulog.Log.Debug().Msg("ws connect ok")
		fmt.Println(pc)
		//go pc.ReadMessage()
	}
}

func (*PluginServiceHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ubx server")
}
