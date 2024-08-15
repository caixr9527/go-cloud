package service

import "github.com/caixr9527/go-cloud/transport"

type GoodsService struct {
	GetUser func(id int) any `method:"GET"`
}

func (g *GoodsService) Client() transport.Client {
	return transport.Client{}
}

func (g *GoodsService) Order() int {
	return 0
}

func (g *GoodsService) Name() string {
	return g.Client().ContentId
}
