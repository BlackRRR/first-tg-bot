package game_logic

import (
	"github.com/BlackRRR/first-tg-bot/assets"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func CreateFieldMarkUp(field [][]string, openField [][]bool, key string) tgbotapi.InlineKeyboardMarkup {
	var returnMarkUp tgbotapi.InlineKeyboardMarkup
	for i := 0; i < assets.Size; i++ {
		returnMarkUp.InlineKeyboard = append(returnMarkUp.InlineKeyboard, createRowButton(field, i, openField, key))
	}
	return returnMarkUp
}

func createRowButton(field [][]string, j int, openField [][]bool, key string) []tgbotapi.InlineKeyboardButton {
	var row []tgbotapi.InlineKeyboardButton
	for i := 0; i < assets.Size; i++ {
		row = append(row, createButton(field, openField, i, j, key))
	}
	return row
}

func createButton(field [][]string, openField [][]bool, i, j int, key string) tgbotapi.InlineKeyboardButton {
	var data string
	counter := Counter(key)
	data = key + "/" + strconv.Itoa(i) + "/" + strconv.Itoa(j)
	if openField[i][j] {
		if field[i][j] == "bomb" {
			if counter == assets.Size*assets.Size {
				return tgbotapi.NewInlineKeyboardButtonData("😍", data)
			}
			return tgbotapi.NewInlineKeyboardButtonData("💣", data)
		}
		return tgbotapi.NewInlineKeyboardButtonData(field[i][j], data)
	}
	return tgbotapi.NewInlineKeyboardButtonData("⏺", data)
}
