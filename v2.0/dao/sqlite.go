package dao

import (
	"database/sql"
	"log"

	//go-sqlite3 驱动
	_ "github.com/mattn/go-sqlite3"
)

//DB sqlite连接
var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "./web.db")
	if err != nil {
		log.Fatal(err)
	}
	//建表
	sqlStmt := `create table if not exists app (
		id text not null primary key,
		name text,
		version text,
		logpath text,
		confpath text,
		enable text,
		cmd_line text,
		work_dir text,
		pid integer not null default 0, 
		mem_percent float not null default 0,
		cpu_percent float not null default 0,
		running integer not null default 0 
		);
	`
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}
