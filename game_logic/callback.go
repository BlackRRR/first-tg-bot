package game_logic

import (
	"github.com/BlackRRR/first-tg-bot/assets"
	"github.com/BlackRRR/first-tg-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

func ActionWithCallback(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	switch callback.Data {
	case "start":
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
	default:
		HandlingGameLogic(callback, bot)
	}
}

func HandlingGameLogic(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	key, i, j := DataSplit(callback.Data)

	if _, exist := assets.Games[key]; !exist {
		return
	}

	if assets.Games[key].PlayingField[i][j] == " " {
		return
	}

	if assets.Games[key].PlayingField[i][j] != "0" && Counter(key) == 0 {
		ReEditField(callback, bot, key)
		ActionWithCallback(callback, bot)
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
		ActionsWithBombUpdate(i, j, key, callback, bot)
		return
	}

	if Counter(key) == assets.Games[key].Size*assets.Games[key].Size-assets.Games[key].BombCounter {
		OpenAllBombsAfterWin(key)
		ActionsWithWin(key, callback, bot)
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

func ActionsWithBombUpdate(i, j int, key string, callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Вы проиграли\nнапишите /sapper для новой игры")

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
	assets.Games[key].FillEmptyField()
}

func ActionsWithWin(key string, callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Вы выйграли нажмите /sapper чтобы начать новую игру")
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
	assets.Games[key].FillEmptyField()
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
	ReplyMarkup := CreateFieldMarkUp(assets.Games[key], key)
	msg := tgbotapi.NewEditMessageReplyMarkup(CallbackChatID, CallbackMsgID, ReplyMarkup)

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
