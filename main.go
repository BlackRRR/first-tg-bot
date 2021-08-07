package main

import (
	"fmt"
	"github.com/BlackRRR/first-tg-bot/game_logic"
	"log"
	"math/rand"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	AdminId       = 872383555
	AdminUserName = "BlackR0_0"
)

func main() {
	rand.Seed(time.Now().Unix())

	bot, updates := startBot()
	game_logic.UploadGame()

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

func SendWorkIsUnderwayMessage(update *tgbotapi.Update, bot *tgbotapi.BotAPI) { //TODO: rename to ...UnderWayMsg; TODO: you don't need to pass the whole update here, just pass the chatId here
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ведуться работы...")
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
	return
}

func SendWorkIsUnderwayCallback(update *tgbotapi.Update, bot *tgbotapi.BotAPI) { //TODO: transmit only the chatId
	msg := tgbotapi.NewCallback(update.CallbackQuery.ID, "Ведуться работы...")
	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		log.Println(err)
	}
	return
}

func checkUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message != nil {

		//db := database.DBConn()
		//_, username := database.GetAllData(db)
		//CheckUsersFromDB(username)
		//database.AddDB(db,update.Message.From.ID,update.Message.From.UserName)

		if game_logic.DeveloperMode && update.Message.From.ID != AdminId {
			SendWorkIsUnderwayMessage(update, bot)
			return
		}

		switch update.Message.Command() {
		case "sapper":
			game_logic.TakeFieldSize(update, bot)
			return
		case "admin":
			if update.Message.From.ID == AdminId {
				CreateAdminButtons(update, bot)
				return
			}
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите /sapper чтобы начать игру")
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
		return
	}

	if update.CallbackQuery != nil {
		if game_logic.DeveloperMode && update.CallbackQuery.From.ID != AdminId {
			SendWorkIsUnderwayCallback(update, bot)
			return
		}

		game_logic.ActionWithCallback(update, bot)
		game_logic.SavingGame()
		return
	}
}

func CreateAdminButtons(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "/admin")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отправить всем сообщение о старте бота", "start")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Включить Режим Администратора", "turn on"),
			tgbotapi.NewInlineKeyboardButtonData("Вылючить Режим Администратора", "turn off"),
		))
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func CheckUsersFromDB(userName []string) {
	for i := range userName {
		fmt.Println(userName[i])
	}
}
