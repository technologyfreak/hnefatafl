package squares

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	piece "github.com/technologyfreak/hnefatafl/piece"
)

const (
	SquaresPerRow = 11
	SquareSize    = 32
)

type Square struct {
	Piece piece.PieceKind

	BgColor    raylib.Color
	PieceColor raylib.Color
	BandColor  raylib.Color

	X int32
	Y int32
}

func ToRowOrCol(i int32) int32 {
	i -= i % SquareSize
	i %= SquaresPerRow

	if i != 0 {
		i = SquaresPerRow - i
	}

	return i
}

func InRowRange(n int32) bool {
	return n >= 0 && n < SquaresPerRow
}

func NewSquare(bgColor raylib.Color, bandColor raylib.Color, x int32, y int32) Square {
	return Square{piece.None, bgColor, raylib.Blank, bandColor, x, y}
}

func (s *Square) AddPiece(piece piece.PieceKind, pieceColor raylib.Color, bandColor raylib.Color) {
	s.Piece = piece
	s.PieceColor = pieceColor
	s.BandColor = bandColor
}

func (s *Square) RemovePiece() {
	s.Piece = piece.None
	s.PieceColor = raylib.Blank
	s.BandColor = raylib.Blank
}

func (s *Square) HasPiece() bool {
	return s.Piece != piece.None
}

func (s1 *Square) IsWestOf(s2 *Square) bool {
	return (s1.X - s2.X) < 0
}

func (s1 *Square) IsEastOf(s2 *Square) bool {
	return (s1.X - s2.X) > 0
}

func (s1 *Square) IsNorthOf(s2 *Square) bool {
	return (s1.Y - s2.Y) < 0
}

func (s1 *Square) IsSouthOf(s2 *Square) bool {
	return (s1.Y - s2.Y) > 0
}

func (s *Square) IsCenter() bool {
	return ToRowOrCol(s.X) == 5 && ToRowOrCol(s.Y) == 5
}

func (s *Square) IsKingsCorner() bool {
	x := ToRowOrCol(s.X)
	y := ToRowOrCol(s.Y)

	return (x == 0 && y == 0) ||
		(x == 0 && y == 10) ||
		(x == 10 && y == 0) ||
		(x == 10 && y == 10)
}
