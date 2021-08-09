package models

import (
	"math/rand"
	"strconv"
)

type Game struct {
	Key                string
	PlayingField       [][]string
	OpenedButtonsField [][]bool
	MessageID          int
	Size               int
	BombCounter        int
}

func (g *Game) FillEmptyField() {
	var field [][]string
	var open [][]bool

	for i := 0; i < g.Size; i++ {
		field = append(field, []string{})
		for j := 0; j < g.Size; j++ {
			field[i] = append(field[i], " ")
		}
	}
	g.PlayingField = field

	for i := 0; i < g.Size; i++ {
		open = append(open, []bool{})
		for j := 0; j < g.Size; j++ {
			open[i] = append(open[i], false)
		}
	}
	g.OpenedButtonsField = open
}

func (g *Game) FillField() {
	for i := 0; i < g.BombCounter; i++ {
		g.deployingBomb()
	}
	for i := 0; i < g.Size; i++ {
		for j := 0; j < g.Size; j++ {
			if g.PlayingField[i][j] == "bomb" {
				continue
			}
			g.PlayingField[i][j] = strconv.Itoa(bombAround(g.PlayingField, i, j, g.Size))
		}
	}
}

func (g *Game) deployingBomb() {
	var flag bool
	for !flag {
		cell := rand.Intn(g.Size * g.Size)
		row := cell % g.Size
		column := cell / g.Size
		if g.PlayingField[row][column] == " " {
			g.PlayingField[row][column] = "bomb"
			flag = true
		}
	}
}

func bombAround(field [][]string, i, j int, Size int) int {
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
