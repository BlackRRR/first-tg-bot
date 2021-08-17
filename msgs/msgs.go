package msgs

import (
	"github.com/BlackRRR/first-tg-bot/language"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func MsgClickAnyCell(callbackID string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewCallback(callbackID, language.LangText(language.UserLang.Language, "open_flag"))
	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		log.Println(err)
	}
}

func MsgFlag(callbackID string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewCallback(callbackID, language.LangText(language.UserLang.Language, "already_flag"))
	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		log.Println(err)
	}
}

func MsgOutOfFlag(callbackID string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewCallback(callbackID, language.LangText(language.UserLang.Language, "out_of_flag"))
	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		log.Println(err)
	}
}

func SendMsgAll(userID int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(userID, language.LangText(language.UserLang.Language, "bot_started_press_sapper"))
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func SendWorkIsUnderWayMsg(ChatId int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(ChatId, language.LangText(language.UserLang.Language, "conducting_work"))
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
	return
}

func SendWorkIsUnderWayCallBack(CallbackId string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewCallback(CallbackId, language.LangText(language.UserLang.Language, "conducting_work"))
	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		log.Println(err)
	}
	return
}

func Rules(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, language.LangText(language.UserLang.Language, "rules"))
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func SelectDifficultGame(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, language.LangText(language.UserLang.Language, "take_level"))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "easy"), "easy"),
			tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "medium"), "medium"),
			tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "hard"), "hard"),
		))
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func TakeLanguage(userID int64, bot *tgbotapi.BotAPI) {
	msgLanguage := tgbotapi.NewMessage(userID, language.LangText(language.UserLang.Language, "select_language"))
	msgLanguage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "russian"), "ru"),
			tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "english"), "en")))
	_, err := bot.Send(msgLanguage)
	if err != nil {
		log.Println(err)
	}
}
