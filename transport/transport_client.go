package transport

import "github.com/caixr9527/go-cloud/common"

type Fallback interface {
}
type Configuration interface {
}

type Client struct {
	Protocol      string // http,rpc
	ServiceName   string // service name
	ContentId     string // content id
	Url           string
	Path          string
	Fallback      Fallback
	Configuration Configuration
}

type Service interface {
	common.Bean
	Client() Client
}

// 复用component.go中的Pool
// factory统一的create和get出入口
