package main

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"math/rand"
)

type Cell struct {
	possibleValues mapset.Set[uint8]
	value          uint8
}

type Sodoku struct {
	Cells [][]Cell
}

func CreateGame() Sodoku {
	cells := make([][]Cell, 9)
	for i := range cells {
		cells[i] = make([]Cell, 9)
	}

	game := Sodoku{Cells: cells}
	game.InitCells()

	return game
}

func (s *Sodoku) InitCells() {
	for i, row := range s.Cells {
		for j, _ := range row {
			s.Cells[i][j].possibleValues = mapset.NewSet[uint8](1, 2, 3, 4, 5, 6, 7, 8, 9)
			s.Cells[i][j].value = 0 // 0 = not collapsed
		}
	}
}

func (s *Sodoku) PrintBoardWithEntropy() {
	fmt.Println("_______________________________________")
	for i, _ := range s.Cells {
		for j, _ := range s.Cells[i] {
			fmt.Print("|")
			if s.Cells[i][j].possibleValues != nil {
				for v := range s.Cells[i][j].possibleValues.Iter() {
					fmt.Print(v)
				}
			}
			fmt.Print("|")
		}
		fmt.Println()
	}
}

func (s *Sodoku) PrintBoard() {
	fmt.Println("_______________________")
	for i, row := range s.Cells {
		fmt.Print("|")
		for j, c := range row {
			if j > 0 && j%3 == 0 {
				fmt.Print(" |")
			}
			if c.value == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print(c.value)
			}
			fmt.Print("|")
		}
		fmt.Println()
		if (i+1)%3 == 0 {
			fmt.Println("=======================")
		}
	}
	fmt.Println()
	fmt.Println()
}

func (s *Sodoku) findMinimumEntropy() (int, int) {
	minEntropy := 11
	minEntropyRow := -1
	minEntropyCol := -1

	for i, row := range s.Cells {
		for j, cell := range row {
			if cell.possibleValues != nil && cell.value == 0 {
				entropy := cell.possibleValues.Cardinality()
				if entropy < minEntropy {
					minEntropy = entropy
					minEntropyRow = i
					minEntropyCol = j
				}
			}
		}
	}

	return minEntropyRow, minEntropyCol
}

func (s *Sodoku) propagateValue(row int, col int, value uint8) {
	// propagate horizontally
	for i := 0; i < 9; i++ {
		if s.Cells[row][i].possibleValues != nil {
			s.Cells[row][i].possibleValues.Remove(value)
		}
	}

	// propagate vertically
	for i := 0; i < 9; i++ {
		if s.Cells[i][col].possibleValues != nil {
			s.Cells[i][col].possibleValues.Remove(value)
		}
	}

	// propagate within "box"
	startRow := row / 3 * 3
	startCol := col / 3 * 3
	endRow := startRow + 2
	endCol := startCol + 2
	for r := startRow; r <= endRow; r++ {
		for c := startCol; c <= endCol; c++ {
			if s.Cells[r][c].possibleValues != nil {
				s.Cells[r][c].possibleValues.Remove(value)
			}
		}
	}
}

func (s *Sodoku) Collapse() (bool, error) {
	minEntropyRow, minEntropyCol := s.findMinimumEntropy()
	foundCollapsable := false

	if minEntropyCol >= 0 && minEntropyRow >= 0 {
		// collapse the cell with the minimum entropy
		cell := &s.Cells[minEntropyRow][minEntropyCol]
		if cell.value == 0 {
			v, _ := cell.possibleValues.Pop()
			if v != 0 {
				cell.value = v
				foundCollapsable = true
			}
		}

		cell.possibleValues = nil

		// propagate the value
		s.propagateValue(minEntropyRow, minEntropyCol, cell.value)

		if !foundCollapsable {
			return false, fmt.Errorf("failed to fully collapse board")
		}

		return true, nil
	}

	return false, nil
}

func (s *Sodoku) ClearSlots(numSlots int) {
	for i := 0; i < numSlots; i++ {
		s.Cells[rand.Int()%9][rand.Int()%9].value = 0
	}
}

func main() {
	game := CreateGame()

	//game.PrintBoard()
	notDone, err := game.Collapse()
	for notDone {
		notDone, err = game.Collapse()

		if err != nil {
			game = CreateGame()
			notDone = true
		}

		//game.PrintBoard()
		//game.PrintBoardWithEntropy()
	}
	game.PrintBoard()

	//game.ClearSlots(50)
	//game.PrintBoard()
}
