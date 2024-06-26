package discover

import (
	"fmt"
	"github.com/caixr9527/go-cloud/common/utils"
	"github.com/caixr9527/go-cloud/common/utils/stringUtils"
	"github.com/caixr9527/go-cloud/component"
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/caixr9527/go-cloud/config"
	logger "github.com/caixr9527/go-cloud/log"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v3"
	"log"
	"math"
	"strings"
	"sync"
)

var once sync.Once

func init() {
	component.RegisterComponent(&Discover{})
}

type Discover struct {
	naming_client.INamingClient
	config_client.IConfigClient
}

func (d *Discover) Refresh() {
	d.refreshConfig()
	d.registerInstance()
}

func (d *Discover) Destroy() {
	l := factory.Get(&logger.Log{})
	d.deregister()
	if d.IConfigClient != nil {
		d.IConfigClient.CloseClient()
		l.Info("destory nacos config client")
	}
	if d.INamingClient != nil {
		d.INamingClient.CloseClient()
		l.Info("destory nacos naming client")
	}
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
		d.INamingClient = namingClient
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
		d.IConfigClient = configClient
	}
	if d.IConfigClient == nil && d.INamingClient == nil {
		return
	}
	factory.Create(d)
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

func (d *Discover) registerInstance() {
	success, err := d.RegisterInstance(d.getRegisterConfig())
	if err != nil {
		log.Println(err)
		return
	}
	if success {
		log.Println("register instance success")
	}
}

func (d *Discover) getRegisterConfig() vo.RegisterInstanceParam {
	param := vo.RegisterInstanceParam{}
	configuration := factory.Get(&config.Configuration{})
	discover := configuration.Discover.Discover
	if stringUtils.IsNotBlank(discover.Ip) {
		param.Ip = discover.Ip
	} else {
		param.Ip = utils.GetRealIp()
	}
	if discover.Port != 0 {
		param.Port = discover.Port
	} else {
		param.Port = uint64(configuration.Server.Port)
	}
	if discover.Weight != 0 {
		param.Weight = discover.Weight
	} else {
		param.Weight = 1
	}
	param.Enable = discover.Enable
	param.Healthy = discover.Healthy
	param.Metadata = discover.Metadata
	if stringUtils.IsNotBlank(discover.ClusterName) {
		param.ClusterName = discover.ClusterName
	}
	if stringUtils.IsNotBlank(discover.ServiceName) {
		param.ServiceName = discover.ServiceName
	} else {
		param.ServiceName = configuration.Server.ServerName
	}
	if stringUtils.IsNotBlank(discover.GroupName) {
		param.GroupName = discover.GroupName
	}
	param.Ephemeral = true
	return param
}

func (d *Discover) deregister() {
	success, err := d.DeregisterInstance(d.getDeregisterConfig())
	l := factory.Get(&logger.Log{})
	if err != nil {
		l.Error(err.Error())
	}
	if success {
		l.Info("deregister instance success")
	}
}

func (d *Discover) getDeregisterConfig() vo.DeregisterInstanceParam {
	param := vo.DeregisterInstanceParam{}
	configuration := factory.Get(&config.Configuration{})
	discover := configuration.Discover.Discover
	if stringUtils.IsNotBlank(discover.Ip) {
		param.Ip = discover.Ip
	} else {
		param.Ip = utils.GetRealIp()
	}
	if discover.Port != 0 {
		param.Port = discover.Port
	} else {
		param.Port = uint64(configuration.Server.Port)
	}
	param.Cluster = discover.ClusterName
	if stringUtils.IsNotBlank(discover.ServiceName) {
		param.ServiceName = discover.ServiceName
	} else {
		param.ServiceName = configuration.Server.ServerName
	}
	param.GroupName = discover.GroupName
	param.Ephemeral = discover.Ephemeral
	return param
}

func (d *Discover) refreshConfig() {
	configuration := factory.Get(&config.Configuration{})
	l := factory.Get(&logger.Log{})
	if !configuration.Discover.EnableConfig {
		return
	}
	dataIds := strings.Split(configuration.Discover.Config.DataIds, ",")
	group := configuration.Discover.Config.Group
	var contents strings.Builder
	for index := range dataIds {
		dataId := dataIds[index]
		content, err := d.GetConfig(vo.ConfigParam{
			DataId: dataId,
			Group:  group,
		})
		if err != nil {
			l.Error(err.Error())
		} else if content != "" {
			contents.WriteString(content)
			contents.WriteString("\n")
			contents.WriteString("---")
		}
	}
	newContents := contents.String()
	if newContents != "" {
		err := yaml.Unmarshal([]byte(newContents), &configuration)
		if err != nil {
			l.Error(err.Error())
			return
		}
		factory.Create(configuration)
		configuration.LoadRemoteCustomConfig(newContents)
	}

	if configuration.Discover.Config.Refresh {
		for index := range dataIds {
			dataId := dataIds[index]
			go func() {
				l.Info(fmt.Sprintf("listening config group: %s, dataId: %s", group, dataId))
				err := d.IConfigClient.ListenConfig(vo.ConfigParam{
					DataId: dataId,
					Group:  group,
					OnChange: func(namespace, group, dataId, data string) {
						l.Info(fmt.Sprintf("config change namespace: %s, group: %s, dataId: %s", namespace, group, dataId))
						err := yaml.Unmarshal([]byte(data), &configuration)
						if err != nil {
							l.Error(err.Error())
							return
						}
						factory.Create(configuration)
						configuration.LoadRemoteCustomConfig(data)
					},
				})
				if err != nil {
					l.Error(err.Error())
				}
			}()
		}
	}
}
