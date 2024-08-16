package orm

import (
	"github.com/caixr9527/go-cloud/common/utils"
	"github.com/caixr9527/go-cloud/config"
	"github.com/caixr9527/go-cloud/factory"
	"github.com/caixr9527/go-cloud/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math"
	"sync"
	"time"
)

var once sync.Once

type Orm struct {
	Db *gorm.DB
}

func (o *Orm) Create() {
	configuration := factory.Get(&config.Configuration{})
	if !configuration.Db.Enable {
		return
	}
	once.Do(func() {
		t := configuration.Db.Type
		switch t {
		case "mysql":
		case "MYSQL":
			o.initMysqlConn()
		}
	})
}

func init() {
	factory.Register(&Orm{})
}

func (o *Orm) Order() int {
	return math.MinInt + 4
}

func (o *Orm) Name() string {
	return utils.ObjName(o)
}

func (o *Orm) Refresh() {

}

func (o *Orm) Destroy() {

}

func (o *Orm) initMysqlConn() {
	configuration := factory.Get(&config.Configuration{})
	logger := factory.Get(&log.Log{})
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
	o.Db = db
	factory.Create(o)
	logger.Info("init mysql conn success")
}
