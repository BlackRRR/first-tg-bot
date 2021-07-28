package assets

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

const (
	GamesSavePath      = "assets/game_save.json"
	DefaultFieldSize   = 8
	DefaultBombCounter = 12
)

var (
	Games = make(map[string]*Game)
)

type Game struct {
	PlayingField       [DefaultFieldSize][DefaultFieldSize]string
	OpenedButtonsField [DefaultFieldSize][DefaultFieldSize]bool
	MessageID          int
}

func (g *Game) FillField() {
	for i := 0; i < DefaultBombCounter; i++ {
		g.deployingBomb()
	}
	for i := 0; i < DefaultFieldSize; i++ {
		for j := 0; j < DefaultFieldSize; j++ {
			if g.PlayingField[i][j] == "bomb" {
				continue
			}
			g.PlayingField[i][j] = strconv.Itoa(bombAround(g.PlayingField, i, j))
		}
	}
}

func (g *Game) deployingBomb() {
	var flag bool
	for !flag {
		cell := rand.Intn(DefaultFieldSize * DefaultFieldSize)
		row := cell % 8
		column := cell / 8

		if g.PlayingField[row][column] == "" {
			g.PlayingField[row][column] = "bomb"
			flag = true
		}
	}
}

func bombAround(field [DefaultFieldSize][DefaultFieldSize]string, i, j int) int {
	var counter int
	for k := -1; k < 2; k++ {
		if i+k < 0 || i+k > DefaultFieldSize-1 {
			continue
		}
		for l := -1; l < 2; l++ {
			if j+l < 0 || j+l > DefaultFieldSize-1 {
				continue
			}
			if field[i+k][j+l] == "bomb" {
				counter++
			}
		}
	}
	return counter
}

func UploadGame() {
	var game map[string]*Game
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
