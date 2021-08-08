package game_logic

import (
	"fmt"
	"github.com/BlackRRR/first-tg-bot/assets"
	"github.com/BlackRRR/first-tg-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func OpenZero(i, j int, key string) {
	for k := -1; k < 2; k++ {
		if i+k < 0 || i+k > assets.Games[key].Size-1 {
			continue
		}
		for l := -1; l < 2; l++ {
			if j+l < 0 || j+l > assets.Games[key].Size-1 {
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

func OpenAllBombsAfterWin(key string) {
	for i := 0; i < assets.Games[key].Size; i++ {
		for j := 0; j < assets.Games[key].Size; j++ {
			if assets.Games[key].OpenedButtonsField[i][j] == false && assets.Games[key].PlayingField[i][j] == "bomb" {
				assets.Games[key].OpenedButtonsField[i][j] = true
			}
		}
	}
}
func Counter(key string) int {
	var counter int
	for i := 0; i < assets.Games[key].Size; i++ {
		for j := 0; j < assets.Games[key].Size; j++ {
			if assets.Games[key].OpenedButtonsField[i][j] {
				counter++
			}
		}
	}
	return counter
}

func ReEditField(update *tgbotapi.Update, bot *tgbotapi.BotAPI, key string) {
	//Size := Games[key].Size
	assets.Games[key] = &models.Game{
		Size: assets.Games[key].Size,
	}

	assets.Games[key].FillEmptyField()
	assets.Games[key].FillField()
	ReplyMarkup := CreateFieldMarkUp(assets.Games[key].PlayingField, assets.Games[key].OpenedButtonsField, key)
	msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, ReplyMarkup)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func printField(field [][]string, key string) {
	for i := 0; i < assets.Games[key].Size; i++ {
		for j := 0; j < assets.Games[key].Size; j++ {
			fmt.Print(field[i][j], " ")
		}
		fmt.Println()
	}
}

func printOpenField(field [][]bool, key string) {
	for i := 0; i < assets.Games[key].Size; i++ {
		for j := 0; j < assets.Games[key].Size; j++ {
			fmt.Print(field[i][j], " ")
		}
		fmt.Println()
	}
}

func SendMsgAll(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
}
