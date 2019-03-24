package commons

import (
	"github.com/Unknwon/goconfig"
	"github.com/cihub/seelog"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

func DBApi(config *goconfig.ConfigFile) *xorm.Engine {
	var engine *xorm.Engine
	var err error
	engine, err = xorm.NewEngine("sqlite3", "smartcar.db")
	if err = engine.Ping(); err != nil {
		seelog.Error(err)
	}
	showORNot, _ := config.Bool("sqlite3", "show_sql")
	engine.ShowSQL(showORNot)
	idleNum, _ := config.Int("sqlite3", "idle_num")
	engine.SetMaxIdleConns(idleNum)
	openNum, _ := config.Int("sqlite3", "open_num")
	engine.SetMaxOpenConns(openNum)
	engine.SetMapper(core.GonicMapper{})
	return engine
}
