package discover

import (
	"github.com/caixr9527/go-cloud/common/utils/stringUtils"
	"github.com/caixr9527/go-cloud/component"
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/caixr9527/go-cloud/config"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"
	"reflect"
	"sync"
)

var once sync.Once

func init() {
	component.RegisterComponent(&discover{})
}

type discover struct {
}

func (d *discover) Order() int {
	return 4
}

func (d *discover) Create(s *component.Singleton) {
	configuration := factory.Get(config.Configuration{})
	if !configuration.Discover.EnableDiscover && !configuration.Discover.EnableConfig {
		return
	}
	once.Do(func() {
		createClient(s)
	})
}

func createClient(s *component.Singleton) {
	configuration := factory.Get(config.Configuration{})
	logger := factory.Get(&zap.Logger{})
	logger.Info("connect nacos")
	clientConfig := clientConf(configuration)
	serverConfigs := serverConfig(configuration)
	if configuration.Discover.EnableDiscover {
		namingClient, err := clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		s.Register(reflect.TypeOf(namingClient).Name(), namingClient)
	}

	if configuration.Discover.EnableConfig {
		configClient, err := clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		s.Register(reflect.TypeOf(configClient).Name(), configClient)
	}
	logger.Info("connect nacos success")
}

func clientConf(configuration config.Configuration) constant.ClientConfig {
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

func serverConfig(configuration config.Configuration) []constant.ServerConfig {
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
