package games

import (
	"fmt"
	"math"
)

const tttt_gridSize = 3

var TTTT_PLAYER_SYMBOLS = [2]string{"X", "O"}

type TTTTGame struct {
	grids         [][]int
	gameStatus    []int
	currentPlayer int
}

func NewTTTTGame() Game {
	numberOfGrids := 3

	gameStatus := []int{numberOfGrids, 0, 0, 0}

	grids := make([][]int, numberOfGrids)

	for i := 0; i < numberOfGrids; i++ {
		grid := make([]int, tttt_gridSize*tttt_gridSize)
		for j := 0; j < tttt_gridSize*tttt_gridSize; j++ {
			grid[j] = IN_PROGRESS
		}
		grids[i] = grid
	}

	tttGame := &TTTTGame{grids: grids, gameStatus: gameStatus, currentPlayer: PLAYER1}

	return tttGame
}

func (tttt *TTTTGame) GetCurrentPlayer() int {
	return tttt.currentPlayer
}

func (tttt *TTTTGame) GetCurrentPlayerSymbol() string {
	return TTTT_PLAYER_SYMBOLS[tttt.currentPlayer-2]
}

func (tttt *TTTTGame) GetOtherPlayerSymbol() string {
	return TTTT_PLAYER_SYMBOLS[(tttt.currentPlayer-1)%2]
}

func (tttt *TTTTGame) UpdateCurrentPlayer() int {

	prevPlayer := tttt.currentPlayer
	if prevPlayer == PLAYER1 {
		tttt.currentPlayer = PLAYER2
	} else {
		tttt.currentPlayer = PLAYER1
	}
	return tttt.currentPlayer
}

func (tttt *TTTTGame) GetInitData() string {
	return ""
}

func (tttt *TTTTGame) PerformMove(coords ...int) (int, []string) {
	return tttt.performMove(coords[0], coords[1], coords[2])
}

func (tttt *TTTTGame) performMove(gridNumber, row, col int) (int, []string) {

	status := tttt.checkGridStatus(gridNumber)
	grid := tttt.grids[gridNumber]
	if status == IN_PROGRESS && grid[row*tttt_gridSize+col] == IN_PROGRESS {
		grid[row*tttt_gridSize+col] = tttt.currentPlayer
		fmt.Printf("playerNo: %d, status: %d\n", tttt.currentPlayer, tttt.checkGridStatus(gridNumber))
		tttt.printBoard()
		gridStatus, gameStatus := tttt.checkStatus(gridNumber)

		gamedata := fmt.Sprintf("%d %d %d", gridNumber, row, col)
		if gridStatus == PLAYER1 {
			gamedata = fmt.Sprintf("%s %d", gamedata, PLAYER1)
		} else if gridStatus == PLAYER2 {
			gamedata = fmt.Sprintf("%s %d", gamedata, PLAYER2)
		}

		return gameStatus, []string{gamedata}
	}
	fmt.Printf("playerNo: %d, status: %d\n", tttt.currentPlayer, status)

	return -1, make([]string, 0)
}

// consider gameStatus as a collection of 4 buckets
// gameStatus[0] corresponds to number of grids that are IN_PROGRESS
// gameStatus[1] corresponds to number of grids that are DRAW
// gameStatus[2] corresponds to number of grids that are won by PLAYER1
// gameStatus[3] corresponds to number of grids that are won by PLAYER2
// every time the status of grid changes from IN_PROGRESS to status_x gameStatus[0] is decremented by 1 and gameStatus[status_x] is incremented by 1
func (tttt *TTTTGame) checkStatus(gridNumber int) (int, int) {

	gridStatus := tttt.checkGridStatus(gridNumber)

	if gridStatus != IN_PROGRESS {
		tttt.gameStatus[0] -= 1
		tttt.gameStatus[gridStatus] += 1
	}

	switch tttt.gameStatus[0] {
	case 0:
		if tttt.gameStatus[2] > tttt.gameStatus[3] {
			return gridStatus, PLAYER1
		} else if tttt.gameStatus[2] < tttt.gameStatus[3] {
			return gridStatus, PLAYER2
		} else {
			return gridStatus, DRAW
		}

	case 1:
		player1Score := float64(tttt.gameStatus[2])
		player2Score := float64(tttt.gameStatus[3])
		if tttt.gameStatus[1] == 0 && int(math.Abs(player1Score-player2Score)) == 2 {
			if player1Score > player2Score {
				return gridStatus, PLAYER1
			} else {
				return gridStatus, PLAYER2
			}
		} else {
			return gridStatus, IN_PROGRESS
		}
	default:
		return gridStatus, IN_PROGRESS
	}
}

func (tttt *TTTTGame) checkGridStatus(gridNumber int) int {
	grid := tttt.grids[gridNumber]

	diag1 := make([]int, tttt_gridSize)
	diag2 := make([]int, tttt_gridSize)

	// checking rows
	for i := 0; i < tttt_gridSize; i++ {
		row := grid[tttt_gridSize*i : tttt_gridSize*(i+1)]
		col := make([]int, tttt_gridSize)
		for j := 0; j < tttt_gridSize; j++ {
			col[j] = grid[j*tttt_gridSize+i]
		}
		checkRowForWin := tttt.checkForWin(row)
		if checkRowForWin == 2 || checkRowForWin == 3 {
			return checkRowForWin
		}
		checkColForWin := tttt.checkForWin(col)
		if checkColForWin == 2 || checkColForWin == 3 {
			return checkColForWin
		}

		diag1[i] = grid[(i*tttt_gridSize)+i]
		diag2[i] = grid[((2-i)*tttt_gridSize)+i]
	}

	checkDiag1ForWin := tttt.checkForWin(diag1)
	if checkDiag1ForWin == 2 || checkDiag1ForWin == 3 {
		return checkDiag1ForWin
	}
	checkDiag2ForWin := tttt.checkForWin(diag2)
	if checkDiag2ForWin == 2 || checkDiag2ForWin == 3 {
		return checkDiag2ForWin
	}

	if tttt.isFull(gridNumber) {
		return DRAW
	}

	return IN_PROGRESS
}

func (tttt *TTTTGame) checkForWin(row []int) int {

	isEqual := true
	prev := row[0]

	for i := 0; i < tttt_gridSize; i++ {
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

func (tttt *TTTTGame) isFull(gridNumber int) bool {

	for i := 0; i < tttt_gridSize*tttt_gridSize; i++ {
		if tttt.grids[gridNumber][i] == 0 {
			return false
		}
	}

	return true
}

func (tttt *TTTTGame) printBoard() {
	for i := 0; i < 3; i++ {
		fmt.Println(fmt.Sprintf("Board %d: ", i), tttt.grids[i])
	}

}
