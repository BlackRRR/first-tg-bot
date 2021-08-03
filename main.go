package main

import (
	"github.com/BlackRRR/first-tg-bot/assets"
	"github.com/BlackRRR/first-tg-bot/game_logic"
	"log"
	"math/rand"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	AvailableSymbolInKey = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	GameKeyLength        = 16
)

func main() {
	rand.Seed(time.Now().Unix())

	bot, updates := startBot()
	assets.UploadGame()

	actionsWithUpdates(updates, bot)
}

func startBot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	botToken := takeBotToken()

	bot, err := tgbotapi.NewBotAPI(botToken)

	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic("Failed to initialize bot: " + err.Error())
	}

	log.Println("The bot is running")

	return bot, updates
}

func takeBotToken() string {
	content, _ := os.ReadFile("botToken.txt")
	return string(content)
}

func actionsWithUpdates(updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) {
	for update := range updates {
		checkUpdate(&update, bot)
	}
}

func checkUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil {
		if update.Message.Command() == "sapper" {
			key := GenerateField()
			NewSapperGame(update.Message, bot, key)
			assets.SavingGame()
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите /sapper чтобы начать игру")
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
		return
	}

	if update.CallbackQuery != nil {
		game_logic.ActionWithCallback(update, bot)
		assets.SavingGame()
		return
	}
}

func NewSapperGame(message *tgbotapi.Message, bot *tgbotapi.BotAPI, key string) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Игра началась")
	msg.ReplyMarkup = game_logic.CreateFieldMarkUp(assets.Games[key].PlayingField, assets.Games[key].OpenedButtonsField, key)
	msgData, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}

	assets.Games[key].MessageID = msgData.MessageID
}

func GenerateField() string {
	key := generateKey()
	assets.Games[key] = &assets.Game{}
	assets.Games[key].FillField()
	return key
}

func generateKey() string {
	var key string
	slice := []rune(AvailableSymbolInKey)
	for i := 0; i < GameKeyLength; i++ {
		key += string(slice[rand.Intn(len(slice))])
	}
	return key
}
