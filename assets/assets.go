package assets //TODO: assets should contain all sorts of text variables and the like, there should not be any methods

import (
	"math/rand"
	"strconv"
)

const (
	GamesSavePath = "assets/game_save.json"
)

var (
	Games       = make(map[string]*Game)
	Size        int //TODO: this variable should not be global, it is unique for each game, so it should be the fields of the Game structure
	BombCounter int
)

type Game struct {
	PlayingField       [][]string
	OpenedButtonsField [][]bool
	MessageID          int
}

func (g *Game) FillEmptyField() {
	var field [][]string
	var open [][]bool

	for i := 0; i < Size; i++ {
		field = append(field, []string{})
		for j := 0; j < Size; j++ {
			field[i] = append(field[i], " ")
		}
	}
	for i := 0; i < Size; i++ {
		open = append(open, []bool{})
		for j := 0; j < Size; j++ {
			open[i] = append(open[i], false)
		}
		g.PlayingField = field //TODO: this action should not be performed several times, it should be performed after for, after line 33 and line 41
		g.OpenedButtonsField = open
	}
}

func (g *Game) FillField() {
	for i := 0; i < BombCounter; i++ {
		g.deployingBomb()
	}
	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
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
		cell := rand.Intn(Size * Size)
		row := cell % Size
		column := cell / Size
		if g.PlayingField[row][column] == " " {
			g.PlayingField[row][column] = "bomb"
			flag = true
		}
	}
}

func bombAround(field [][]string, i, j int) int {
	var counter int
	for k := -1; k < 2; k++ {
		if i+k < 0 || i+k > Size-1 {
			continue
		}
		for l := -1; l < 2; l++ {
			if j+l < 0 || j+l > Size-1 {
				continue
			}
			if field[i+k][j+l] == "bomb" {
				counter++
			}
		}
	}
	return counter
}
