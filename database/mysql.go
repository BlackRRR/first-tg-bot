package database

import (
	"database/sql"
	"github.com/BlackRRR/first-tg-bot/cfg"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type User struct { // use the `json` tag only if you are going to transfer it somewhere, you don't need to transfer it anywhere, so you can remove it
	UserID   int    // UserID	if the structure field starts with a small letter it is visible only inside the folder
	UserName string // UserName	the first letters of different words must begin with capital letters
}

func DBConn() *sql.DB {
	db, err := sql.Open("mysql", cfg.DBCfg.User+cfg.DBCfg.Password+"@/"+cfg.DBCfg.Name)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func AddDB(db *sql.DB, userID int, userName string) {
	insert, err := db.Prepare("INSERT INTO users (user_id, user_name) VALUES (?,?)")
	if err != nil {
		log.Print(err)
	}

	_, err = insert.Exec(userID, userName)
	if err != nil {
		log.Println(err)
	}

	defer db.Close()
}

func GetAllData(db *sql.DB) (Users []User) { // you have an array of users that already contains UserName, why pass them separately?
	resp, err := db.Query("SELECT * FROM `users`")
	if err != nil {
		log.Println(err)
	}

	defer resp.Close()

	for resp.Next() {
		var user User
		err = resp.Scan(&user.UserID, &user.UserName)
		if err != nil {
			log.Println(err)
		}

		Users = append(Users, user)
	}

	return Users
}
