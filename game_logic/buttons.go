package game_logic

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func CreateFieldMarkUp(field [][]string, openField [][]bool, key string) tgbotapi.InlineKeyboardMarkup {
	var returnMarkUp tgbotapi.InlineKeyboardMarkup
	for i := 0; i < Games[key].Size; i++ {
		returnMarkUp.InlineKeyboard = append(returnMarkUp.InlineKeyboard, createRowButton(field, i, openField, key))
	}
	return returnMarkUp
}

func createRowButton(field [][]string, j int, openField [][]bool, key string) []tgbotapi.InlineKeyboardButton {
	var row []tgbotapi.InlineKeyboardButton
	for i := 0; i < Games[key].Size; i++ {
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
			if counter == Games[key].Size*Games[key].Size {
				return tgbotapi.NewInlineKeyboardButtonData("😍", data)
			}
			return tgbotapi.NewInlineKeyboardButtonData("💣", data)
		}
		return tgbotapi.NewInlineKeyboardButtonData(field[i][j], data)
	}
	return tgbotapi.NewInlineKeyboardButtonData("⏺", data)
}
