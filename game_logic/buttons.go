package game_logic

import (
	"github.com/BlackRRR/first-tg-bot/assets"
	"github.com/BlackRRR/first-tg-bot/language"
	"github.com/BlackRRR/first-tg-bot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)

func CreateFieldMarkUp(game *models.Game, key string) tgbotapi.InlineKeyboardMarkup {
	var returnMarkUp tgbotapi.InlineKeyboardMarkup
	for i := 0; i < assets.Games[key].Size; i++ {
		returnMarkUp.InlineKeyboard = append(returnMarkUp.InlineKeyboard, createRowButton(game, i, key))
	}
	returnMarkUp.InlineKeyboard = append(returnMarkUp.InlineKeyboard, CreateRowFlag(game, assets.Games[key].Size+1, 0, key))
	return returnMarkUp
}

func createRowButton(game *models.Game, j int, key string) []tgbotapi.InlineKeyboardButton {
	var row []tgbotapi.InlineKeyboardButton
	for i := 0; i < assets.Games[key].Size; i++ {
		row = append(row, createButton(game, i, j, key))
	}
	return row
}

func CreateRowFlag(game *models.Game, i, j int, key string) []tgbotapi.InlineKeyboardButton {
	var row []tgbotapi.InlineKeyboardButton
	row = append(row, CreateButtonFlag(game, i, j, key))
	return row
}

func createButton(game *models.Game, i, j int, key string) tgbotapi.InlineKeyboardButton {
	var data string
	counter := Counter(key)
	data = key + "/" + strconv.Itoa(i) + "/" + strconv.Itoa(j)

	if game.OpenedButtonsField[i][j] == "❌" {
		return tgbotapi.NewInlineKeyboardButtonData("❌", data)
	}

	if game.OpenedButtonsField[i][j] == "flag" {
		return tgbotapi.NewInlineKeyboardButtonData("🚩", data)
	}

	if game.OpenedButtonsField[i][j] == "false" {
		return tgbotapi.NewInlineKeyboardButtonData("⏺", data)
	}

	if !(game.PlayingField[i][j] == "bomb") {
		return tgbotapi.NewInlineKeyboardButtonData(game.PlayingField[i][j], data)
	}

	if counter == assets.Games[key].Size*assets.Games[key].Size {
		return tgbotapi.NewInlineKeyboardButtonData("😍", data)
	}
	return tgbotapi.NewInlineKeyboardButtonData("💣", data)
}

func CreateButtonFlag(game *models.Game, i, j int, key string) tgbotapi.InlineKeyboardButton {
	var data string
	data = key + "/" + strconv.Itoa(i) + "/" + strconv.Itoa(j)

	if game.Flag == "true" {
		return tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "flag_true"), data)
	} else if game.Flag == "false" {
		return tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "flag_false"), data)
	} else {
		return tgbotapi.NewInlineKeyboardButtonData(language.LangText(language.UserLang.Language, "open_flag"), data)
	}
}
