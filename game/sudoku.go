package game

import (
	"fmt"
	"io/ioutil"
)

type Sudoku [9][9]int
type SudokuPartial [9][9](map[int]bool)

func LoadSudoku(filename string) Sudoku {
	rawData, _ := ioutil.ReadFile(filename)
	rawSudokuData := make([]int, 9*9)
	var sudokuData Sudoku
	i := 0
	for _, val := range rawData {
		if val == 32 || val == 10 {
			continue
		}
		rawSudokuData[i] = int(val) - 48
		i++
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			sudokuData[i][j] = rawSudokuData[i*9+j]
		}
	}

	return sudokuData
}

func unsolved(sudoku SudokuPartial) bool {
	for _, row := range sudoku {
		for _, cell := range row {
			if len(cell) > 1 {
				return true
			}
		}
	}
	return false
}

func initializeSudokuSolver(sudoku Sudoku) SudokuPartial {
	var part SudokuPartial
	for x, row := range sudoku {
		for y, cell := range row {
			part[x][y] = make(map[int]bool)
			if cell != 0 {
				part[x][y][cell] = true
			} else {
				for i := 1; i <= 9; i++ {
					part[x][y][i] = true
				}
			}
		}
	}
	return part
}

func horSolHelper(sp *SudokuPartial, x int, y int) bool {
	changed := false
	for yCur := 0; yCur < 9; yCur++ {
		if yCur != y && len(sp[x][yCur]) == 1 {
			for key := range sp[x][yCur] {
				if sp[x][y][key] {
					delete(sp[x][y], key)
					changed = true
					fmt.Println("[Hor] For [", x, ",", y, "], we remove the value", key)
					fmt.Println(sp[x][y])
				}
			}
		}
	}
	return changed
}

func verSolHelper(sp *SudokuPartial, x int, y int) bool {
	changed := false

	for xCur := 0; xCur < 9; xCur++ {
		if xCur != x && len(sp[xCur][y]) == 1 {
			for key := range sp[xCur][y] {
				if sp[x][y][key] {
					delete(sp[x][y], key)
					changed = true

					fmt.Println("[Ver] For [", x, ",", y, "], we remove the value", key)
				}
			}
		}
	}
	return changed
}

func blockSolHelper(sp *SudokuPartial, x int, y int) bool {
	xBlock := int(x / 3)
	yBlock := int(y / 3)
	changed := false

	for xDiv := 0; xDiv < 3; xDiv++ {
		for yDiv := 0; yDiv < 3; yDiv++ {
			solvedCell := len(sp[xBlock*3+xDiv][yBlock*3+yDiv]) == 1
			sameCell := (xBlock*3+xDiv == x && yBlock*3+yDiv == y)
			if solvedCell && !sameCell {
				for key := range sp[xBlock*3+xDiv][yBlock*3+yDiv] {
					if sp[x][y][key] {
						delete(sp[x][y], key)
						changed = true

						fmt.Println("[Blk] For [", x, ",", y, "], we remove the value", key, " because it is already in [", xBlock*3+xDiv, ",", yBlock*3+yDiv, "]")
					}
				}
			}
		}
	}
	return changed
}

func horizontalSolve(sp *SudokuPartial) bool {
	return generalSolve(sp, horSolHelper)
}

func verticalSolve(sp *SudokuPartial) bool {
	return generalSolve(sp, verSolHelper)
}

func blockSolve(sp *SudokuPartial) bool {
	return generalSolve(sp, blockSolHelper)
}

func generalSolve(sp *SudokuPartial, f func(*SudokuPartial, int, int) bool) bool {
	changed := false

	for x, row := range sp {
		for y, cell := range row {
			if len(cell) > 1 {
				changed = changed || f(sp, x, y)
			}
		}
	}

	return changed
}

func printPartSolve(sp *SudokuPartial) {
	for i, row := range sp {
		for j, cell := range row {
			if len(cell) > 1 {
				fmt.Print("_")

			} else {
				for key := range cell {
					fmt.Print(key)
				}
			}

			if (j+1)%3 == 0 && j <= 5 {
				fmt.Print(" | ")

			} else {
				fmt.Print(" ")

			}
		}
		fmt.Println()
		if (i+1)%3 == 0 && i <= 5 {
			fmt.Println("---------------------")

		}
	}
}

func SolveSudoku(sudoku Sudoku) Sudoku {
	var sol Sudoku
	partSol := initializeSudokuSolver(sudoku)

	i := 1
	for unsolved(partSol) {
		fmt.Println("Iteration", i)

		improvement := horizontalSolve(&partSol) || verticalSolve(&partSol) || blockSolve(&partSol)
		if !improvement {
			fmt.Println("Currently unsolvable :(")
			break
		}
		i++
		printPartSolve(&partSol)
		fmt.Println()
	}

	return sol
}
