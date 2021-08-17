package database

import (
	"database/sql"
	"github.com/BlackRRR/first-tg-bot/cfg"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var RatioChange *Ratio

type User struct {
	UserID       int64
	UserName     string
	Wins         int
	Losses       int
	WinLossRatio float64
	Language     string
}

type Ratio struct {
	Wins   int
	Losses int
	Ratio  float64
}

func DBConn() *sql.DB {
	db, err := sql.Open("mysql", cfg.DBCfg.User+cfg.DBCfg.Password+"@/"+cfg.DBCfg.Name)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func AddUser(db *sql.DB, userID int64, userName string, language string) {
	var (
		wins   = 0
		losses = 0
		ratio  = 0.0
	)
	insert, err := db.Prepare("INSERT INTO users (user_id, user_name,wins,losses,win_loss_ratio,language) VALUES (?,?,?,?,?,?)")
	if err != nil {
		log.Print(err)
	}

	_, err = insert.Exec(userID, userName, wins, losses, ratio, language)
	if err != nil {
		log.Println(err)
	}

	RatioChange = &Ratio{
		Wins:   wins,
		Losses: losses,
		Ratio:  ratio,
	}

	defer db.Close()
}

func GetAllData(db *sql.DB) (Users []User) {
	resp, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Println(err)
	}

	defer resp.Close()

	for resp.Next() {
		var user User
		err = resp.Scan(&user.UserID, &user.UserName, &user.Wins, &user.Losses, &user.WinLossRatio, &user.Language)
		if err != nil {
			log.Println(err)
		}

		Users = append(Users, user)
	}

	return Users
}

func UpdateTable(db *sql.DB, wins int, loss int, ratio float64, userID int64) {
	update, err := db.Prepare("UPDATE users SET wins=?,losses=?,win_loss_ratio=? where user_id=?;")
	if err != nil {
		log.Print(err)
	}

	_, err = update.Exec(wins, loss, ratio, userID)
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
}
