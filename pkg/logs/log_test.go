package logs

import (
	"errors"
	"testing"
	"time"
)

func TestNewLog(t *testing.T) {

	l := New(Config{
		Level:    "debug",
		FilePath: "test.log",
		Fields: map[string]string{
			"f1": "foo",
			"f2": "bar",
		},
		ToStdout: true,
	})

	ulog := l.Hook(AddFieldHook{})
	ulog.Info().Msg("test info")

	ulog.Info().Dur("dur", time.Second).Msg("test")
	ulog.Debug().Dur("dur", time.Second).Msg("test")
	ulog.Err(errors.New("err 测试的错误信息")).Msgf("xxx的结果为：%+v", true)
}

func TestMain(m *testing.M) {
	m.Run()
}
