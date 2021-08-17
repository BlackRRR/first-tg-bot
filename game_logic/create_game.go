package game_logic

import (
	"github.com/BlackRRR/first-tg-bot/assets"
	"github.com/BlackRRR/first-tg-bot/language"
	"github.com/BlackRRR/first-tg-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
)

const (
	AvailableSymbolInKey = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	GameKeyLength        = 16
)

func TakeFieldSize(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	Text := language.LangText(language.UserLang.Language, "take_field_size")
	ReplyMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("5X5", "5"),
			tgbotapi.NewInlineKeyboardButtonData("6X6", "6"),
			tgbotapi.NewInlineKeyboardButtonData("7X7", "7"),
			tgbotapi.NewInlineKeyboardButtonData("8X8", "8"),
		))

	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      callback.Message.Chat.ID,
			MessageID:   callback.Message.MessageID,
			ReplyMarkup: &ReplyMarkup,
		},
		Text: Text,
	}
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func NewSapperGame(callback *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI, key string) {
	Text := language.FormatText(language.UserLang.Language, "game_started", assets.Games[key].BombCounter, assets.Games[key].FlagCounter)
	ReplyMarkup := CreateFieldMarkUp(assets.Games[key], key)

	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      callback.Message.Chat.ID,
			MessageID:   callback.Message.MessageID,
			ReplyMarkup: &ReplyMarkup,
		},
		Text: Text,
	}
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
