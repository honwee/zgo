package boot

import (
	"path"

	"github.com/dimiro1/banner"
	"github.com/mattn/go-colorable"

	"zgo/pkg/config"
	"zgo/pkg/logs"
	"zgo/pkg/ws"
)

func init() {
	showBanner()
	//new log
	logs.NewLog(logs.Config{
		Level:    "debug",
		FilePath: path.Clean("./"),
		ToStdout: true,
	})
	ws.WebsocketHub = ws.CreateHubFactory()
	go ws.WebsocketHub.Run()
}

func showBanner() {

	templ := `{{ .Title "ZGO" "" 4 }}
   {{ .AnsiColor.BrightCyan }} * Copyright (C) 2021 UnionTech Software Technology Co., Ltd. All rights reserved.
    * @Email  : chenhongwei@uniontech.com{{ .AnsiColor.Default }}
    * GoVersion: {{ .GoVersion }}
    * GOOS: {{ .GOOS }}
    * GOARCH: {{ .GOARCH }}
    * Version:{{ .AnsiColor.BrightGreen }} ` + config.Version + `{{ .AnsiColor.Default }}
    * StartTime:{{ .AnsiColor.BrightRed }} {{ .Now "2006-01-02 15:04:05" }} {{ .AnsiColor.Default }}
 
`

	banner.InitString(colorable.NewColorableStdout(), true, true, templ)
}
