package game_logic

import (
	"encoding/json"
	"fmt"
	"github.com/BlackRRR/first-tg-bot/assets"
	"log"
	"os"
)

func UploadGame() {
	var game map[string]*assets.Game
	data, err := os.ReadFile(assets.GamesSavePath)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(data, &game)
	if err != nil {
		fmt.Println(err)
	}

	assets.Games = game
}

func SavingGame() {
	dataSave, err := json.MarshalIndent(assets.Games, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	err = os.WriteFile(assets.GamesSavePath, dataSave, 0600)
	if err != nil {
		log.Fatalln(err)
	}
}
