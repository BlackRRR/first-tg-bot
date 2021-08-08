package game_logic // rename file name

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func OpenZero(i, j int, key string) {
	for k := -1; k < 2; k++ {
		if i+k < 0 || i+k > Games[key].Size-1 {
			continue
		}
		for l := -1; l < 2; l++ {
			if j+l < 0 || j+l > Games[key].Size-1 {
				continue
			}

			row := i + k
			col := j + l

			if Games[key].OpenedButtonsField[row][col] {
				continue
			}

			Games[key].OpenedButtonsField[row][col] = true
			if Games[key].PlayingField[row][col] == "0" {
				OpenZero(row, col, key)
			}
		}
	}
}

func OpenAllBombsAfterWin(key string) {
	for i := 0; i < Games[key].Size; i++ {
		for j := 0; j < Games[key].Size; j++ {
			if Games[key].OpenedButtonsField[i][j] == false && Games[key].PlayingField[i][j] == "bomb" {
				Games[key].OpenedButtonsField[i][j] = true
			}
		}
	}
}
func Counter(key string) int {
	var counter int
	for i := 0; i < Games[key].Size; i++ {
		for j := 0; j < Games[key].Size; j++ {
			if Games[key].OpenedButtonsField[i][j] {
				counter++
			}
		}
	}
	return counter
}

func ReEditField(update *tgbotapi.Update, bot *tgbotapi.BotAPI, key string) {
	Size := Games[key].Size
	Games[key] = &Game{}
	Games[key].Size = Size
	Games[key].FillEmptyField()
	Games[key].FillField()
	ReplyMarkup := CreateFieldMarkUp(Games[key].PlayingField, Games[key].OpenedButtonsField, key)
	msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, ReplyMarkup)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func printField(field [][]string, key string) {
	for i := 0; i < Games[key].Size; i++ {
		for j := 0; j < Games[key].Size; j++ {
			fmt.Print(field[i][j], " ")
		}
		fmt.Println()
	}
}

func printOpenField(field [][]bool, key string) {
	for i := 0; i < Games[key].Size; i++ {
		for j := 0; j < Games[key].Size; j++ {
			fmt.Print(field[i][j], " ")
		}
		fmt.Println()
	}
}
func SendMsgAll(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
}
