package game

import (
	"testing"
)

func TestPoint(t *testing.T) {
	point := Point{10, 10}

	_, err := point.getCoordinate()
	if err == nil {
		t.Fatal("Points should not be encodable for values > 10")
	}

}

func TestInstruction(t *testing.T) {
	in1 := Instruction{'A', '0'}
	in1m, err := ParseInstruction("A0")

	if err != nil {
		t.Fatal("Failed to properly parse the instruction")
	}

	if in1m != in1 {
		t.Fatal("Parsed Instruction should be equivalent to the actual function")
	}
}

func TestBoard(t *testing.T) {
	var board Board

	for ridx, row := range board {
		for _, col := range row {
			if col != 0 {
				t.Fatalf("board at Point{%d, %d}\n", col, ridx)
			}
		}
	}

	err := board.Place(Point{0, 0}, PATROLBOAT, RIGHT)
	if err != nil {
		t.Fatal("Placement should not fail")
	}

  for i := 0; i < int(PATROLBOAT.Size()); i++ {
    if board[0][i] != 1 {
      t.Fatal("Failed to properly mark spot")
    }
  }

  err = board.Place(Point{0,0}, SUBMARINE, RIGHT)
  if err == nil {
    t.Log(board)
    t.Fatal("Placement should not be possible")
  }
}
