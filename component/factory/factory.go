package factory

import (
	"github.com/caixr9527/go-cloud/component"
	"github.com/caixr9527/go-cloud/config"
	"gorm.io/gorm"
)

func GetConf() config.Configuration {
	return component.SinglePool.Get("config").(config.Configuration)
}

//func GetRedis() redis.Redis {
//	return component.SinglePool.Get("redis").(redis.Redis)
//}

//func GetLogger() zap.Logger {
//	return component.SinglePool.Get("logger").(zap.Logger)
//}

//func GetConfigClient() config_client.IConfigClient {
//	return component.SinglePool.Get("configClient").(config_client.IConfigClient)
//}
//func GetNamingClient() naming_client.INamingClient {
//	return component.SinglePool.Get("namingClient").(naming_client.INamingClient)
//}

func GetDb() gorm.DB {
	return component.SinglePool.Get("db").(gorm.DB)
}
