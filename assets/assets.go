package assets

import (
	"encoding/json"
	"fmt"
	"github.com/BlackRRR/first-tg-bot/models"
	"log"
	"os"
)

const (
	GamesSavePath = "assets/game_save.json"
)

var (
	Games         = make(map[string]*models.Game)
	DeveloperMode = true // true = admin, false = all users
)

func UploadGame() {
	var game map[string]*models.Game
	data, err := os.ReadFile(GamesSavePath)
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
	err = os.WriteFile(GamesSavePath, dataSave, 0600)
	if err != nil {
		log.Fatalln(err)
	}
}
