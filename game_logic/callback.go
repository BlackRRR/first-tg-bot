package game_logic

import (
	"database/sql"
	"github.com/BlackRRR/first-tg-bot/assets"
	"github.com/BlackRRR/first-tg-bot/database"
	"github.com/BlackRRR/first-tg-bot/language"
	"github.com/BlackRRR/first-tg-bot/models"
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
		return
	case "turn on":
		CheckDeveloperMode(callback, bot, true)
		return
	case "turn off":
		CheckDeveloperMode(callback, bot, false)
		return
	case "5", "6", "7", "8":
		TakeCallBackFieldSize(callback, bot, callback.Data)
		return
	case "ru":
		language.ReturnLanguage("ru")
		SendMsgLanguage(callback, bot, "ru")
	case "en":
		language.ReturnLanguage("en")
		SendMsgLanguage(callback, bot, "en")
	default:
		HandlingGameLogic(callback, bot, users, db)
	}
}

func HandlingGameLogic(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, users []database.User, db *sql.DB) {
	key, i, j := DataSplit(callback.Data)

	if assets.Games[key] == nil {
		SendGameOver(callback, bot)
	}

	if _, exist := assets.Games[key]; !exist {
		return
	}

	if assets.Games[key].PlayingField[i][j] == " " {
		return
	}

	if assets.Games[key].PlayingField[i][j] != "0" && Counter(key) == 0 {
		ReEditField(callback, bot, key)
		ActionWithCallback(callback, bot, users, db)
		return
	}

	if assets.Games[key].OpenedButtonsField[i][j] {
		return
	}

	assets.Games[key].OpenedButtonsField[i][j] = true

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

	CallEditMessage(key, callback.Message.Chat.ID, callback.Message.MessageID, bot)

}

func TakeCallBackFieldSize(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, size string) {
	bombCounter := map[int]int{
		5: 5,
		6: 6,
		7: 8,
		8: 12,
	}

	intSize, _ := strconv.Atoi(size)
	key := GenerateField(intSize, bombCounter[intSize])

	NewSapperGame(callback, bot, key)
	assets.SavingGame()
}

func ActionsWithBombUpdate(i, j int, key string, callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, db *sql.DB) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, language.LangText(language.UserLang.Language, "you_lost"))

	assets.Games[key].OpenedButtonsField[i][j] = true

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
	database.RatioChange.Ratio = float64(database.RatioChange.Wins / database.RatioChange.Losses)
	database.UpdateTable(db, database.RatioChange.Wins, database.RatioChange.Losses, database.RatioChange.Ratio, callback.Message.From.ID)
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
	database.RatioChange.Ratio = float64(database.RatioChange.Wins / database.RatioChange.Losses)
	database.UpdateTable(db, database.RatioChange.Wins, database.RatioChange.Losses, database.RatioChange.Ratio, callback.Message.From.ID)
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
	msg := tgbotapi.NewEditMessageReplyMarkup(CallbackChatID, CallbackMsgID, ReplyMarkup)

	if _, err := bot.Send(msg); err != nil {
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
