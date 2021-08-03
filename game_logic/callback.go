package game_logic

import (
	"fmt"
	"github.com/BlackRRR/first-tg-bot/assets"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

func ActionWithCallback(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := strings.Split(update.CallbackQuery.Data, "/")
	key := data[0]
	i, _ := strconv.Atoi(data[1])
	j, _ := strconv.Atoi(data[2])

	if assets.Games[key].PlayingField[i][j] != "0" && Counter(key) == 0 && assets.Games[key].PlayingField[i][j] != "" {
		assets.Games[key] = &assets.Game{}
		assets.Games[key].FillField()
		ReEditField(update, bot, key)
		ActionWithCallback(update, bot)
		return
	}

	if _, exist := assets.Games[key]; !exist {
		return
	}

	if assets.Games[key].OpenedButtonsField[i][j] {
		return
	}

	if assets.Games[key].PlayingField[i][j] == "" {
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
	fmt.Println(counter)
	if counter == DefaultFieldSize*DefaultFieldSize-assets.DefaultBombCounter {
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
	return

}

func OpenZero(i, j int, key string) {
	for k := -1; k < 2; k++ {
		if i+k < 0 || i+k > DefaultFieldSize-1 {
			continue
		}
		for l := -1; l < 2; l++ {
			if j+l < 0 || j+l > DefaultFieldSize-1 {
				continue
			}

			row := i + k
			col := j + l

			if assets.Games[key].OpenedButtonsField[row][col] {
				continue
			}

			assets.Games[key].OpenedButtonsField[row][col] = true
			if assets.Games[key].PlayingField[row][col] == "0" {
				OpenZero(row, col, key)
			}
		}
	}
}

func printField(field [DefaultFieldSize][DefaultFieldSize]string) {
	for i := 0; i < DefaultFieldSize; i++ {
		for j := 0; j < DefaultFieldSize; j++ {
			fmt.Print(field[i][j], " ")
		}
		fmt.Println()
	}
}

func printOpenField(field [DefaultFieldSize][DefaultFieldSize]bool) {
	for i := 0; i < DefaultFieldSize; i++ {
		for j := 0; j < DefaultFieldSize; j++ {
			fmt.Print(field[i][j], " ")
		}
		fmt.Println()
	}
}

func OpenAllBombsAfterWin(key string) {
	for i := 0; i < assets.DefaultFieldSize; i++ {
		for j := 0; j < assets.DefaultFieldSize; j++ {
			if assets.Games[key].OpenedButtonsField[i][j] == false && assets.Games[key].PlayingField[i][j] == "bomb" {
				assets.Games[key].OpenedButtonsField[i][j] = true
			}
		}
	}
}
func Counter(key string) int {
	var counter int
	for i := 0; i < DefaultFieldSize; i++ {
		for j := 0; j < DefaultFieldSize; j++ {
			if assets.Games[key].OpenedButtonsField[i][j] {
				counter++
			}
		}
	}
	return counter
}

func ReEditField(update *tgbotapi.Update, bot *tgbotapi.BotAPI, key string) {
	ReplyMarkup := CreateFieldMarkUp(assets.Games[key].PlayingField, assets.Games[key].OpenedButtonsField, key)
	msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, ReplyMarkup)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
