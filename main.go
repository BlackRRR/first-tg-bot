package main

import (
	"encoding/json"
	"fmt"
	"github.com/BlackRRR/first-tg-bot/assets"
	"github.com/BlackRRR/first-tg-bot/database"
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

func SendWorkIsUnderWayMsg(ChatId int64, bot *tgbotapi.BotAPI) { // rename to ...UnderWayMsg; : you don't need to pass the whole update here, just pass the chatId here
	msg := tgbotapi.NewMessage(ChatId, "Ведуться работы...")
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
	return
}

func SendWorkIsUnderWayCallBack(CallbackId string, bot *tgbotapi.BotAPI) { // transmit only the chatId
	msg := tgbotapi.NewCallback(CallbackId, "Ведуться работы...")
	if _, err := bot.AnswerCallbackQuery(msg); err != nil {
		log.Println(err)
	}
	return
}

func checkUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	printUpdate(update)
	if update.Message != nil {

		db := database.DBConn()
		users := database.GetAllData(db)
		CheckUsersFromDB(users)
		//database.AddDB(db,update.Message.From.ID,update.Message.From.UserName)

		if assets.DeveloperMode && update.Message.From.ID != AdminId {
			SendWorkIsUnderWayMsg(update.Message.Chat.ID, bot)
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
		if assets.DeveloperMode && update.CallbackQuery.From.ID != AdminId {
			SendWorkIsUnderWayCallBack(update.CallbackQuery.ID, bot)
			return
		}

		game_logic.ActionWithCallback(update.CallbackQuery, bot)
		assets.SavingGame()
		return
	}
}

func printUpdate(update *tgbotapi.Update) {
	updateData, err := json.MarshalIndent(update, "", "  ")
	if err != nil {
		log.Println(err)
	}

	log.Println(string(updateData))
}

func CreateAdminButtons(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Панель администратора ⚙️")
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

func CheckUsersFromDB(users []database.User) {
	var user database.User
	for i := range users {
		user = users[i]
		fmt.Println(user.UserID)
	}
}

func SendMsgAll(userID int, bot *tgbotapi.BotAPI) {

}
