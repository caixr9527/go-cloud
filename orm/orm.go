package orm

import (
	"github.com/caixr9527/go-cloud/config"
	logger "github.com/caixr9527/go-cloud/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var DB *gorm.DB
var once sync.Once

func InitDbConn() {
	if !config.Cfg.Db.Enable {
		return
	}
	once.Do(func() {
		t := config.Cfg.Db.Type
		switch t {
		case "mysql":
		case "MYSQL":
			initMysqlConn()
		}
	})
}

func initMysqlConn() {
	logger.Log.Info("init mysql conn")
	mysqlCfg := config.Cfg.Db.Mysql
	db, err := gorm.Open(mysql.Open(mysqlCfg.Url), &gorm.Config{
		PrepareStmt: mysqlCfg.PrepareStmt,
	})
	if err != nil {
		panic(err)
	}
	conn, err := db.DB()
	if err != nil {
		panic(err)
	}
	if mysqlCfg.MaxOpenConns != 0 {
		conn.SetMaxOpenConns(mysqlCfg.MaxOpenConns)
	}
	if mysqlCfg.MaxLifetime != 0 {
		conn.SetConnMaxLifetime(time.Duration(mysqlCfg.MaxLifetime) * time.Millisecond)
	}
	if mysqlCfg.MaxIdleConns != 0 {
		conn.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	}
	if mysqlCfg.MaxIdleTime != 0 {
		conn.SetConnMaxIdleTime(time.Duration(mysqlCfg.MaxIdleTime) * time.Millisecond)
	}
	DB = db
	logger.Log.Info("init mysql conn success")
}
