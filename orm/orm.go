package orm

import (
	"github.com/caixr9527/go-cloud/component"
	"github.com/caixr9527/go-cloud/component/factory"
	"github.com/caixr9527/go-cloud/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math"
	"reflect"
	"sync"
	"time"
)

var once sync.Once

type orm struct {
}

func (o *orm) Create(s *component.Singleton) {
	configuration := factory.Get(config.Configuration{})
	if !configuration.Db.Enable {
		return
	}
	once.Do(func() {
		t := configuration.Db.Type
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
	configuration := factory.Get(config.Configuration{})
	logger := factory.Get(&zap.Logger{})
	logger.Info("init mysql conn")
	mysqlCfg := configuration.Db.Mysql
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
	s.Register(reflect.TypeOf(db).String(), db)
	logger.Info("init mysql conn success")
}
