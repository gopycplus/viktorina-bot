package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var gameDb *sql.DB
var questionDb *sql.DB

func ConnectGameDb() bool {
	var err error
	gameDb, err = sql.Open("sqlite3", "../game.db")
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func ConnectHistoryDb() bool {
	var err error
	questionDb, err = sql.Open("sqlite3", "../history.db")
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
