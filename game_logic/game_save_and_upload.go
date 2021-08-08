package game_logic //TODO: this code should rather be in the assets folder, since this is not the logic of the game, but the logic of the launch

import (
	"encoding/json"
	"fmt"
	"github.com/BlackRRR/first-tg-bot/assets"
	"log"
	"os"
)

func UploadGame() {
	var game map[string]*Game
	data, err := os.ReadFile(assets.GamesSavePath)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(data, &game)
	if err != nil {
		fmt.Println(err)
	}

	Games = game
}

func SavingGame() {
	dataSave, err := json.MarshalIndent(Games, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}
	err = os.WriteFile(assets.GamesSavePath, dataSave, 0600)
	if err != nil {
		log.Fatalln(err)
	}
}
