package database

import (
	"github.com/BlackRRR/first-tg-bot/language"
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
	msg := tgbotapi.NewMessage(userID, language.LangText(language.UserLang.Language, "bot_started_press_sapper"))
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	msgLanguage := tgbotapi.NewMessage(userID, language.LangText(language.UserLang.Language, "select_language"))
	msgLanguage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "russian"), "ru"),
			tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "english"), "en")))
	_, err = bot.Send(msgLanguage)
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
