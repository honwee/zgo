/**
 * Copyright (C) 2021 UnionTech Software Technology Co., Ltd. All rights reserved.
 * @author 陈弘唯
 * @Email  : chenhongwei@uniontech.com
 * @date 2021/4/29 下午8:52
 */

package config

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	Log struct {
		Level    string // debug info warn error fatal
		FilePath string // 日志输出文件路径
		ToStdout bool   // 是否输出到Stdout, 默认不输出到Stdout
	}

	Store struct {
		Address  string `toml:"address"`
		LocalDir string `toml:"localdir"`
		Token    string `toml:"token"`
	}
}

var (
	// C 加载后的配置信息
	C = config{}
	// Model 程序运行模式
	Model string
	// Version go build -ldflags "-X 'ubx/server/configs.Version=${VERSION}'"
	Version string
)

func LoadConfig(name, suffix string) error {
	viper.SetConfigName(name)
	viper.SetConfigType(suffix)
	//rootPath, err := GetRootPath()
	rootPath := "/usr/bin/ubx"
	//if err != nil || len(rootPath) == 0 {
	//	fmt.Println("GetRootPath:", rootPath)
	//	return err
	//}
	viper.AddConfigPath(path.Join(rootPath, "configs"))

	// 设定默认值
	viper.SetDefault("Log.MaxSize", 100)
	viper.SetDefault("Log.MaxBackups", 3)

	viper.SetDefault("addr.host", "")
	viper.SetDefault("addr.port", "8017")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	// 解析配置文件
	if err := viper.Unmarshal(&C); err != nil {
		return err
	}
	return nil
}

// 获取user home
func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		fmt.Println("Getenv get user home:", home)
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	if result := strings.TrimSpace(stdout.String()); result != "" && result != "~" {
		fmt.Println("shell get user home:", result)
		return result, nil
	}
	return "/root", nil
}
