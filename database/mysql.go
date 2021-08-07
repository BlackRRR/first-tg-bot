package database

import (
	"database/sql"
	"github.com/BlackRRR/first-tg-bot/cfg"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type User struct { //TODO: use the `json` tag only if you are going to transfer it somewhere, you don't need to transfer it anywhere, so you can remove it
	userid   int    `json:"user_id"`   //TODO: UserID	if the structure field starts with a small letter it is visible only inside the folder
	username string `json:"user_name"` //TODO: UserName	the first letters of different words must begin with capital letters
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

func GetAllData(db *sql.DB) (Users []User, username []string) { //TODO: you have an array of users that already contains username, why pass them separately?
	resp, err := db.Query("SELECT * FROM `users`")
	if err != nil {
		log.Println(err)
	}

	defer resp.Close()

	for resp.Next() {
		var user User
		var UserName string
		err = resp.Scan(&user.userid, &user.username)
		if err != nil {
			log.Println(err)
		}

		UserName = user.username
		username = append(username, UserName)
		Users = append(Users, user)
	}

	return Users, username
}
