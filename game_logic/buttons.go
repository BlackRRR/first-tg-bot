package game_logic

import (
	"github.com/BlackRRR/first-tg-bot/assets"
	"github.com/BlackRRR/first-tg-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func CreateFieldMarkUp(game *models.Game, key string) tgbotapi.InlineKeyboardMarkup {
	var returnMarkUp tgbotapi.InlineKeyboardMarkup
	for i := 0; i < assets.Games[key].Size; i++ {
		returnMarkUp.InlineKeyboard = append(returnMarkUp.InlineKeyboard, createRowButton(game, i, key))
	}
	return returnMarkUp
}

func createRowButton(game *models.Game, j int, key string) []tgbotapi.InlineKeyboardButton {
	var row []tgbotapi.InlineKeyboardButton
	for i := 0; i < assets.Games[key].Size; i++ {
		row = append(row, createButton(game, i, j, key))
	}
	return row
}

func createButton(game *models.Game, i, j int, key string) tgbotapi.InlineKeyboardButton {
	var data string
	counter := Counter(key)
	data = key + "/" + strconv.Itoa(i) + "/" + strconv.Itoa(j)
	if game.OpenedButtonsField[i][j] {
		if game.PlayingField[i][j] == "bomb" {
			if counter == assets.Games[key].Size*assets.Games[key].Size {
				return tgbotapi.NewInlineKeyboardButtonData("😍", data)
			}
			return tgbotapi.NewInlineKeyboardButtonData("💣", data)
		}
		return tgbotapi.NewInlineKeyboardButtonData(game.PlayingField[i][j], data)
	}
	return tgbotapi.NewInlineKeyboardButtonData("⏺", data)
}
