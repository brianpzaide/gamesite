package games

import "fmt"

const nttt_gridSize = 3

var NTTT_PLAYER_SYMBOLS = [2]string{"X", "O"}

type NTTTGame struct {
	grids         [][]int
	currentPlayer int
}

func NewNTTTGame() Game {

	numberOfGrids := 10

	grids := make([][]int, numberOfGrids)

	for i := 0; i < numberOfGrids; i++ {
		grid := make([]int, nttt_gridSize*nttt_gridSize)
		for j := 0; j < nttt_gridSize*nttt_gridSize; j++ {
			grid[j] = IN_PROGRESS
		}
		grids[i] = grid
	}

	nttGame := &NTTTGame{grids: grids, currentPlayer: PLAYER1}

	return nttGame
}

func (nttt *NTTTGame) GetCurrentPlayer() int {
	return nttt.currentPlayer
}

func (nttt *NTTTGame) GetCurrentPlayerSymbol() string {
	return NTTT_PLAYER_SYMBOLS[nttt.currentPlayer-2]
}

func (nttt *NTTTGame) GetOtherPlayerSymbol() string {
	return NTTT_PLAYER_SYMBOLS[(nttt.currentPlayer-1)%2]
}

func (nttt *NTTTGame) UpdateCurrentPlayer() int {

	prevPlayer := nttt.currentPlayer
	if prevPlayer == PLAYER1 {
		nttt.currentPlayer = PLAYER2
	} else {
		nttt.currentPlayer = PLAYER1
	}
	return nttt.currentPlayer
}

func (nttt *NTTTGame) GetInitData() string {
	return ""
}

func (nttt *NTTTGame) PerformMove(coords ...int) (int, []string) {
	return nttt.performMove(coords[0], coords[1], coords[2])
}

func (nttt *NTTTGame) performMove(gridNumber, row, col int) (int, []string) {
	status := nttt.checkGridStatus(gridNumber)
	grid := nttt.grids[gridNumber]
	fmt.Printf("NTTT gridNo: %d row: %d, col: %d\n", gridNumber, row, status)
	if status == IN_PROGRESS && grid[row*nttt_gridSize+col] == IN_PROGRESS {
		grid[row*nttt_gridSize+col] = nttt.currentPlayer
		fmt.Printf("playerNo: %d, status: %d\n", nttt.currentPlayer, nttt.checkGridStatus(gridNumber))
		nttt.printBoard()
		gridStatus, gameStatus := nttt.checkStatus(gridNumber)
		gamedata := fmt.Sprintf("%d %d %d", gridNumber, row, col)
		if gridStatus == PLAYER1 {
			gamedata = fmt.Sprintf("%s %d", gamedata, PLAYER1)
		} else if gridStatus == PLAYER2 {
			gamedata = fmt.Sprintf("%s %d", gamedata, PLAYER2)
		}
		return gameStatus, []string{gamedata}
	}
	fmt.Printf("playerNo: %d, status: %d\n", nttt.currentPlayer, status)
	return -1, make([]string, 0)
}

func (nttt *NTTTGame) checkStatus(gridNumber int) (int, int) {
	mainGridNo := 9

	gridStatus := nttt.checkGridStatus(gridNumber)

	col := gridNumber % 3
	row := (gridNumber - col) / 3

	if gridStatus != IN_PROGRESS {

		nttt.grids[mainGridNo][row*nttt_gridSize+col] = gridStatus
	}

	return gridStatus, nttt.checkGridStatus(mainGridNo)

}

func (nttt *NTTTGame) checkGridStatus(gridNumber int) int {
	grid := nttt.grids[gridNumber]

	diag1 := make([]int, nttt_gridSize)
	diag2 := make([]int, nttt_gridSize)

	// checking rows
	for i := 0; i < nttt_gridSize; i++ {
		row := grid[nttt_gridSize*i : nttt_gridSize*(i+1)]
		col := make([]int, nttt_gridSize)
		for j := 0; j < nttt_gridSize; j++ {
			col[j] = grid[j*nttt_gridSize+i]
		}
		checkRowForWin := nttt.checkForWin(row)
		if checkRowForWin == 2 || checkRowForWin == 3 {
			return checkRowForWin
		}
		checkColForWin := nttt.checkForWin(col)
		if checkColForWin == 2 || checkColForWin == 3 {
			return checkColForWin
		}

		diag1[i] = grid[(i*nttt_gridSize)+i]
		diag2[i] = grid[((2-i)*nttt_gridSize)+i]
	}

	checkDiag1ForWin := nttt.checkForWin(diag1)
	if checkDiag1ForWin == 2 || checkDiag1ForWin == 3 {
		return checkDiag1ForWin
	}
	checkDiag2ForWin := nttt.checkForWin(diag2)
	if checkDiag2ForWin == 2 || checkDiag2ForWin == 3 {
		return checkDiag2ForWin
	}

	if nttt.isFull(gridNumber) {
		return DRAW
	}

	return IN_PROGRESS
}

func (nttt *NTTTGame) checkForWin(row []int) int {

	isEqual := true
	prev := row[0]

	for i := 0; i < nttt_gridSize; i++ {
		if prev != row[i] {
			isEqual = false
			break
		}
	}

	if isEqual {
		return prev
	}

	return IN_PROGRESS
}

func (nttt *NTTTGame) isFull(gridNumber int) bool {

	for i := 0; i < nttt_gridSize*nttt_gridSize; i++ {
		if nttt.grids[gridNumber][i] == 0 {
			return false
		}
	}

	return true
}

func (nttt *NTTTGame) printBoard() {
	for i := 0; i < 10; i++ {
		fmt.Println(fmt.Sprintf("Board %d: ", i), nttt.grids[i])
	}

}
