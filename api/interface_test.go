package api

import (
	"testing"

	_ "zgo/internal/boot"
)

func TestMain(m *testing.M) {

	New().v1beta1().Create()
}
