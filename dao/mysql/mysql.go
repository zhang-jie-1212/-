package mysql

import (
	"bluebell/settings"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

//sqlx连接数据库
func Init(cfg *settings.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		//zap.Error(err)，将err转为键值对 field形式；前面的参数是message信息
		zap.L().Error("connect mysql failed", zap.Error(err))
	}
	//设置最大连接数
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	//设置最大空闲连接数
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return err
}
func Close() {
	db.Close()
}
