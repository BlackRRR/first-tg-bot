package database

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func CheckUsersFromDBAndSendMsg(users []User, bot *tgbotapi.BotAPI) {
	var user User
	for i := range users {
		user = users[i]
		SendMsgAll(user.UserID, bot)
	}
}

func SendMsgAll(userID int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(userID, "Бот запущен \U0001F973, напишите /sapper чтобы начать игру")
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func CheckIdenticalValues(update *tgbotapi.Update, users []User) bool {
	var user User
	var flag = true
	for i := range users {
		user = users[i]
		if user.UserID == update.Message.Chat.ID {
			flag = false
		}
	}
	return flag
}
