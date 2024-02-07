package main

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

const (
	squaresPerRow = 11
	squareSize    = 32
	pieceRadius   = squareSize / 2
)

type PieceKind int8

const (
	None PieceKind = iota
	Pawn
	King
)

type Square struct {
	Piece PieceKind

	HasPiece bool

	BgColor    raylib.Color
	PieceColor raylib.Color
	BandColor  raylib.Color

	Radius int32
	Size   int32
	X      int32
	Y      int32
}

func NewSquare(bgColor raylib.Color, x int32, y int32) Square {
	return Square{None, false, bgColor, raylib.Blank, raylib.Blank, pieceRadius, squareSize, x, y}
}

func (s *Square) AddPiece(piece PieceKind, pieceColor raylib.Color, bandColor raylib.Color) {
	s.HasPiece = true
	s.Piece = piece
	s.PieceColor = pieceColor
	s.BandColor = bandColor
}

type Board struct {
	Squares [squaresPerRow][squaresPerRow]Square
}

func NewBoard() Board {
	var b Board
	var toggle raylib.Color = raylib.Beige

	for i := int32(0); i < squaresPerRow; i++ {
		for j := int32(0); j < squaresPerRow; j++ {
			if toggle == raylib.Beige {
				toggle = raylib.Brown
			} else {
				toggle = raylib.Beige
			}

			b.Squares[i][j] = NewSquare(toggle, i*squareSize, j*squareSize)
			b.Squares[i][j].BandColor = toggle
		}
	}

	// Black Pieces
	for i := 3; i < 8; i++ {
		b.Squares[i][0].AddPiece(Pawn, raylib.Black, raylib.White)
		b.Squares[0][i].AddPiece(Pawn, raylib.Black, raylib.White)
		b.Squares[i][squaresPerRow-1].AddPiece(Pawn, raylib.Black, raylib.White)
		b.Squares[squaresPerRow-1][i].AddPiece(Pawn, raylib.Black, raylib.White)

	}

	b.Squares[5][1].AddPiece(Pawn, raylib.Black, raylib.White)
	b.Squares[1][5].AddPiece(Pawn, raylib.Black, raylib.White)
	b.Squares[5][squaresPerRow-2].AddPiece(Pawn, raylib.Black, raylib.White)
	b.Squares[squaresPerRow-2][5].AddPiece(Pawn, raylib.Black, raylib.White)

	// White Pieces
	for i := 4; i < 7; i++ {
		for j := 4; j < 7; j++ {
			b.Squares[i][j].AddPiece(Pawn, raylib.White, raylib.Black)
		}
	}

	b.Squares[5][5].AddPiece(King, raylib.Yellow, raylib.Black)
	b.Squares[5][3].AddPiece(Pawn, raylib.White, raylib.Black)
	b.Squares[3][5].AddPiece(Pawn, raylib.White, raylib.Black)
	b.Squares[5][squaresPerRow-4].AddPiece(Pawn, raylib.White, raylib.Black)
	b.Squares[squaresPerRow-4][5].AddPiece(Pawn, raylib.White, raylib.Black)

	return b
}

type Game struct {
	FrameCounter int32
	ScreenWidth  int32
	ScreenHeight int32

	Win        bool
	FirstTurn  bool
	FirstClick bool

	GameBoard Board

	Turn raylib.Color

	Selected     *Square
	PrevSelected *Square
}

func (g *Game) Init() {
	g.ScreenWidth = squareSize*squaresPerRow + 1
	g.ScreenHeight = squareSize*squaresPerRow + 1
	g.FrameCounter = 0

	g.Win = false
	g.FirstTurn = true
	g.FirstClick = true

	g.GameBoard = NewBoard()

	g.Turn = raylib.Black
}

func (g *Game) SelectSquare() {
	g.PrevSelected = g.Selected

	x := raylib.GetMouseX()
	y := raylib.GetMouseY()
	x -= x % squareSize
	y -= y % squareSize

	row := x % squaresPerRow
	col := y % squaresPerRow

	if row != 0 {
		row = squaresPerRow - row
	}

	if col != 0 {
		col = squaresPerRow - col
	}

	if row < squaresPerRow &&
		col < squaresPerRow &&
		row > -1 && col > -1 {
		g.Selected = &g.GameBoard.Squares[row][col]
	} else {
		g.Selected = nil
	}
}

func (g *Game) Update() {
	if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
		g.SelectSquare()
	}
	g.FrameCounter++
}

func (g *Game) Draw() {
	raylib.BeginDrawing()
	if g.FirstTurn {
		for i := 0; i < squaresPerRow; i++ {
			for j := 0; j < squaresPerRow; j++ {
				s := &g.GameBoard.Squares[i][j]

				raylib.DrawRectangle(s.X, s.Y, squareSize, squareSize, s.BgColor)

				if s.Piece != None {
					raylib.DrawCircle(s.X+pieceRadius, s.Y+pieceRadius, pieceRadius, s.PieceColor)
					raylib.DrawEllipseLines(s.X+pieceRadius, s.Y+pieceRadius, float32(pieceRadius), float32(pieceRadius), s.BandColor)
				}
			}
		}

		g.FirstTurn = false
	}

	if g.Selected != nil {
		raylib.DrawEllipseLines(g.Selected.X+pieceRadius, g.Selected.Y+pieceRadius, float32(pieceRadius), float32(pieceRadius), raylib.Green)

		if g.PrevSelected != nil {
			raylib.DrawEllipseLines(g.PrevSelected.X+pieceRadius, g.PrevSelected.Y+pieceRadius, float32(pieceRadius), float32(pieceRadius), g.PrevSelected.BandColor)
		}
	}

	raylib.EndDrawing()
}

func main() {
	game := Game{}
	game.Init()

	raylib.InitWindow(int32(game.ScreenWidth), int32(game.ScreenHeight), "Hnefatafl")
	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		game.Draw()
		game.Update()
	}

	raylib.CloseWindow()
}
