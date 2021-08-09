package game_logic

import (
	"github.com/BlackRRR/first-tg-bot/assets"
	"github.com/BlackRRR/first-tg-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
)

const (
	AvailableSymbolInKey = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	GameKeyLength        = 16
)

func TakeFieldSize(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите размер поля")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("5X5", "5"),
			tgbotapi.NewInlineKeyboardButtonData("6X6", "6"),
			tgbotapi.NewInlineKeyboardButtonData("7X7", "7"),
			tgbotapi.NewInlineKeyboardButtonData("8X8", "8"),
		))
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func NewSapperGame(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, key string) {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Игра началась")
	msg.ReplyMarkup = CreateFieldMarkUp(assets.Games[key], key)
	msgData, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

	assets.Games[key].MessageID = msgData.MessageID
}

func generateKey() string {
	var key string
	slice := []rune(AvailableSymbolInKey)
	for i := 0; i < GameKeyLength; i++ {
		key += string(slice[rand.Intn(len(slice))])
	}
	return key
}

func GenerateField(size int, bombCounter int) string {
	key := generateKey()
	game := &models.Game{
		Size:        size,
		BombCounter: bombCounter,
	}
	game.FillEmptyField()
	game.FillField()
	assets.Games[key] = game
	return key
}
