package games

import (
	"fmt"
	"math/rand"
	"strconv"
)

const (
	MIN            = 0
	MAX            = 10
	maxit_gridSize = 6
)

var MAXIT_PLAYER_SYMBOLS = [2]string{"V", "H"}

type MaxitGame struct {
	gameStatus int
	// stores the maxit board
	grid []cell
	// stores the scores of each player
	scores                                []int
	currentPlayer, currentRow, currentCol int
}
type cell struct {
	status, value int
}

// func init() {
// 	rand.Seed(time.Now().Unix())
// }

func randInt() int {
	return MIN + rand.Intn(MAX-MIN)
}

func NewMaxitGame() Game {
	scores := []int{0, 0}
	grid := make([]cell, maxit_gridSize*maxit_gridSize)
	for j := 0; j < maxit_gridSize*maxit_gridSize; j++ {
		grid[j].value = randInt()
	}
	maxitGame := &MaxitGame{grid: grid,
		gameStatus:    IN_PROGRESS,
		currentPlayer: PLAYER1,
		scores:        scores,
		currentRow:    0,
		currentCol:    0}
	return maxitGame
}

func (maxit *MaxitGame) GetCurrentPlayer() int {
	return maxit.currentPlayer
}

func (maxit *MaxitGame) GetCurrentPlayerSymbol() string {
	return MAXIT_PLAYER_SYMBOLS[maxit.currentPlayer-2]
}

func (maxit *MaxitGame) GetOtherPlayerSymbol() string {
	return MAXIT_PLAYER_SYMBOLS[(maxit.currentPlayer-1)%2]
}

func (maxit *MaxitGame) UpdateCurrentPlayer() int {
	return maxit.currentPlayer
}

func (maxit *MaxitGame) GetInitData() string {
	s := ""

	for i := 0; i < len(maxit.grid); i++ {
		s += fmt.Sprintf("%s ", strconv.Itoa(maxit.grid[i].value))
	}
	return s
}

func (maxit *MaxitGame) PerformMove(coords ...int) (int, []string) {

	return maxit.performMove(coords[0], coords[1])
}

func (maxit *MaxitGame) performMove(row, col int) (int, []string) {
	cell := row*maxit_gridSize + col
	if maxit.grid[cell].status == IN_PROGRESS {
		// if current player is V(vertical)
		if maxit.currentPlayer == PLAYER1 {
			if maxit.currentCol == col {
				// updating the scores
				maxit.scores[maxit.currentPlayer-2] += maxit.grid[cell].value
				// updating the cell status
				maxit.grid[cell].status = maxit.currentPlayer
				maxit.currentRow = row
			}
			// else currentPlayer is H(horizontal)
		} else {
			if maxit.currentRow == row {
				// updating the scores
				maxit.scores[maxit.currentPlayer-2] += maxit.grid[cell].value
				// updating the cell status
				maxit.grid[cell].status = maxit.currentPlayer
				maxit.currentCol = col
			}
		}
		status, nextPlayer := maxit.getStatus(row, col)
		if status == IN_PROGRESS {
			maxit.currentPlayer = nextPlayer
			maxit.gameStatus = status
			gamedata := fmt.Sprintf("%d %d", row, col)
			return status, []string{fmt.Sprintf("%ds%d", maxit.scores[PLAYER1-2], maxit.scores[PLAYER2-2]), gamedata}
			// return status, []string{gamedata}
		}

		return status, make([]string, 0)
	}
	return -1, make([]string, 0)
}

func (maxit *MaxitGame) getStatus(row, col int) (int, int) {
	status, nextPlayer := IN_PROGRESS, PLAYER1
	player1Moves, player2Moves := 0, 0
	for i := 0; i < maxit_gridSize; i++ {
		// checking the column if moves are available for vertical player
		if maxit.grid[i*maxit_gridSize+col].status == IN_PROGRESS {
			player1Moves += 1
		}
		// checking the row if moves are available for vertical player
		if maxit.grid[row*maxit_gridSize+i].status == IN_PROGRESS {
			player2Moves += 1
		}
	}
	// fmt.Println("player1Moves:", player1Moves, "player2Moves:", player2Moves)
	if player1Moves == 0 && player2Moves == 0 {
		// fmt.Println("both players have no moves")
		if maxit.scores[0] > maxit.scores[1] {
			status = PLAYER1
		} else if maxit.scores[0] < maxit.scores[1] {
			status = PLAYER2
		} else {
			status = DRAW
		}
		nextPlayer = -1
	} else if player1Moves == 0 {
		// fmt.Println("player1 has no moves")
		nextPlayer = PLAYER2
	} else if player2Moves == 0 {
		// fmt.Println("player2 has no moves")
		nextPlayer = PLAYER1
	} else {
		if maxit.currentPlayer == PLAYER1 {
			nextPlayer = PLAYER2
		} else {
			nextPlayer = PLAYER1
		}
	}

	return status, nextPlayer
}

// func (maxit *MaxitGame) printBoard() {
// 	for i := 0; i < maxit_gridSize; i++ {
// 		for j := 0; j < maxit_gridSize; j++ {
// 			fmt.Printf("%d %d", maxit.grid[i*maxit_gridSize+j].status, maxit.grid[i*maxit_gridSize+j].value)
// 		}
// 		fmt.Println()
// 	}
// }
