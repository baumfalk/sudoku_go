package main

import (
	"baumfalk/sudoku/game"
	"fmt"
)

func main() {
	// read sudoku
	sudoku := game.LoadSudoku("s05a.txt")
	solution := game.SolveSudoku(sudoku)
	fmt.Println("----------------------------------------------------")
	fmt.Println(sudoku)
	fmt.Println(solution)

}
