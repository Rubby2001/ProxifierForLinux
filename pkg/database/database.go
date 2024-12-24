package database

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"xorm.io/xorm"
)

var Engine *xorm.Engine
var DatabaseName string = "./rules.db"

type Rules struct {
	Type     string `json:"type"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func ConnectDatabase() {
	var err error
	Engine, err = xorm.NewEngine("sqlite3", DatabaseName)
	if err != nil {
		log.Fatalf("连接sqlite3数据库失败: %v", err)
	}
	err = Engine.Sync2(new(Rules))
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
}
