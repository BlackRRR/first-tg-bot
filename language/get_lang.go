package language

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

const (
	AvailableLangPath = "language/"
)

var (
	AvailableLang = []string{"ru", "en"}

	Language = make([]map[string]string, 2)

	UserLang *UserLanguage

	LangChange = false
)

type UserLanguage struct {
	Language string
}

func LangParsing() {
	for i, lang := range AvailableLang {
		bytes, _ := os.ReadFile(AvailableLangPath + lang + ".json")
		_ = json.Unmarshal(bytes, &Language[i])
	}
}

func LangText(lang, key string) string {
	index := LangIndex(lang)
	return Language[index][key]
}

func LangIndex(lang string) int {
	for i, elem := range AvailableLang {
		if elem == lang {
			return i
		}
	}
	return 0
}

func ReturnLanguage(lang string) {
	UserLang.Language = lang
}

func GenerateLanguage() {
	Lang := &UserLanguage{
		Language: "ru",
	}
	UserLang = Lang
}

func TakeLanguage(userID int64, bot *tgbotapi.BotAPI) {
	msgLanguage := tgbotapi.NewMessage(userID, LangText(UserLang.Language, "select_language"))
	msgLanguage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(LangText(UserLang.Language, "russian"), "ru"),
			tgbotapi.NewInlineKeyboardButtonData(LangText(UserLang.Language, "english"), "en")))
	_, err := bot.Send(msgLanguage)
	if err != nil {
		log.Println(err)
	}
}
