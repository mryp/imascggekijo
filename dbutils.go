package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" //dbrで使用する
	"github.com/gocraft/dbr"
)

//DBデータ
const (
	DBUserID     = "root"
	DBPassword   = "root"
	DBHostName   = "127.0.0.1"
	DBPortNumber = "3306"
	DBName       = "imascggekijo"
)

//ConnectDB DB接続
func ConnectDB() *dbr.Session {
	db, err := dbr.Open("mysql", DBUserID+":"+DBPassword+"@tcp("+DBHostName+":"+DBPortNumber+")/"+DBName+"?parseTime=true", nil)
	if err != nil {
		fmt.Printf("connectDB err=%v\n", err)
		return nil
	}

	dbsession := db.NewSession(nil)
	return dbsession
}
