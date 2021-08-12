package top

import (
	"fmt"
	"github.com/BlackRRR/first-tg-bot/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

var TextTop = "Номер  Имя игрока  победы  поражения  коэф. \n"

func CreateTop(users []database.User, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var user database.User
	for i := range users {
		user = users[i]
		AddToTop(i, user.UserName, user.Wins, user.Losses, user.WinLossRatio)
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, TextTop)
	if _, err := bot.Send(msg); err != nil {
		fmt.Println(err)
	}
}

func AddToTop(place int, userName string, wins int, losses int, ratio float64) {
	pl := strconv.Itoa(place + 1)
	win := strconv.Itoa(wins)
	loss := strconv.Itoa(losses)
	WLRatio := strconv.Itoa(int(ratio))
	TextTop = TextTop + pl + "." + "              " + userName + "             " + win + "               " + loss + "                 " + WLRatio + "\n"
}
