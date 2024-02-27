package board

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	piece "github.com/technologyfreak/hnefatafl/piece"
	square "github.com/technologyfreak/hnefatafl/square"
)

type Board struct {
	Squares [square.SquaresPerRow][square.SquaresPerRow]square.Square
}

func NewBoard() Board {
	var b Board
	var toggle raylib.Color = raylib.Beige

	for row := int32(0); row < square.SquaresPerRow; row++ {
		for col := int32(0); col < square.SquaresPerRow; col++ {
			if toggle == raylib.Beige {
				toggle = raylib.Brown
			} else {
				toggle = raylib.Beige
			}

			b.Squares[row][col] = square.NewSquare(toggle, row*square.SquareSize, col*square.SquareSize)
		}
	}

	// Black Pieces
	for row := 3; row < 8; row++ {
		b.Squares[row][0].AddPiece(piece.BlackPawn)
		b.Squares[0][row].AddPiece(piece.BlackPawn)
		b.Squares[row][square.SquaresPerRow-1].AddPiece(piece.BlackPawn)
		b.Squares[square.SquaresPerRow-1][row].AddPiece(piece.BlackPawn)

	}

	b.Squares[5][1].AddPiece(piece.BlackPawn)
	b.Squares[1][5].AddPiece(piece.BlackPawn)
	b.Squares[5][square.SquaresPerRow-2].AddPiece(piece.BlackPawn)
	b.Squares[square.SquaresPerRow-2][5].AddPiece(piece.BlackPawn)

	// White Pieces
	for row := 4; row < 7; row++ {
		for col := 4; col < 7; col++ {
			b.Squares[row][col].AddPiece(piece.WhitePawn)
		}
	}

	b.Squares[5][5].AddPiece(piece.King | piece.WhitePawn)
	b.Squares[5][3].AddPiece(piece.WhitePawn)
	b.Squares[3][5].AddPiece(piece.WhitePawn)
	b.Squares[5][square.SquaresPerRow-4].AddPiece(piece.WhitePawn)
	b.Squares[square.SquaresPerRow-4][5].AddPiece(piece.WhitePawn)

	return b
}
