package database

import (
	"database/sql"
	"github.com/BlackRRR/first-tg-bot/cfg"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type User struct {
	userid   int    `json:"user_id"`
	username string `json:"user_name"`
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

func GetAllData(db *sql.DB) (Users []User, username []string) {
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
