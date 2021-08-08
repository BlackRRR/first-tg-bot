package database

import (
	"database/sql"
	"github.com/BlackRRR/first-tg-bot/cfg"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type User struct {
	UserID   int
	UserName string
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
