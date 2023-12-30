package util

type Turn = bool

const (
  PLAYER1 Turn = false
  PLAYER2 Turn = true
)

type Board [10][10]byte

type Game struct {
  p1 Board
  p2 Board
  PlayerTurn Turn
}
