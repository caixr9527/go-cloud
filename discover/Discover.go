package discover

import (
	"fmt"
	"github.com/caixr9527/go-cloud/common/utils"
	"github.com/caixr9527/go-cloud/common/utils/stringUtils"
	"github.com/caixr9527/go-cloud/component"
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/caixr9527/go-cloud/config"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"math"
	"sync"
)

var once sync.Once

func init() {
	component.RegisterComponent(&Discover{})
}

type Discover struct {
	NamingClient naming_client.INamingClient
	ConfigClient config_client.IConfigClient
}

func (d *Discover) Refresh() {

}

func (d *Discover) Destroy() {

}

func (d *Discover) Order() int {
	return math.MinInt + 1
}

func (d *Discover) Name() string {
	return utils.ObjName(d)
}

func (d *Discover) Create() {
	configuration := factory.Get(&config.Configuration{})
	if !configuration.Discover.EnableDiscover && !configuration.Discover.EnableConfig {
		return
	}
	once.Do(func() {
		log.Println("connect nacos")
		d.createClient()
		log.Println("connect nacos success")
	})
}

func (d *Discover) createClient() {
	configuration := factory.Get(&config.Configuration{})
	clientConfig := d.clientConf(configuration)
	serverConfigs := d.serverConfig(configuration)
	if configuration.Discover.EnableDiscover {
		namingClient, err := clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		d.NamingClient = namingClient
	}

	if configuration.Discover.EnableConfig {
		configClient, err := clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		d.ConfigClient = configClient
		factory.Create(d)
	}

}

func (d *Discover) clientConf(configuration *config.Configuration) constant.ClientConfig {
	cf := configuration.Discover.Client
	clientConfig := constant.ClientConfig{}
	if cf.TimeoutMs != 0 {
		clientConfig.TimeoutMs = cf.TimeoutMs
	}
	if stringUtils.IsNotBlank(cf.NamespaceId) {
		clientConfig.NamespaceId = cf.NamespaceId
	}
	if stringUtils.IsNotBlank(cf.Endpoint) {
		clientConfig.Endpoint = cf.Endpoint
	}
	if stringUtils.IsNotBlank(cf.RegionId) {
		clientConfig.RegionId = cf.RegionId
	}
	if stringUtils.IsNotBlank(cf.AccessKey) {
		clientConfig.AccessKey = cf.AccessKey
	}
	if stringUtils.IsNotBlank(cf.SecretKey) {
		clientConfig.SecretKey = cf.SecretKey
	}
	clientConfig.OpenKMS = cf.OpenKMS
	if stringUtils.IsNotBlank(cf.CacheDir) {
		clientConfig.CacheDir = cf.CacheDir
	}
	if cf.UpdateThreadNum != 0 {
		clientConfig.UpdateThreadNum = cf.UpdateThreadNum
	}
	clientConfig.NotLoadCacheAtStart = cf.NotLoadCacheAtStart
	clientConfig.UpdateCacheWhenEmpty = cf.UpdateCacheWhenEmpty
	if stringUtils.IsNotBlank(cf.Username) {
		clientConfig.Username = cf.Username
	}
	if stringUtils.IsNotBlank(cf.Password) {
		clientConfig.Password = cf.Password
	}
	if stringUtils.IsNotBlank(cf.LogDir) {
		clientConfig.LogDir = cf.LogDir
	}
	if stringUtils.IsNotBlank(cf.LogLevel) {
		clientConfig.LogLevel = cf.LogLevel
	}
	return clientConfig
}

func (d *Discover) serverConfig(configuration *config.Configuration) []constant.ServerConfig {
	s := configuration.Discover.Server
	var serverConfigs []constant.ServerConfig
	for index := range s {
		sf := s[index]
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			Scheme:      sf.Scheme,
			ContextPath: sf.ContextPath,
			IpAddr:      sf.IpAddr,
			Port:        sf.Port,
			GrpcPort:    sf.GrpcPort,
		})
	}
	return serverConfigs
}
