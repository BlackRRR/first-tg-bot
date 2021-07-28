package game_logic

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

const (
	DefaultFieldSize = 8
)

func CreateFieldMarkUp(field [DefaultFieldSize][DefaultFieldSize]string, openField [DefaultFieldSize][DefaultFieldSize]bool, key string) tgbotapi.InlineKeyboardMarkup {
	var returnMarkUp tgbotapi.InlineKeyboardMarkup
	for i := 0; i < DefaultFieldSize; i++ {
		returnMarkUp.InlineKeyboard = append(returnMarkUp.InlineKeyboard, createRowButton(field, i, openField, key))
	}
	return returnMarkUp
}

func createRowButton(field [DefaultFieldSize][DefaultFieldSize]string, j int, openField [DefaultFieldSize][DefaultFieldSize]bool, key string) []tgbotapi.InlineKeyboardButton {
	var row []tgbotapi.InlineKeyboardButton
	for i := 0; i < DefaultFieldSize; i++ {
		row = append(row, createButton(field, openField, i, j, key))
	}
	return row
}

func createButton(field [DefaultFieldSize][DefaultFieldSize]string, openField [DefaultFieldSize][DefaultFieldSize]bool, i, j int, key string) tgbotapi.InlineKeyboardButton {
	var data string
	data = key + "/" + strconv.Itoa(i) + "/" + strconv.Itoa(j) + "/" + field[i][j]
	if openField[i][j] {
		if field[i][j] == "bomb" {
			return tgbotapi.NewInlineKeyboardButtonData("😍", data)
		}
		return tgbotapi.NewInlineKeyboardButtonData(field[i][j], data)
	}
	return tgbotapi.NewInlineKeyboardButtonData("⏺", data)
}
