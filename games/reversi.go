package games

import (
	"fmt"
	"strings"
)

const reversi_gridSize = 8

var REVERSI_PLAYER_SYMBOLS = [2]string{"White", "Black"}

// reversi_dir represents list of directions starting from right at index 0 to downright at index 7
// right 	 = 0
// left 	 = 1
// up 		 = 2
// down 	 = 3
// upright 	 = 4
// downleft  = 5
// upleft 	 = 6
// downright = 7
var reversi_dir = [8]int{0, 1, 2, 3, 4, 5, 6, 7}

// dirRate represents howm to increment/decrement row and col to move in a particular direction
var reversi_dirRate = map[int][2]int{
	//moving right
	//increment row by 0 and col by 1
	0: {0, 1},

	//moving left
	//increment row by 0 and col by -1
	1: {0, -1},

	//moving up
	//increment row by -1 and col by 0
	2: {-1, 0},

	//moving down
	//increment row by 1 and col by 0
	3: {1, 0},

	//moving upright
	//increment row by -1 and col by 1
	4: {-1, 1},

	//moving downleft
	//increment row by 1 and col by -1
	5: {1, -1},

	//moving upleft
	//increment row by -1 and col by -1
	6: {-1, -1},

	//moving downright
	//increment row by 1 and col by 1
	7: {1, 1},
}

type Reversi struct {
	board         []int
	currentPlayer int
	gameStatus    int
	score         map[int]int
}

type reversi_cell struct {
	R, C int
}

func NewReversiGame() Game {

	board := make([]int, reversi_gridSize*reversi_gridSize)
	for j := 0; j < reversi_gridSize*reversi_gridSize; j++ {
		if j == 27 || j == 36 {
			board[j] = PLAYER1
		} else if j == 28 || j == 35 {
			board[j] = PLAYER2
		} else {
			board[j] = IN_PROGRESS
		}
	}
	reversiGame := &Reversi{board: board, currentPlayer: PLAYER1, gameStatus: IN_PROGRESS, score: map[int]int{PLAYER1: 2, PLAYER2: 2}}

	return reversiGame
}

func (ri *Reversi) UpdateCurrentPlayer() int {
	return ri.currentPlayer
}

func (ri *Reversi) GetCurrentPlayer() int {
	return ri.currentPlayer
}

func (ri *Reversi) GetCurrentPlayerSymbol() string {
	return REVERSI_PLAYER_SYMBOLS[ri.currentPlayer-2]
}

func (ri *Reversi) GetOtherPlayerSymbol() string {
	return REVERSI_PLAYER_SYMBOLS[(ri.currentPlayer-1)%2]
}

func (ri *Reversi) GetInitData() string {
	return ""
}

func (ri *Reversi) PerformMove(coords ...int) (int, []string) {
	return ri.performMove(coords[0], coords[1])
}

func (ri *Reversi) performMove(row, col int) (int, []string) {
	if ri.board[row*reversi_gridSize+col] == IN_PROGRESS {
		allowed := false

		currPlayer := ri.currentPlayer
		opponent := 0
		if currPlayer == PLAYER1 {
			opponent = PLAYER2
		} else {
			opponent = PLAYER1
		}

		// c stores direction as key and slice of cells in that direction as value
		c := make(map[int][]reversi_cell)

		for _, dir := range reversi_dir {
			ok, cellsAlongDir := ri.checkInDir(ri.currentPlayer, opponent, row, col, dir)
			if ok {
				c[dir] = cellsAlongDir
				allowed = true
			}
		}

		if !allowed {
			return -1, make([]string, 0)
		}

		ri.board[row*reversi_gridSize+col] = currPlayer
		count := 1
		for _, cells := range c {
			for _, cell := range cells {
				ri.board[cell.R*reversi_gridSize+cell.C] = ri.currentPlayer
				count++
			}
		}

		ri.score[currPlayer] = ri.score[currPlayer] + count
		ri.score[opponent] = ri.score[opponent] - count + 1

		data := make([]string, 0)
		for j := 0; j < reversi_gridSize*reversi_gridSize; j++ {
			data = append(data, fmt.Sprintf("%d", ri.board[j]))
		}

		movesForWhite, movesForBlack := ri.checkStatus()
		if ri.currentPlayer == PLAYER1 {
			if movesForBlack {
				ri.currentPlayer = PLAYER2
			}
		} else {
			if movesForWhite {
				ri.currentPlayer = PLAYER1
			}
		}

		return ri.gameStatus, []string{fmt.Sprintf("%d-%d", ri.score[PLAYER1], ri.score[PLAYER2]), strings.Join(data, " ")}
	}
	return -1, make([]string, 0)
}

func (ri *Reversi) checkInDir(player, opponent, row, col, dir int) (bool, []reversi_cell) {
	validMove := false

	r := row
	c := col
	row_rate := reversi_dirRate[dir][0]
	col_rate := reversi_dirRate[dir][1]
	cells := make([]reversi_cell, 0)
	initial := true
	r = r + row_rate
	c = c + col_rate
	for r >= 0 && r <= 7 && c >= 0 && c <= 7 {
		if ri.board[r*reversi_gridSize+c] == player {
			if !initial {
				validMove = true
			}
			return validMove, cells
		} else if ri.board[r*reversi_gridSize+c] == opponent {
			rprime := r
			cprime := c
			cells = append(cells, reversi_cell{R: rprime, C: cprime})
			initial = false
			r = r + row_rate
			c = c + col_rate
		} else {
			return validMove, cells
		}
	}

	return validMove, cells
}

func (ri *Reversi) checkStatus() (bool, bool) {
	movesAvailableForWhite := false
	movesAvailableForBlack := false

	//for all the empty cells in the board check if black or white fits
	for row := 0; row <= 7; row++ {
		if movesAvailableForWhite && movesAvailableForBlack {
			break
		}
		for col := 0; col <= 7; col++ {
			if movesAvailableForWhite && movesAvailableForBlack {
				break
			}
			if ri.board[row*reversi_gridSize+col] == IN_PROGRESS {
				w := ri.isFits(PLAYER1, row, col)
				b := ri.isFits(PLAYER2, row, col)
				if w {
					movesAvailableForWhite = true
				}
				if b {
					movesAvailableForBlack = true
				}
			}

		}
	}

	if movesAvailableForWhite || movesAvailableForBlack {
		ri.gameStatus = IN_PROGRESS
	} else {
		if ri.score[PLAYER1] > ri.score[PLAYER2] {
			ri.gameStatus = PLAYER1
		} else if ri.score[PLAYER1] < ri.score[PLAYER2] {
			ri.gameStatus = PLAYER2
		} else {
			ri.gameStatus = DRAW
		}
		return false, false
	}
	return movesAvailableForWhite, movesAvailableForBlack
}

func (ri *Reversi) isFits(player, row, col int) bool {
	fits := false
	var opponent int
	if player == PLAYER1 {
		opponent = PLAYER2
	} else {
		opponent = PLAYER1
	}
	for _, dir := range reversi_dir {
		fits, _ = ri.checkInDir(player, opponent, row, col, dir)
		if fits {
			return true
		}
	}
	return false
}
