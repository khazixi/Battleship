package game

import "errors"

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

func (b *Board) Validate(coordinate Coordinate, p Piece, d Direction) bool {
	switch d {
	case LEFT:
		if p.Size() > coordinate.getX() {
			return false
		} else {
			for cur := coordinate.getX(); cur > p.Size()-coordinate.getX(); cur-- {
				if b[coordinate.getY()][coordinate.getX()] == 1 {
					return false
				}
			}
		}
	case RIGHT:
		if p.Size()+coordinate.getX() > 9 {
			return false
		} else {
			for cur := coordinate.getX(); cur < p.Size()+coordinate.getX(); cur++ {
				if b[coordinate.getY()][coordinate.getX()] == 1 {
					return false
				}
			}
		}

	case UP:
		if p.Size() > coordinate.getY() {
			return false
		} else {
			for cur := coordinate.getY(); cur > p.Size()-coordinate.getY(); cur-- {
				if b[coordinate.getY()][coordinate.getX()] == 1 {
					return false
				}
			}
		}

	case DOWN:
		if p.Size()+coordinate.getY() > 9 {
			return false
		} else {
			for cur := coordinate.getY(); cur < p.Size()+coordinate.getY(); cur++ {
				if b[coordinate.getY()][coordinate.getX()] == 1 {
					return false
				}
			}

		}
	}

	return true
}

func (b *Board) Place(coordinate Coordinate, p Piece, d Direction) error {
	point, err := coordinate.getCoordinate()
	if err != nil {
		return err
	}

  if !b.Validate(coordinate, p, d) {
    return errors.New("Invalid Placement for the board")
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
		if point.X+p.Size()-1 < 0 {
			return errors.New("Piece Cannot be placed in this direction")
		}

		for i := point.X; i < point.X+p.Size(); i++ {
			b[point.Y][i] = 1
		}
	case DOWN:
		if point.Y-p.Size()+1 < 0 {
			return errors.New("Piece Cannot be placed in this direction")
		}

		for i := point.Y; i > point.Y-p.Size(); i-- {
			b[i][point.X] = 1
		}
	case UP:
		if point.Y+p.Size()-1 < 0 {
			return errors.New("Piece Cannot be placed in this direction")
		}

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
