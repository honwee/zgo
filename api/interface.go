package api

import (
	"zgo/api/v1"
	"zgo/api/v1beta1"
)

type Interface interface {
	v1() v1.Interface
	v1beta1() v1beta1.Interface
}

type group struct {
}

func New() Interface {
	return &group{}
}

func (g *group) v1() v1.Interface {
	return v1.New()
}

func (g *group) v1beta1() v1beta1.Interface {
	return v1beta1.New()
}
