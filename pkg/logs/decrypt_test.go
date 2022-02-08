package logs

import (
	"fmt"
	"testing"
)

func TestDecryptLog(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{filePath: "/home/honwee/.UBX/ubx-server.log/___go_build_ubx_server.20220106.log"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DecryptLog(tt.args.filePath)
			if err != nil {
				fmt.Println(err)
				return
			}
		})
	}
}

func TestDecryptLogToFile(t *testing.T) {
	type args struct {
		filePath string
		target   string
	}
	tests := []struct {
		name string
		args args
	}{
		{"test", args{
			filePath: "/home/honwee/.UBX/ubx-server.log/___go_build_ubx_server.20220106.log",
			target:   "/home/honwee/.UBX/ubx-server.log/test.log",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DecryptLogToFile(tt.args.filePath, tt.args.target)
			if err != nil {
				fmt.Println(err)
				return
			}

		})
	}
}
