package main

import (
	"encoding/json"
	"github.com/BlackRRR/first-tg-bot/assets"
	"github.com/BlackRRR/first-tg-bot/database"
	"github.com/BlackRRR/first-tg-bot/game_logic"
	"github.com/BlackRRR/first-tg-bot/language"
	"github.com/BlackRRR/first-tg-bot/msgs"
	"github.com/BlackRRR/first-tg-bot/top"
	"log"
	"math/rand"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	language.LangParsing()
	for update := range updates {
		checkUpdate(&update, bot)
	}
}

func checkUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if !language.LangChange {
		language.GenerateLanguage()
	}

	db := database.DBConn()
	users := database.GetAllData(db)

	if update.Message != nil {

		flag := database.CheckIdenticalValues(update, users)
		if flag {
			database.AddUser(db, update.Message.Chat.ID, update.Message.From.UserName, update.Message.From.LanguageCode)
		}

		if !language.LangChange {
			language.ReturnLanguage(update.Message.From.LanguageCode)
		}

		if database.RatioChange == nil {
			database.AddRatio(users, update.Message.Chat.ID)
		}

		if assets.DeveloperMode && update.Message.From.ID != assets.AdminId {
			msgs.SendWorkIsUnderWayMsg(update.Message.Chat.ID, bot)
			return
		}

		switch update.Message.Command() {
		case "sapper":
			msgs.SelectDifficultGame(update, bot)
			return
		case "admin":
			if update.Message.From.ID == assets.AdminId {
				CreateAdminButtons(update, bot)
				return
			}
		case "language":
			msgs.TakeLanguage(update.Message.Chat.ID, bot)
			return
		case "top":
			top.CreateTop(users, update, bot)
			return
		case "rules":
			msgs.Rules(update, bot)
			return
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, language.LangText(language.UserLang.Language, "start_game"))
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
		return
	}

	if update.CallbackQuery != nil {
		if assets.DeveloperMode && update.CallbackQuery.From.ID != assets.AdminId {
			msgs.SendWorkIsUnderWayCallBack(update.CallbackQuery.ID, bot)
			return
		}

		game_logic.ActionWithCallback(update.CallbackQuery, bot, users, db)
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
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, language.LangText(language.UserLang.Language, "admin_Panel"))
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "send_all_start"), "start")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "turning_administrator_mode_on"), "turn on"),
			tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "turning_administrator_mode_off"), "turn off"),
		))
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
