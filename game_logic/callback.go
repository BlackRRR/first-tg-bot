package game_logic

import (
	"database/sql"
	"fmt"
	"github.com/BlackRRR/first-tg-bot/assets"
	"github.com/BlackRRR/first-tg-bot/database"
	"github.com/BlackRRR/first-tg-bot/language"
	"github.com/BlackRRR/first-tg-bot/models"
	"github.com/BlackRRR/first-tg-bot/msgs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

func ActionWithCallback(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, users []database.User, db *sql.DB) {
	switch callback.Data {
	case "start":
		SendAdminBotStart(callback.ID, bot, callback.From.ID)
		database.CheckUsersFromDBAndSendMsg(users, bot)
	case "turn on":
		CheckDeveloperMode(callback, bot, true)
	case "turn off":
		CheckDeveloperMode(callback, bot, false)
	case "5", "6", "7", "8":
		TakeCallBackFieldSize(callback, bot, callback.Data, models.Level)
	case "ru":
		language.ReturnLanguage("ru")
		SendMsgLanguage(callback, bot, "ru")
	case "en":
		language.ReturnLanguage("en")
		SendMsgLanguage(callback, bot, "en")
	case "easy":
		ReturnLevelGameSendMsg("easy", callback.ID, bot)
		TakeFieldSize(callback, bot)
	case "medium":
		ReturnLevelGameSendMsg("medium", callback.ID, bot)
		TakeFieldSize(callback, bot)
	case "hard":
		ReturnLevelGameSendMsg("hard", callback.ID, bot)
		TakeFieldSize(callback, bot)
	default:
		HandlingGameLogic(callback, bot, users, db)
	}
}

func HandlingGameLogic(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, users []database.User, db *sql.DB) {

	key, i, j := DataSplit(callback.Data)

	if assets.Games[key] == nil {
		SendGameOver(callback, bot)
		return
	}

	if assets.Games[key].FlagCounter == 0 && assets.Games[key].Flag == "true" && i != assets.Games[key].Size+1 {
		msgs.MsgOutOfFlag(callback.ID, bot)
		return
	}

	if i == assets.Games[key].Size+1 && assets.Games[key].Flag == "false" {
		assets.Games[key].Flag = "true"
		ReEditForFlags(callback, bot, key)
		return
	}

	if assets.Games[key].Flag == "true" && i == assets.Games[key].Size+1 {
		assets.Games[key].Flag = "false"
		ReEditForFlags(callback, bot, key)
		return
	}

	if assets.Games[key].Flag == "true" && assets.Games[key].OpenedButtonsField[i][j] == "flag" {
		if assets.Games[key].FlagCounter != models.FlagCounter {
			assets.Games[key].FlagCounter += 1
		}
		assets.Games[key].OpenedButtonsField[i][j] = "false"
		ReEditForFlags(callback, bot, key)
		return
	}

	if i != assets.Games[key].Size+1 && assets.Games[key].Flag == "true" {
		if assets.Games[key].FlagCounter != 0 {
			assets.Games[key].FlagCounter -= 1
		}
		assets.Games[key].OpenedButtonsField[i][j] = "flag"
		ReEditForFlags(callback, bot, key)
		return
	}

	if i == assets.Games[key].Size+1 && assets.Games[key].Flag == "" {
		msgs.MsgClickAnyCell(callback.ID, bot)
		return
	}

	if _, exist := assets.Games[key]; !exist {
		return
	}

	if assets.Games[key].PlayingField[i][j] == " " {
		return
	}

	if assets.Games[key].PlayingField[i][j] == "bomb" && Counter(key) == 0 {
		assets.Games[key].Flag = "false"
		fmt.Println(assets.Games[key].FlagCounter)
		ReEditField(callback, bot, key)
		ActionWithCallback(callback, bot, users, db)
		assets.Games[key].FlagCounter = models.FlagCounter
		ReEditForFlags(callback, bot, key)
		return
	}

	if assets.Games[key].OpenedButtonsField[i][j] == "true" {
		return
	}

	if assets.Games[key].OpenedButtonsField[i][j] == "flag" {
		msgs.MsgFlag(callback.ID, bot)
		return
	}

	assets.Games[key].OpenedButtonsField[i][j] = "true"

	switch assets.Games[key].PlayingField[i][j] {
	case "0":
		OpenZero(i, j, key)
	case "bomb":
		OpenAllBombsAfterWin(key)
		ActionsWithBombUpdate(i, j, key, callback, bot, db)
		return
	}

	if Counter(key) == assets.Games[key].Size*assets.Games[key].Size-assets.Games[key].BombCounter {
		OpenAllBombsAfterWin(key)
		ActionsWithWin(key, callback, bot, db)
		return
	}

	if Counter(key) == 0 {
		assets.Games[key].FlagCounter = models.FlagCounter
	}

	assets.Games[key].Flag = "false"
	CallEditMessage(key, callback.Message.Chat.ID, callback.Message.MessageID, bot)

}

func TakeCallBackFieldSize(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, size string, level map[int]int) {

	bombCounter := level

	intSize, _ := strconv.Atoi(size)
	key := GenerateField(intSize, bombCounter[intSize])

	switch models.Difficult {
	case "easy":
		switch size {
		case "5":
			assets.Games[key].FlagCounter = 15
			models.FlagCounter = assets.Games[key].FlagCounter
		case "6":
			assets.Games[key].FlagCounter = 20
			models.FlagCounter = assets.Games[key].FlagCounter
		case "7":
			assets.Games[key].FlagCounter = 25
			models.FlagCounter = assets.Games[key].FlagCounter
		case "8":
			assets.Games[key].FlagCounter = 35
			models.FlagCounter = assets.Games[key].FlagCounter
		}
	case "medium":
		switch size {
		case "5":
			assets.Games[key].FlagCounter = 10
			models.FlagCounter = assets.Games[key].FlagCounter
		case "6":
			assets.Games[key].FlagCounter = 12
			models.FlagCounter = assets.Games[key].FlagCounter
		case "7":
			assets.Games[key].FlagCounter = 15
			models.FlagCounter = assets.Games[key].FlagCounter
		case "8":
			assets.Games[key].FlagCounter = 18
			models.FlagCounter = assets.Games[key].FlagCounter
		}
	case "hard":
		switch size {
		case "5":
			assets.Games[key].FlagCounter = 10
			models.FlagCounter = assets.Games[key].FlagCounter
		case "6":
			assets.Games[key].FlagCounter = 12
			models.FlagCounter = assets.Games[key].FlagCounter
		case "7":
			assets.Games[key].FlagCounter = 16
			models.FlagCounter = assets.Games[key].FlagCounter
		case "8":
			assets.Games[key].FlagCounter = 20
			models.FlagCounter = assets.Games[key].FlagCounter
		}
	}

	msg := tgbotapi.NewCallback(callback.ID, size+"X"+size)
	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		log.Println(err)
	}

	NewSapperGame(callback, bot, key)
	assets.SavingGame()
}

func ActionsWithBombUpdate(i, j int, key string, callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, db *sql.DB) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, language.LangText(language.UserLang.Language, "you_lost"))

	assets.Games[key].OpenedButtonsField[i][j] = "true"

	ReplyMarkup := CreateFieldMarkUp(assets.Games[key], key)
	msgAboutBomb := tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, ReplyMarkup)

	if _, err := bot.Send(msgAboutBomb); err != nil {
		log.Println(err)
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}

	assets.Games[key] = &models.Game{
		Size: assets.Games[key].Size,
	}

	delete(assets.Games, key)
	database.RatioChange.Losses += 1
	database.RatioChange.Ratio = float64(database.RatioChange.Wins) / float64(database.RatioChange.Losses)
	database.UpdateTable(db, database.RatioChange.Wins, database.RatioChange.Losses, database.RatioChange.Ratio, callback.Message.Chat.ID)
}

func ActionsWithWin(key string, callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, db *sql.DB) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, language.LangText(language.UserLang.Language, "you_win"))
	ReplyMarkup := CreateFieldMarkUp(assets.Games[key], key)
	msgAboutBomb := tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, ReplyMarkup)

	if _, err := bot.Send(msgAboutBomb); err != nil {
		log.Println(err)
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}

	assets.Games[key] = &models.Game{
		Size: assets.Games[key].Size,
	}

	delete(assets.Games, key)
	database.RatioChange.Wins += 1
	database.RatioChange.Ratio = float64(database.RatioChange.Wins) / float64(database.RatioChange.Losses)
	database.UpdateTable(db, database.RatioChange.Wins, database.RatioChange.Losses, database.RatioChange.Ratio, callback.Message.Chat.ID)
}

func CheckDeveloperMode(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, developerMode bool) {
	assets.DeveloperMode = developerMode
	if developerMode {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, language.LangText(language.UserLang.Language, "administrator_mode_on"))
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	} else {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, language.LangText(language.UserLang.Language, "administrator_mode_off"))
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func DataSplit(callbackData string) (key string, i, j int) {
	data := strings.Split(callbackData, "/")
	key = data[0]
	i, _ = strconv.Atoi(data[1])
	j, _ = strconv.Atoi(data[2])
	return key, i, j
}

func CallEditMessage(key string, CallbackChatID int64, CallbackMsgID int, bot *tgbotapi.BotAPI) {

	ReplyMarkup := CreateFieldMarkUp(assets.Games[key], key)
	msgMark := tgbotapi.NewEditMessageReplyMarkup(CallbackChatID, CallbackMsgID, ReplyMarkup)

	if _, err := bot.Send(msgMark); err != nil {
		log.Println(err)
	}
}

func SendAdminBotStart(callbackID string, bot *tgbotapi.BotAPI, UserID int) {
	assets.DeveloperMode = false
	if UserID == assets.AdminId {
		msg := tgbotapi.NewCallback(callbackID, language.LangText(language.UserLang.Language, "bot_started_admin_mode_off"))
		if _, err := bot.AnswerCallbackQuery(msg); err != nil {
			log.Println(err)
		}
	}
}

func SendGameOver(callBack *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewCallback(callBack.ID, language.LangText(language.UserLang.Language, "game_over"))
	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		log.Println(err)
	}
}

func SendMsgLanguage(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, lang string) {
	language.LangChange = true
	if lang == "ru" {
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, language.LangText(lang, "you_take_russian"))
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
		return
	}
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, language.LangText(lang, "you_take_english"))
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func ReturnLevelGameSendMsg(difficult string, callbackID string, bot *tgbotapi.BotAPI) {
	var Text string
	models.Difficult = difficult
	switch difficult {
	case "easy":
		models.Level = map[int]int{5: 4, 6: 4, 7: 6, 8: 9}
		Text = language.LangText(language.UserLang.Language, "easy")
	case "medium":
		models.Level = map[int]int{5: 7, 6: 9, 7: 12, 8: 15}
		Text = language.LangText(language.UserLang.Language, "medium")
	case "hard":
		models.Level = map[int]int{5: 10, 6: 12, 7: 16, 8: 20}
		Text = language.LangText(language.UserLang.Language, "hard")
	}

	msg := tgbotapi.NewCallback(callbackID, Text)
	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		log.Println(err)
	}
	return
}
