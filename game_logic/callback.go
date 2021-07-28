package game_logic

import (
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
	//for k := -1; k < 2; k++ {
	//	if i+k < 0 || i+k > DefaultFieldSize-1 {
	//		continue
	//	}
	//	for l := -1; l < 2; l++ {
	//		if j+l < 0 || j+l > DefaultFieldSize-1 {
	//			continue
	//		}
	//		if !assets.Games[key].OpenedButtonsField[i+k][j+l] {
	//			assets.Games[key].OpenedButtonsField[i+k][j+l] = true
	//			ReplyMarkup := CreateFieldMarkUp(assets.Games[key].PlayingField, assets.Games[key].OpenedButtonsField, key)
	//			msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, assets.Games[key].MessageID, ReplyMarkup)
	//			_, err := bot.Send(msg)
	//			if err != nil {
	//				log.Println(err)
	//			}
	//			if assets.Games[key].PlayingField[i+k][j+l] == "0" {
	//				OpenZero(i+k, j+l, update, bot, key)
	//			}
	//		}
	//
	//	}
	//	assets.Games[key].OpenedButtonsField[i+k][j+l] = true
	//	ReplyMarkup := CreateFieldMarkUp(assets.Games[key].PlayingField, assets.Games[key].OpenedButtonsField, key)
	//	msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, assets.Games[key].MessageID, ReplyMarkup)
	//	_, err := bot.Send(msg)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	if assets.Games[key].PlayingField[i+k][j+l] == "0" {
	//		OpenZero(i+k, j+l, update, bot, key)
	//	}
	//}

	for k := -1; k < 2; k++ {
		if i+k < 0 || i+k > DefaultFieldSize-1 {
			continue
		}
		for l := -1; l < 2; l++ {
			if j+l < 0 || j+l > DefaultFieldSize-1 {
				continue
			}
			assets.Games[key].OpenedButtonsField[i+k][j+l] = true
			if assets.Games[key].PlayingField[i+k][j+l] == "0" {
				OpenZero(i+k, j+l, key)
			}
		}
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
