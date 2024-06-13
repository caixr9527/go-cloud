package orm

import (
	"github.com/caixr9527/go-cloud/component"
	"github.com/caixr9527/go-cloud/config"
	logger "github.com/caixr9527/go-cloud/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math"
	"sync"
	"time"
)

var DB *gorm.DB
var once sync.Once

type orm struct {
}

func (o *orm) Create(s *component.Singleton) {
	if !config.Cfg.Db.Enable {
		return
	}
	once.Do(func() {
		t := config.Cfg.Db.Type
		switch t {
		case "mysql":
		case "MYSQL":
			initMysqlConn(s)
		}
	})
}

func init() {
	component.RegisterComponent(&orm{})
}

func (o *orm) Order() int {
	return math.MinInt + 2
}

func initMysqlConn(s *component.Singleton) {
	logger.Log.Info("init mysql conn")
	mysqlCfg := config.Cfg.Db.Mysql
	g := &gorm.Config{
		PrepareStmt:                              mysqlCfg.PrepareStmt,
		SkipDefaultTransaction:                   mysqlCfg.SkipDefaultTransaction,
		FullSaveAssociations:                     mysqlCfg.FullSaveAssociations,
		DryRun:                                   mysqlCfg.DryRun,
		DisableAutomaticPing:                     mysqlCfg.DisableAutomaticPing,
		DisableForeignKeyConstraintWhenMigrating: mysqlCfg.DisableForeignKeyConstraintWhenMigrating,
		IgnoreRelationshipsWhenMigrating:         mysqlCfg.IgnoreRelationshipsWhenMigrating,
		DisableNestedTransaction:                 mysqlCfg.DisableNestedTransaction,
		AllowGlobalUpdate:                        mysqlCfg.AllowGlobalUpdate,
		QueryFields:                              mysqlCfg.QueryFields,
		TranslateError:                           mysqlCfg.TranslateError,
	}
	if mysqlCfg.CreateBatchSize != 0 {
		g.CreateBatchSize = mysqlCfg.CreateBatchSize
	}
	db, err := gorm.Open(mysql.Open(mysqlCfg.Url), g)
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
		conn.SetConnMaxLifetime(time.Duration(mysqlCfg.MaxLifetime))
	}
	if mysqlCfg.MaxIdleConns != 0 {
		conn.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	}
	if mysqlCfg.MaxIdleTime != 0 {
		conn.SetConnMaxIdleTime(time.Duration(mysqlCfg.MaxIdleTime))
	}
	s.Register("db", db)
	logger.Log.Info("init mysql conn success")
}
