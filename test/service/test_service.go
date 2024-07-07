package service

import "github.com/caixr9527/go-cloud/transport"

type GoodsService struct {
	GetUser func(id int) any
}

func (g *GoodsService) Client() transport.Client {
	return transport.Client{}
}

func (g *GoodsService) Register() {

}
