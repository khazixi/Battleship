package game

import "errors"

type Coordinate interface {
	getCoordinate() (Point, error)
}

type Instruction [2]byte

type Point struct {
	Y, X int
}

func (p Point) getCoordinate() (Point, error) {
	if p.X < 0 || p.X > 9 {
		return Point{}, errors.New("Invalid Coorindate Placement")
	} else if p.Y < 0 || p.Y > 9 {
		return Point{}, errors.New("Invalid Coorindate Placement")
	}
	return p, nil
}

func (i Instruction) getCoordinate() (Point, error) {
	if ('A' <= i[0] && i[0] <= 'J') && ('0' <= i[1] && i[1] <= '9') {
		return Point{X: int(i[1] - '0'), Y: int(i[1] - 'A')}, nil
	}
	return Point{}, errors.New("Invalid Insturction")
}

func ValidInstruction(s string) bool {
	switch {
	case s[0] < 'A':
		fallthrough
	case s[0] > 'J':
		fallthrough
	case s[1] < '0':
		fallthrough
	case s[1] > '9':
		fallthrough
	case len(s) != 2:
		return false
	default:
		return true
	}
}

func ParseInstruction(s string) (Instruction, error) {
	switch {
	case s[0] < 'A':
		fallthrough
	case s[0] > 'J':
		fallthrough
	case s[1] < '0':
		fallthrough
	case s[1] > '9':
		fallthrough
	case len(s) != 2:
		return [2]byte{}, errors.New("Failed to properly parse the Instruction")
	default:
		return [2]byte{s[0], s[1]}, nil
	}
}
