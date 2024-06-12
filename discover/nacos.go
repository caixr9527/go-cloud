package discover

import (
	"github.com/caixr9527/go-cloud/common/utils/stringUtils"
	"github.com/caixr9527/go-cloud/config"
	"github.com/caixr9527/go-cloud/log"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"sync"
)

var ConfigClient config_client.IConfigClient
var NamingClient naming_client.INamingClient
var once sync.Once

func Init() {
	if !config.Cfg.Discover.EnableDiscover && !config.Cfg.Discover.EnableConfig {
		return
	}
	once.Do(func() {
		createClient()
	})
}

func createClient() {
	clientConfig := clientConf()
	serverConfigs := serverConfig()
	if config.Cfg.Discover.EnableDiscover {
		namingClient, err := clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)
		if err != nil {
			log.Log.Error(err.Error())
			return
		}
		NamingClient = namingClient
	}

	if config.Cfg.Discover.EnableConfig {
		configClient, err := clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		)
		if err != nil {
			log.Log.Error(err.Error())
			return
		}

		ConfigClient = configClient
	}
}

func clientConf() constant.ClientConfig {
	cf := config.Cfg.Discover.Client
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

func serverConfig() []constant.ServerConfig {
	s := config.Cfg.Discover.Server
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
