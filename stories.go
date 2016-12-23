package main

import (
	"fmt"
	"time"

	"github.com/gocraft/dbr"
)

const tableName = "stories"

//StoriesTable 劇場話数テーブル
type StoriesTable struct {
	ID     int64     `db:"id"`
	Number int       `db:"number"`
	Title  string    `db:"title"`
	Uptime time.Time `db:"uptime"`
}

//InsertStory 劇場話数情報を追加する
func InsertStory(number int, title string) error {
	if number == 0 || title == "" {
		return fmt.Errorf("パラメーターエラー")
	}
	record := StoriesTable{Number: number, Title: title, Uptime: time.Now()}

	session := ConnectDB()
	if session == nil {
		return fmt.Errorf("DB接続失敗")
	}
	defer session.Close()

	selectRecords := selectStoryRecord(session, number)
	if selectRecords == nil || len(selectRecords) == 0 {
		_, err := session.InsertInto(tableName).
			Columns("number", "title", "uptime").
			Record(record).
			Exec()
		if err != nil {
			return err
		}
	} else {
		_, err := session.Update(tableName).
			Set("title", record.Title).
			Set("uptime", record.Uptime).
			Where("number = ?", record.Number).
			Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func selectStoryRecord(session *dbr.Session, number int) []StoriesTable {
	var recordList []StoriesTable
	_, err := session.Select("*").From(tableName).Where("number = ?", number).Load(&recordList)
	if err != nil {
		fmt.Printf("SelectStory err=" + err.Error())
		return nil
	}

	return recordList
}

//SelectStory 劇場話数を取得する
func SelectStory(number int) (result StoriesTable) {
	session := ConnectDB()
	if session == nil {
		fmt.Printf("SelectStory DB接続失敗")
		return
	}
	defer session.Close()

	var table []StoriesTable
	_, err := session.Select("*").From(tableName).Where("number = ?", number).Load(&table)
	if err != nil {
		fmt.Printf("SelectStory err=" + err.Error())
		return
	}

	result = table[0]
	return
}
