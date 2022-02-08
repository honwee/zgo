/**
 * Copyright (C) 2021 UnionTech Software Technology Co., Ltd. All rights reserved.
 * @author 陈弘唯
 * @Email  : chenhongwei@uniontech.com
 * @date 2021/12/28 上午11:29
 */

package logs

import (
	"crypto/md5" // #nosec
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/longbozhan/timewriter"
	"github.com/rs/zerolog"

	"zgo/pkg/encryption"
)

type (
	// Logger 别名
	Logger = zerolog.Logger
	// Level 别名
	Level = zerolog.Level
)

var (
	Log    *Logger
	aesKey [16]byte
)

func NewLog(c Config) {
	Log = New(Config{
		Level:    c.Level,
		FilePath: c.FilePath,
		ToStdout: c.ToStdout,
		MaxAge:   c.MaxAge,
	})
}

var levelMap = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"panic": zerolog.PanicLevel,
	"fatal": zerolog.FatalLevel,
}

// Config 可用在配置文件中
type Config struct {
	Level    string
	FilePath string            // 日志文件路径
	MaxAge   int               // days
	Fields   map[string]string // slog的初始化字段(session)
	ToStdout bool              // 默认不输出到Stdout
}

// New log
func New(c Config) *Logger {
	if c.Level == "" {
		c.Level = "error"
	}
	_, ok := levelMap[c.Level]
	if !ok {
		c.Level = "error"
	}
	// init zerolog format
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z07:00"
	zerolog.DurationFieldInteger = true
	zerolog.TimestampFieldName = "timestamp"
	zerolog.DurationFieldUnit = time.Millisecond
	out := &timewriter.TimeWriter{
		Dir:        c.FilePath,
		Compress:   true,
		ReserveDay: c.MaxAge,
	}

	zerolog.SetGlobalLevel(levelMap[c.Level])

	var zOut io.Writer
	//w := defaultConsoleWriter()
	w := zerolog.NewConsoleWriter()
	w.Out = out
	zOut = w
	// 同时输出到文件和控制台
	if c.ToStdout {
		zOut = zerolog.MultiLevelWriter(zOut, func() io.Writer {

			//return defaultConsoleWriter()
			return zerolog.NewConsoleWriter()
		}())
	}
	zc := zerolog.New(zOut).With().Timestamp().Caller()
	for k, v := range c.Fields {
		zc = zc.Str(k, v)
	}
	slog := zc.Logger()
	return &slog
}

//Output console
func defaultConsoleWriter() zerolog.ConsoleWriter {
	/* #nosec */
	var key string
	w := zerolog.NewConsoleWriter()
	w.TimeFormat = "2006-01-02 15:04:05.000"
	w.NoColor = true
	w.FormatTimestamp = func(i interface{}) string {
		s, ok := i.(string)
		if !ok {
			return fmt.Sprintf("%v ", i)
		}

		sArr := strings.Split(s, ".")
		if len(sArr) > 1 && len(sArr[1]) > 3 {
			key = sArr[1][:3]
		}
		/* #nosec */
		aesKey = md5.Sum([]byte(key + "ubx2022"))
		return fmt.Sprintf("%s", i)
	}
	w.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("[%-4s] >", i))
	}
	w.FormatCaller = func(i interface{}) string {
		s, ok := i.(string)
		if !ok {
			return fmt.Sprintf("%-10v >", i)
		}

		encrypt, err := encryption.AesCFBEncrypt([]byte(s), aesKey[:])
		if err != nil {
			fmt.Println(err)
			return fmt.Sprintf("%s >", i)
		}
		return fmt.Sprintf("%-10s >", base64.StdEncoding.EncodeToString(encrypt))

	}
	w.FormatMessage = func(i interface{}) string {
		s, ok := i.(string)
		if !ok {
			return fmt.Sprintf("%v >", i)
		}

		encrypt, err := encryption.AesCFBEncrypt([]byte(s), aesKey[:])
		if err != nil {
			fmt.Println(err)
			return fmt.Sprintf("%s >", i)
		}

		return fmt.Sprintf("%s >", base64.StdEncoding.EncodeToString(encrypt))

	}
	w.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("> [%s:", i)
	}
	w.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s]", i))
	}
	w.FormatErrFieldName = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s:", i))
	}
	w.FormatErrFieldValue = func(i interface{}) string {
		s, ok := i.(string)
		if !ok {
			return fmt.Sprintf("%v", i)
		}

		encrypt, err := encryption.AesCFBEncrypt([]byte(s), aesKey[:])
		if err != nil {
			fmt.Println(err)
			return fmt.Sprintf("%s", i)
		}

		return fmt.Sprintf("%s", base64.StdEncoding.EncodeToString(encrypt))

	}
	return w
}

// AddFieldHook todo test
type AddFieldHook struct {
}

// Run todo test
func (AddFieldHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {

	if level == zerolog.DebugLevel {
		e.Str("encryption", "true")

	}
}
