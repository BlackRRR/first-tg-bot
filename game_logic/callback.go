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
	keyFromOpenButton := data[0]
	i, _ := strconv.Atoi(data[1])
	j, _ := strconv.Atoi(data[2])

	if _, exist := assets.Games[keyFromOpenButton]; !exist {
		return
	}

	if assets.Games[keyFromOpenButton].OpenedButtonsField[i][j] {
		return
	}

	assets.Games[keyFromOpenButton].OpenedButtonsField[i][j] = true
	if assets.Games[data[0]].PlayingField[i][j] == "0" {
		OpenZero(i, j, data[0])
	}

	if data[3] == "bomb" {
		OpenAllBombsAfterWin(data[0], update, bot)
		ActionsWithBombUpdate(update, bot)
		return
	}

	counter := Counter(data[0])
	if counter == DefaultFieldSize*DefaultFieldSize-assets.DefaultBombCounter {
		OpenAllBombsAfterWin(data[0], update, bot)
		ActionsWithWin(update, bot)
		return
	}

	ReplyMarkup := CreateFieldMarkUp(assets.Games[keyFromOpenButton].PlayingField, assets.Games[keyFromOpenButton].OpenedButtonsField, keyFromOpenButton)
	msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, assets.Games[keyFromOpenButton].MessageID, ReplyMarkup)

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
	return
}

func ActionsWithBombUpdate(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	data := strings.Split(update.CallbackQuery.Data, "/")
	keyFromOpenButton := data[0]
	i, _ := strconv.Atoi(data[1])
	j, _ := strconv.Atoi(data[2])
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы проиграли\nнапишите /sapper для новой игры")

	assets.Games[keyFromOpenButton].OpenedButtonsField[i][j] = true

	ReplyMarkup := CreateFieldMarkUp(assets.Games[keyFromOpenButton].PlayingField, assets.Games[keyFromOpenButton].OpenedButtonsField, keyFromOpenButton)
	msg1 := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, assets.Games[keyFromOpenButton].MessageID, ReplyMarkup)

	if _, err := bot.Send(msg1); err != nil {
		log.Println(err)
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}

	assets.Games[keyFromOpenButton] = &assets.Game{}
	return
}

func ActionsWithWin(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Вы выйграли нажмите /sapper чтобы начать новую игру")

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
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

func OpenAllBombsAfterWin(key string, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	for i := 0; i < assets.DefaultFieldSize; i++ {
		for j := 0; j < assets.DefaultFieldSize; j++ {
			if assets.Games[key].OpenedButtonsField[i][j] == false && assets.Games[key].PlayingField[i][j] == "bomb" {
				assets.Games[key].OpenedButtonsField[i][j] = true
				ReplyMarkup := CreateFieldMarkUp(assets.Games[key].PlayingField, assets.Games[key].OpenedButtonsField, key)
				msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, assets.Games[key].MessageID, ReplyMarkup)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
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
