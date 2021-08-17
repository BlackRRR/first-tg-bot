package database

import (
	"github.com/BlackRRR/first-tg-bot/msgs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func CheckUsersFromDBAndSendMsg(users []User, bot *tgbotapi.BotAPI) {
	var user User
	for i := range users {
		user = users[i]
		msgs.SendMsgAll(user.UserID, bot)
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

func AddRatio(users []User, UserID int64) {
	var user User
	for i := range users {
		user = users[i]
		if user.UserID == UserID {
			RatioChange = &Ratio{
				Wins:   user.Wins,
				Losses: user.Losses,
				Ratio:  user.WinLossRatio,
			}
		}
	}
}
