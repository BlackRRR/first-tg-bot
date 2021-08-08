package game_logic

import (
	"github.com/BlackRRR/first-tg-bot/assets"
	"github.com/BlackRRR/first-tg-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

// developer mode is not part of the game logic, it should rather lie in assets

//TODO: Refactor ActionWithCallBack func
/* example
func checkCallbackQuery(botLang string, callbackQuery *tgbotapi.CallbackQuery) {
	switch strings.Split(callbackQuery.Data, "/")[0] {
	case "moreMoney":
		GetBonus(botLang, callbackQuery)
	case "withdrawal_exit":
		CheckSubsAndWithdrawal(botLang, callbackQuery, callbackQuery.From.ID)
	case "change_lang":
		ChangeLanguage(botLang, callbackQuery)
	case "admin":
		admin.AnalyseAdminCallback(botLang, callbackQuery)
	}
}

remove the pieces of logic into separate functions and call them depending on the data
*/

func ActionWithCallback(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.CallbackQuery.Data == "start" {
		return
	}

	switch update.CallbackQuery.Data {
	case "turn on":
		assets.DeveloperMode = true
		CheckDeveloperMode(update, bot, assets.DeveloperMode)
		return
	case "turn off":
		assets.DeveloperMode = false
		CheckDeveloperMode(update, bot, assets.DeveloperMode)
		return
	}

	if update.CallbackQuery.Data == "5" || update.CallbackQuery.Data == "6" || update.CallbackQuery.Data == "7" || update.CallbackQuery.Data == "8" {
		TakeCallBackFieldSize(update, bot)
		return
	}

	key, i, j := DataSplit(update.CallbackQuery.Data)

	if _, exist := assets.Games[key]; !exist {
		return
	}

	if assets.Games[key].PlayingField[i][j] == " " {
		return
	}

	if assets.Games[key].PlayingField[i][j] != "0" && Counter(key) == 0 {
		ReEditField(update, bot, key)
		ActionWithCallback(update, bot)
		return
	}

	if assets.Games[key].OpenedButtonsField[i][j] {
		return
	}

	assets.Games[key].OpenedButtonsField[i][j] = true

	if assets.Games[key].PlayingField[i][j] == "0" {
		OpenZero(i, j, key)
	}

	if assets.Games[key].PlayingField[i][j] == "bomb" {
		OpenAllBombsAfterWin(key)
		ActionsWithBombUpdate(i, j, key, update, bot)
		return
	}

	if Counter(key) == assets.Games[key].Size*assets.Games[key].Size-assets.BombCounter {
		OpenAllBombsAfterWin(key)
		ActionsWithWin(key, update, bot)
		return
	}

	CallEditMessage(key, update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, bot)

	return
}

func TakeCallBackFieldSize(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	var key string
	switch update.CallbackQuery.Data {
	case "5":
		assets.BombCounter = 5 //TODO: delete global variable BombCounter
		key = GenerateField(5)
	case "6":
		assets.BombCounter = 6
		key = GenerateField(6)
	case "7":
		assets.BombCounter = 8
		key = GenerateField(7)
	case "8":
		assets.BombCounter = 12
		key = GenerateField(8)
	default:
		return
	}

	NewSapperGame(update, bot, key)
	assets.SavingGame()
}

func ActionsWithBombUpdate(i, j int, key string, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы проиграли\nнапишите /sapper для новой игры")

	assets.Games[key].OpenedButtonsField[i][j] = true

	ReplyMarkup := CreateFieldMarkUp(assets.Games[key].PlayingField, assets.Games[key].OpenedButtonsField, key)
	msgAboutBomb := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, ReplyMarkup)

	if _, err := bot.Send(msgAboutBomb); err != nil {
		log.Println(err)
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}

	assets.Games[key] = &models.Game{
		Size: assets.Games[key].Size,
	}
	assets.Games[key].FillEmptyField()
	return
}

func ActionsWithWin(key string, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы выйграли нажмите /sapper чтобы начать новую игру")
	ReplyMarkup := CreateFieldMarkUp(assets.Games[key].PlayingField, assets.Games[key].OpenedButtonsField, key)
	msgAboutBomb := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, ReplyMarkup)

	if _, err := bot.Send(msgAboutBomb); err != nil {
		log.Println(err)
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}

	assets.Games[key] = &models.Game{
		Size: assets.Games[key].Size,
	}
	assets.Games[key].FillEmptyField()
	return

}

func CheckDeveloperMode(update *tgbotapi.Update, bot *tgbotapi.BotAPI, developerMode bool) {
	if developerMode {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Режим Администрации включен")
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	} else {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Режим Администрации выключен")
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
	ReplyMarkup := CreateFieldMarkUp(assets.Games[key].PlayingField, assets.Games[key].OpenedButtonsField, key)
	msg := tgbotapi.NewEditMessageReplyMarkup(CallbackChatID, CallbackMsgID, ReplyMarkup)

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
