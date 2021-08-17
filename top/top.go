package top

import (
	"fmt"
	"github.com/BlackRRR/first-tg-bot/database"
	"github.com/BlackRRR/first-tg-bot/language"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"sort"
	"strconv"
)

var (
	TextTop string
	Place   UserPlace
)

type UserPlace struct {
	UserName     string
	Wins         int
	Losses       int
	WinLossRatio float64
}

func CreateTop(users []database.User, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	TextTop = language.LangText(language.UserLang.Language, "text_top")
	var user database.User
	var masPlace []UserPlace
	for i := range users {
		user = users[i]

		place := ReturnUserPlace(user)

		Place = place
		masPlace = append(masPlace, Place)
	}

	sort.SliceStable(masPlace, func(i, j int) bool {
		return masPlace[i].WinLossRatio > masPlace[j].WinLossRatio
	})

	AddToTop(masPlace)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, TextTop)
	if _, err := bot.Send(msg); err != nil {
		fmt.Println(err)
	}
}

func ReturnUserPlace(user database.User) UserPlace {
	place := UserPlace{
		UserName:     user.UserName,
		Wins:         user.Wins,
		Losses:       user.Losses,
		WinLossRatio: user.WinLossRatio,
	}
	return place
}

func AddToTop(masPlace []UserPlace) {
	for i := range masPlace {
		var user UserPlace
		user = masPlace[i]

		AddToText(i, user.UserName, user.Wins, user.Losses, user.WinLossRatio)
	}
}

func AddToText(place int, userName string, wins int, losses int, ratio float64) {
	var spaces string
	pl := strconv.Itoa(place + 1)
	win := strconv.Itoa(wins)
	loss := strconv.Itoa(losses)
	WLRatio := fmt.Sprintf("%.3f", ratio)

	for i := 0; i < (21 - len(userName)); i++ {
		spaces += " "
	}

	TextTop = TextTop + pl + "." + "   " + userName + spaces + win + "        " + loss + "      " + WLRatio + "\n"
}
