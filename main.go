package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Post struct {
	ID      int    `db:"id"`
	UserID  int    `db:"user_id"`
	Content string `db:"content"`
}

func InitDB() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", "./test.db")
	if err != nil {
		log.Fatalln(err)
	}

	// テーブルの作成
	schema := `
	DROP TABLE IF EXISTS users;
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);

	DROP TABLE IF EXISTS posts;
	CREATE TABLE posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func main() {
	db := InitDB()
	defer db.Close()
}
