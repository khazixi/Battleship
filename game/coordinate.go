package util

import "errors"

type Coordinate interface {
	getCoordinate() (Point, error)
}

type Point struct {
	X, Y int
}

func (p Point) getCoordinate() (Point, error) {
	if p.X < 0 || p.X > 9 {
		return Point{}, errors.New("Invalid Coorindate Placement")
	} else if p.Y < 0 || p.Y > 9 {
		return Point{}, errors.New("Invalid Coorindate Placement")
	}
	return p, nil
}
