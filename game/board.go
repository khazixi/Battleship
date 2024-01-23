package game

import (
	"errors"
)

type Turn bool

type Piece byte

type Direction int

const (
	PLAYER1 Turn = false
	PLAYER2 Turn = true
)

const (
	// 5 Spaces Long
	CARRIER Piece = iota

	// 4 Spaces Long
	BATTLESHIP

	// 3 Spaces Long
	DESTROYER

	// 3 Spaces Long
	SUBMARINE

	// 2 Spaces Long
	PATROLBOAT
)

const (
	UP Direction = iota
	LEFT
	DOWN
	RIGHT
)

type Board [10][10]byte

type Game struct {
	P1         Board
	P2         Board
	PlayerTurn Turn
}

type Transmit struct {
	Piece
	Direction
	Coordinate
}

func (p Piece) Size() byte {
	switch p {
	case CARRIER:
		return 5
	case BATTLESHIP:
		return 4
	case DESTROYER:
		return 3
	case SUBMARINE:
		return 3
	case PATROLBOAT:
		return 2
	default: // NOTE: This invariant should not be reached
		return 0
	}
}

func (b *Board) Reset() {
	for i := 0; i < 10; i++ {
		for ii := 0; ii < 10; ii++ {
			b[i][ii] = 0
		}
	}
}

func (b *Board) Validate(coordinate Coordinate, p Piece, d Direction) error {
	point, err := coordinate.getCoordinate()
	if err != nil {
		return err
	}
	switch d {
	case LEFT:
		if p.Size()-1 > point.X {
			return errors.New("Placement of Piece is out of range")
		} else {
			for cur := point.X; cur > p.Size()-point.X-1; cur-- {
				if b[point.Y][point.X] == 1 {
					return errors.New("Other Boat Detected")
				}
			}
		}

	case RIGHT:
		if p.Size()+point.X-1 > 9 {
			return errors.New("Other Boat Detected")
		} else {
			for cur := point.X; cur < p.Size()+point.X; cur++ {
				if b[point.Y][point.X] == 1 {
					return errors.New("Other Boat Detected")
				}
			}
		}

	case UP:
		if p.Size()-1 > point.Y {
			return errors.New("Other Boat Detected")
		} else {
			for cur := point.Y; cur > p.Size()-point.Y; cur-- {
				if b[point.Y][point.X] == 1 {
					return errors.New("Other Boat Detected")
				}
			}
		}

	case DOWN:
		if p.Size()+point.Y-1 > 9 {
			return errors.New("Other Boat Detected")
		} else {
			for cur := point.Y; cur < p.Size()+point.Y; cur++ {
				if b[point.Y][point.X] == 1 {
					return errors.New("Other Boat Detected")
				}
			}

		}
	}

	return nil
}

func (b *Board) Place(coordinate Coordinate, p Piece, d Direction) error {
	point, err := coordinate.getCoordinate()
	if err != nil {
		return err
	}

	err = b.Validate(coordinate, p, d)

	if err != nil {
		return err
	}

	// WARNING: Hopefully this is correct I don't know
	switch d {
	case LEFT:
		if point.X-p.Size()+1 < 0 {
			return errors.New("Piece Cannot be placed in this direction")
		}

		for i := point.X; i > point.X-p.Size(); i-- {
			b[point.Y][i] = 1
		}

	case RIGHT:
		for i := point.X; i < point.X+p.Size(); i++ {
			b[point.Y][i] = 1
		}
	case DOWN:
		for i := point.Y; i > point.Y-p.Size(); i-- {
			b[i][point.X] = 1
		}
	case UP:
		for i := point.Y; i < point.Y+p.Size(); i++ {
			b[i][point.X] = 1
		}
	default:
		return errors.New("Invalid Ditection")
	}

	return nil
}

func (b *Board) Mark(coordinate Coordinate) bool {
	point, err := coordinate.getCoordinate()
	if err != nil {
		return false
	}

	b[point.Y][point.X] += 2

	return b[point.Y][point.X] == 3
}

func (b Board) HasWin() bool {
	count := 0
	for i := 0; i < 10; i++ {
		for ii := 0; ii < 10; ii++ {
			if b[i][ii] == 3 {
				count++
			}
		}
	}

	return count == 17
}

func CreateGame() *Game {
	return &Game{
		P1:         [10][10]byte{},
		P2:         [10][10]byte{},
		PlayerTurn: PLAYER1,
	}
}
