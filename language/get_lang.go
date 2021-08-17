package language

import (
	"encoding/json"
	"fmt"
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

func FormatText(lang, key string, values ...interface{}) string {
	formatText := LangText(lang, key)
	return fmt.Sprintf(formatText, values...)
}
