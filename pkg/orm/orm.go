package orm

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
)

const (
	DefaultMaxOpenConn = 100
	DefaultMaxIdleConn = 10
	DefaultMaxLifetime = 3600
)

type Config struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  int
	LogLevel     logger.LogLevel
}

type ormLog struct {
	LogLevel logger.LogLevel
}

func (l *ormLog) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l

}

func (l *ormLog) Info(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel < logger.Info {
		return
	}
	logx.WithContext(ctx).Infof(s, i...)
}

func (l *ormLog) Warn(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel < logger.Warn {
		return
	}
	logx.WithContext(ctx).Errorf(s, i...)
}

func (l *ormLog) Error(ctx context.Context, s string, i ...interface{}) {
	if l.LogLevel < logger.Error {
		return
	}
	logx.WithContext(ctx).Errorf(s, i...)
}

func (l *ormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	logx.WithContext(ctx).WithDuration(elapsed).Infof("[%.3fms] [rows: %v] %s", float64(elapsed.Nanoseconds())/1e6, rows, sql)
}

func NewMysql(conf *Config) (*gorm.DB, error) {
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = DefaultMaxIdleConn
	}
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = DefaultMaxOpenConn
	}
	if conf.MaxLifetime == 0 {
		conf.MaxLifetime = DefaultMaxLifetime
	}

	db, err := gorm.Open(mysql.Open(conf.DSN), &gorm.Config{
		NamingStrategy: &schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true, // 禁用表名复数
		},
		//Logger: logger.Default.LogMode(logger.Info),
		Logger: &ormLog{conf.LogLevel},
	})
	if err != nil {
		return nil, err
	}
	sdb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sdb.SetMaxIdleConns(conf.MaxIdleConns)
	sdb.SetMaxOpenConns(conf.MaxOpenConns)
	sdb.SetConnMaxLifetime(time.Second * time.Duration(conf.MaxLifetime))

	//err = db.Use(NewCustomPlugin())
	//if err != nil {
	//	return nil, err
	//}
	err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics()))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MustNewMysql(conf *Config) *gorm.DB {
	db, err := NewMysql(conf)
	logx.Must(err)
	return db
}
