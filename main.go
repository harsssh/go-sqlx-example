package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// camel case でないならタグは不要
type User struct {
	ID   int
	Name string
}

type Post struct {
	ID      int
	UserID  int `db:"user_id"`
	Content string
}

func main() {
	db := InitDB()
	defer db.Close()

	BulkInsert(db)
	SelectUsers(db)
	InQuery(db)
	JoinQuery(db)
}

func InitDB() *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", "./test.db")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to the database")

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
	log.Println("Created tables")

	return db
}

func BulkInsert(db *sqlx.DB) {
	users := []User{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Charlie"},
	}
	result, err := db.NamedExec("INSERT INTO users (name) VALUES (:name)", users)
	if err != nil {
		log.Fatalln(err)
	}
	rowsAffected, err := result.RowsAffected()
	log.Printf("Insert users: %d\n", rowsAffected)

	// Alice has 2 posts, Bob has 1 post, Charlie has no post
	posts := []Post{
		{UserID: 1, Content: "Hello, Alice"},
		{UserID: 1, Content: "Nice to meet you"},
		{UserID: 2, Content: "Hello, Bob"},
	}
	result, err = db.NamedExec("INSERT INTO posts (user_id, content) VALUES (:user_id, :content)", posts)
	if err != nil {
		log.Fatalln(err)
	}
	rowsAffected, err = result.RowsAffected()
	log.Printf("Insert posts: %d\n", rowsAffected)
}

func SelectUsers(db *sqlx.DB) {
	users := []User{}
	err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("All users:", users)
}

func InQuery(db *sqlx.DB) {
	userIDs := []int{1, 2}
	query, args, _ := sqlx.In("SELECT * FROM users WHERE id IN (?)", userIDs)
	query = db.Rebind(query)

	var users []User
	if err := db.Select(&users, query, args...); err != nil {
		log.Fatalln(err)
	}

	log.Println("Selected users:", users)
}

func JoinQuery(db *sqlx.DB) {
	// users.id, posts.id のタグが被るので, 少なくとも一方のタグは必須
	// マッピング先が一意ならタグ, AS は不要
	type T struct {
		User `db:"user"`
		Post
	}

	// LEFT JOIN だと NULL をマッピングできなくてエラーになる
	// *Post を埋め込んでもダメ
	// refs: https://github.com/jmoiron/sqlx/issues/162
	query := `
		SELECT
			users.id AS "user.id",
			users.name AS "user.name",
			posts.*
		FROM users
		INNER JOIN posts ON users.id = posts.user_id
	`
	var result []T
	if err := db.Select(&result, query); err != nil {
		log.Fatalln(err)
	}

	log.Println("Joined result:", result)
}
