package games

import (
	"fmt"
	"math"
	"strings"
)

const poc_gridSize = 8

var POC_PLAYER_SYMBOLS = [2]string{"White", "Black"}

type POC struct {
	board         []int
	currentPlayer int
	gameStatus    int
	initialPawns  map[int]bool
}

func NewPOCGame() Game {
	board := make([]int, poc_gridSize*poc_gridSize)
	initialPawns := make(map[int]bool)
	for j := 0; j < poc_gridSize*poc_gridSize; j++ {
		if j < 16 {
			board[j] = PLAYER1
			initialPawns[j] = true
		}
		if j >= 48 {
			board[j] = PLAYER2
			initialPawns[j] = true
		}
	}
	pocGame := &POC{board: board, currentPlayer: PLAYER1, gameStatus: IN_PROGRESS, initialPawns: initialPawns}
	return pocGame
}

func (poc *POC) UpdateCurrentPlayer() int {
	return poc.currentPlayer
}

func (poc *POC) GetCurrentPlayer() int {
	return poc.currentPlayer
}

func (poc *POC) GetCurrentPlayerSymbol() string {
	return POC_PLAYER_SYMBOLS[poc.currentPlayer-2]
}

func (poc *POC) GetOtherPlayerSymbol() string {
	return POC_PLAYER_SYMBOLS[(poc.currentPlayer-1)%2]
}

func (poc *POC) GetInitData() string {
	return ""
}

func (poc *POC) PerformMove(coords ...int) (int, []string) {
	return poc.performMove(coords[0], coords[1], coords[2], coords[3])
}

func (poc *POC) performMove(fromRow, fromCol, toRow, toCol int) (int, []string) {
	currPlayer := poc.currentPlayer
	// make move only for the correct current player
	if poc.board[fromRow*poc_gridSize+fromCol] == currPlayer {
		if poc.isValid(fromRow, fromCol, toRow, toCol) {
			poc.board[fromRow*poc_gridSize+fromCol] = 0

			// if it a diagonal first move then current player captures the opponents pawn
			if (poc.board[toRow*poc_gridSize+toCol] == 0) &&
				(poc.board[toRow*poc_gridSize+fromCol] != currPlayer) {
				poc.board[toRow*poc_gridSize+fromCol] = 0
			}

			poc.board[toRow*poc_gridSize+toCol] = currPlayer

			// delete the pawn from the initialPawns list
			delete(poc.initialPawns, fromRow*poc_gridSize+fromCol)

			data := make([]string, 0)
			for j := 0; j < poc_gridSize*poc_gridSize; j++ {
				data = append(data, fmt.Sprintf("%d", poc.board[j]))
			}

			movesForWhite, movesForBlack := poc.checkStatus()
			if currPlayer == PLAYER1 {
				if movesForBlack {
					poc.currentPlayer = PLAYER2
				}
			} else {
				if movesForWhite {
					poc.currentPlayer = PLAYER1
				}
			}

			return poc.gameStatus, []string{strings.Join(data, " ")}
		}
		return -1, make([]string, 0)
	}
	return -1, make([]string, 0)
}

func (poc *POC) isValid(fromRow, fromCol, toRow, toCol int) bool {
	if !poc.checkBounds(fromRow, fromCol, toRow, toCol) {
		return false
	}
	startPoint := poc_gridSize*fromRow + fromCol
	endPoint := poc_gridSize*toRow + toCol

	player, otherPlayer := poc.board[startPoint], poc.board[endPoint]
	// white pawn can only move downwards
	if player == PLAYER1 && (toRow-fromRow > 0) {
		return poc.isValidCheck(player, otherPlayer, startPoint, endPoint)
	}
	// black pawn can only move upwords
	if player == PLAYER2 && (toRow-fromRow < 0) {
		return poc.isValidCheck(player, otherPlayer, startPoint, endPoint)
	}
	return false
}

func (poc *POC) isValidCheck(player, otherPlayer, startPoint, endPoint int) bool {
	fromRow, fromCol := startPoint/poc_gridSize, startPoint%poc_gridSize
	toRow, toCol := endPoint/poc_gridSize, endPoint%poc_gridSize

	if fromCol == toCol {
		// check if the pawn is allowed to take two steps(only in it's first move)
		if math.Abs(float64(fromRow)-float64(toRow)) == 2 {
			if ok := poc.initialPawns[startPoint]; ok {
				betweenCell := 0
				if player == PLAYER1 {
					betweenCell = poc.board[poc_gridSize*(fromRow+1)+fromCol]
				}
				if player == PLAYER2 {
					betweenCell = poc.board[poc_gridSize*(fromRow-1)+fromCol]
				}
				// make sure it cannot jump over/on other pawn
				return (betweenCell == 0) && (otherPlayer == 0)
			}
			return false
		}
		return otherPlayer == 0
	} else {
		// checking for diagonal move to an empty square (allowed only if it is it's first move)
		if otherPlayer == 0 {
			if ok := poc.initialPawns[startPoint]; ok {
				p := poc.board[poc_gridSize*toRow+fromCol]
				if p != 0 && p != player {
					return true
				}
				return false
			}
			return false
		}
		return otherPlayer != player
	}
}

func (poc *POC) checkBounds(fromRow, fromCol, toRow, toCol int) bool {
	return (fromRow >= 0) && (fromRow < poc_gridSize) &&
		(fromCol >= 0) && (fromCol < poc_gridSize) &&
		(toRow >= 0) && (toRow < poc_gridSize) &&
		(toCol >= 0) && (toCol < poc_gridSize)
}

func (poc *POC) checkStatus() (bool, bool) {
	// PLAYER1 is white and PLAYER2 is black
	for col := 0; col < 8; col++ {
		// check if white has reached the respective end row
		if poc.board[poc_gridSize*7+col] == PLAYER1 {
			poc.gameStatus = PLAYER1
			return false, false
		}
		// check if black has reached the respective end row
		if poc.board[col] == PLAYER2 {
			poc.gameStatus = PLAYER2
			return false, false
		}
	}

	movesAvailableForWhite := false
	movesAvailableForBlack := false

	//for all the currently occupied cells check if the corresponding pawns can make a move
	for row := 0; row <= 7; row++ {
		if movesAvailableForWhite && movesAvailableForBlack {
			break
		}
		for col := 0; col <= 7; col++ {
			if movesAvailableForWhite && movesAvailableForBlack {
				break
			}
			player := poc.board[row*poc_gridSize+col]
			if player == PLAYER1 {
				whiteHasMoves := poc.hasMoves(player, row, col)
				if whiteHasMoves {
					movesAvailableForWhite = whiteHasMoves
				}

			}
			if player == PLAYER2 {
				blackHasMoves := poc.hasMoves(player, row, col)
				if blackHasMoves {
					movesAvailableForBlack = blackHasMoves
				}
			}
		}
	}

	if movesAvailableForWhite || movesAvailableForBlack {
		poc.gameStatus = IN_PROGRESS
	} else {
		poc.gameStatus = DRAW
	}
	return movesAvailableForWhite, movesAvailableForBlack
}

func (poc *POC) hasMoves(player, fromRow, fromCol int) bool {
	var dirs = []int{-1, 0, 1}
	if player == PLAYER1 {
		for _, d := range dirs {
			if poc.isValid(fromRow, fromCol, fromRow+1, fromCol+d) {
				return true
			}
		}
		// if the pawn is initial, it has the ability to take two steps froward
		if poc.isValid(fromRow, fromCol, fromRow+2, fromCol) {
			return true
		}
		return false
	} else {
		for _, d := range dirs {
			if poc.isValid(fromRow, fromCol, fromRow-1, fromCol+d) {
				return true
			}
		}
		// if the pawn is initial, it has the ability to take two steps froward
		if poc.isValid(fromRow, fromCol, fromRow-2, fromCol) {
			return true
		}
		return false
	}
}
