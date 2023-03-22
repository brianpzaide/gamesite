package games

const (
	IN_PROGRESS = 0
	DRAW        = 1
	PLAYER1     = 2
	PLAYER2     = 3
)

type Game interface {
	PerformMove(coords ...int) (int, []string)
	UpdateCurrentPlayer() int
	GetCurrentPlayer() int
	GetCurrentPlayerSymbol() string
	GetOtherPlayerSymbol() string
	GetInitData() string
}
