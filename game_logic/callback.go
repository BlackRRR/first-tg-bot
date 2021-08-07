package game_logic

import (
	"github.com/BlackRRR/first-tg-bot/assets"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

var DeveloperMode bool //true = admin, false = all users

func ActionWithCallback(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.CallbackQuery.Data == "start" {
		return
	}

	switch update.CallbackQuery.Data {
	case "turn on":
		DeveloperMode = true
		CheckDeveloperMode(update, bot, DeveloperMode)
		return
	case "turn off":
		DeveloperMode = false
		CheckDeveloperMode(update, bot, DeveloperMode)
		return
	}

	if update.CallbackQuery.Data == "5" || update.CallbackQuery.Data == "6" || update.CallbackQuery.Data == "7" || update.CallbackQuery.Data == "8" {
		TakeCallBackFieldSize(update, bot)
		return
	}

	data := strings.Split(update.CallbackQuery.Data, "/")
	key := data[0]
	i, _ := strconv.Atoi(data[1])
	j, _ := strconv.Atoi(data[2])

	if _, exist := assets.Games[key]; !exist {
		return
	}

	if assets.Games[key].PlayingField[i][j] == " " {
		return
	}

	if assets.Games[key].PlayingField[i][j] != "0" && Counter(key) == 0 {
		assets.Games[key] = &assets.Game{}
		assets.Games[key].FillEmptyField()
		assets.Games[key].FillField()
		ReEditField(update, bot, key)
		ActionWithCallback(update, bot)
		return
	}

	if assets.Games[key].OpenedButtonsField[i][j] {
		return
	}

	assets.Games[key].OpenedButtonsField[i][j] = true

	if assets.Games[data[0]].PlayingField[i][j] == "0" {
		OpenZero(i, j, data[0])
	}

	if assets.Games[data[0]].PlayingField[i][j] == "bomb" {
		OpenAllBombsAfterWin(data[0])
		ActionsWithBombUpdate(i, j, data[0], update, bot)
		return
	}

	counter := Counter(data[0])
	if counter == assets.Size*assets.Size-assets.BombCounter {
		OpenAllBombsAfterWin(data[0])
		ActionsWithWin(data[0], update, bot)
		return
	}

	ReplyMarkup := CreateFieldMarkUp(assets.Games[key].PlayingField, assets.Games[key].OpenedButtonsField, key)
	msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, ReplyMarkup)

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
	return
}

func TakeCallBackFieldSize(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	switch update.CallbackQuery.Data {
	case "5":
		assets.Size = 5
		assets.BombCounter = 5
		key := GenerateField()
		NewSapperGame(update, bot, key)
		SavingGame()
		return
	case "6":
		assets.Size = 6
		assets.BombCounter = 6
		key := GenerateField()
		NewSapperGame(update, bot, key)
		SavingGame()
		return
	case "7":
		assets.Size = 7
		assets.BombCounter = 8
		key := GenerateField()
		NewSapperGame(update, bot, key)
		SavingGame()
		return
	case "8":
		assets.Size = 8
		assets.BombCounter = 12
		key := GenerateField()
		NewSapperGame(update, bot, key)
		SavingGame()
		return
	default:
		return
	}
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

	assets.Games[key] = &assets.Game{}
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

	assets.Games[key] = &assets.Game{}
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
