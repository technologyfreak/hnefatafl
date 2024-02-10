package main

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

const (
	squaresPerRow = 11
	squareSize    = 32
	pieceRadius   = squareSize / 2
)

func toRowOrCol(i int32) int32 {
	i -= i % squareSize
	i %= squaresPerRow

	if i != 0 {
		i = squaresPerRow - i
	}

	return i
}

type PieceKind int8

const (
	None PieceKind = iota
	Pawn
	King
)

type Square struct {
	Piece PieceKind

	BgColor    raylib.Color
	PieceColor raylib.Color
	BandColor  raylib.Color

	X int32
	Y int32
}

func NewSquare(bgColor raylib.Color, bandColor raylib.Color, x int32, y int32) Square {
	return Square{None, bgColor, raylib.Blank, bandColor, x, y}
}

func (s *Square) AddPiece(piece PieceKind, pieceColor raylib.Color, bandColor raylib.Color) {
	s.Piece = piece
	s.PieceColor = pieceColor
	s.BandColor = bandColor
}

func (s *Square) RemovePiece() {
	s.Piece = None
	s.PieceColor = raylib.Blank
	s.BandColor = raylib.Blank
}

func (s *Square) HasPiece() bool {
	return s.Piece != None
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

type Board struct {
	Squares [squaresPerRow][squaresPerRow]Square
}

func NewBoard() Board {
	var b Board
	var toggle raylib.Color = raylib.Beige

	for row := int32(0); row < squaresPerRow; row++ {
		for col := int32(0); col < squaresPerRow; col++ {
			if toggle == raylib.Beige {
				toggle = raylib.Brown
			} else {
				toggle = raylib.Beige
			}

			b.Squares[row][col] = NewSquare(toggle, toggle, row*squareSize, col*squareSize)
		}
	}

	// Black Pieces
	for row := 3; row < 8; row++ {
		b.Squares[row][0].AddPiece(Pawn, raylib.Black, raylib.White)
		b.Squares[0][row].AddPiece(Pawn, raylib.Black, raylib.White)
		b.Squares[row][squaresPerRow-1].AddPiece(Pawn, raylib.Black, raylib.White)
		b.Squares[squaresPerRow-1][row].AddPiece(Pawn, raylib.Black, raylib.White)

	}

	b.Squares[5][1].AddPiece(Pawn, raylib.Black, raylib.White)
	b.Squares[1][5].AddPiece(Pawn, raylib.Black, raylib.White)
	b.Squares[5][squaresPerRow-2].AddPiece(Pawn, raylib.Black, raylib.White)
	b.Squares[squaresPerRow-2][5].AddPiece(Pawn, raylib.Black, raylib.White)

	// White Pieces
	for row := 4; row < 7; row++ {
		for col := 4; col < 7; col++ {
			b.Squares[row][col].AddPiece(Pawn, raylib.White, raylib.Black)
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

	MovePhase uint8

	Win       bool
	FirstTurn bool

	GameBoard Board

	Turn raylib.Color

	Selected     *Square
	PrevSelected *Square
}

func (g *Game) Init() {
	g.ScreenWidth = squareSize*squaresPerRow + 1
	g.ScreenHeight = squareSize*squaresPerRow + 1
	g.FrameCounter = 0

	g.MovePhase = 0

	g.Win = false
	g.FirstTurn = true

	g.GameBoard = NewBoard()

	g.Turn = raylib.Black
}

func (g *Game) SelectSquare() {
	if g.MovePhase == 0 {
		g.Selected = nil
		g.PrevSelected = nil
	} else {
		g.PrevSelected = g.Selected
	}

	row := raylib.GetMouseX()
	col := raylib.GetMouseY()

	row = toRowOrCol(row)
	col = toRowOrCol(col)

	if row < squaresPerRow &&
		col < squaresPerRow &&
		row > -1 &&
		col > -1 {
		g.Selected = &g.GameBoard.Squares[row][col]
		g.MovePhase++
	}
}

func (g *Game) TryMoveWest(wPiece *Square, blank *Square) int32 {
	if blank.IsWestOf(wPiece) &&
		wPiece.Y == blank.Y {

		x := wPiece.X - squareSize
		col := wPiece.Y

		col = toRowOrCol(col)

		for x > blank.X && x > 0 {
			selected := &g.GameBoard.Squares[toRowOrCol(x)][col]

			if selected.HasPiece() {
				return -1
			}

			x -= squareSize
		}

		g.MovePhase++
		return 0
	}

	return -1
}

func (g *Game) TryMoveEast(wPiece *Square, blank *Square) int32 {
	if blank.IsEastOf(wPiece) &&
		wPiece.Y == blank.Y {

		x := wPiece.X + squareSize
		col := wPiece.Y

		col = toRowOrCol(col)

		for x < blank.X && toRowOrCol(x) < squaresPerRow {
			selected := &g.GameBoard.Squares[toRowOrCol(x)][col]

			if selected.HasPiece() {
				return -1
			}

			x += squareSize
		}

		g.MovePhase++
		return 0
	}

	return -1
}

func (g *Game) TryMoveNorth(wPiece *Square, blank *Square) int32 {
	if blank.IsNorthOf(wPiece) &&
		wPiece.X == blank.X {

		y := wPiece.Y - squareSize
		row := wPiece.X

		row = toRowOrCol(row)

		for y > blank.Y && y > 0 {
			selected := &g.GameBoard.Squares[row][toRowOrCol(y)]

			if selected.HasPiece() {
				return -1
			}

			y -= squareSize
		}

		g.MovePhase++
		return 0
	}

	return -1
}

func (g *Game) TryMoveSouth(wPiece *Square, blank *Square) int32 {
	if blank.IsSouthOf(wPiece) &&
		wPiece.X == blank.X {

		y := wPiece.Y + squareSize
		row := wPiece.X

		row = toRowOrCol(row)

		for y < blank.Y && toRowOrCol(y) < squaresPerRow {
			selected := &g.GameBoard.Squares[row][toRowOrCol(y)]

			if selected.HasPiece() {
				return -1
			}

			y += squareSize
		}

		g.MovePhase++
		return 0
	}

	return -1
}

func (g *Game) ValidateMove() {
	if g.MovePhase == 1 {
		if g.Selected.HasPiece() {
			g.MovePhase++
		} else {
			g.MovePhase = 0
		}
	}

	if g.MovePhase == 3 {
		if !g.Selected.HasPiece() {
			g.MovePhase++
		} else {
			g.MovePhase = 0
		}
	}

	if g.MovePhase == 4 {
		if g.TryMoveWest(g.PrevSelected, g.Selected) == -1 &&
			g.TryMoveEast(g.PrevSelected, g.Selected) == -1 &&
			g.TryMoveNorth(g.PrevSelected, g.Selected) == -1 &&
			g.TryMoveSouth(g.PrevSelected, g.Selected) == -1 {
			g.MovePhase = 0
		}
	}
}

func (g *Game) Update() {
	if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
		g.SelectSquare()
		g.ValidateMove()

		// Swap piece to no location
		if g.MovePhase == 5 {
			g.Selected.AddPiece(g.PrevSelected.Piece, g.PrevSelected.PieceColor, g.PrevSelected.BandColor)
			g.PrevSelected.RemovePiece()
			g.MovePhase++
		}
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

				if s.HasPiece() {
					raylib.DrawCircle(s.X+pieceRadius, s.Y+pieceRadius, pieceRadius, s.PieceColor)
					raylib.DrawEllipseLines(s.X+pieceRadius, s.Y+pieceRadius, float32(pieceRadius), float32(pieceRadius), s.BandColor)
				}
			}
		}

		g.FirstTurn = false
	}

	if g.MovePhase == 2 {
		// TODO - Add back selected piece highlighting
	}

	if g.MovePhase == 6 {
		raylib.DrawRectangle(g.PrevSelected.X, g.PrevSelected.Y, squareSize, squareSize, g.PrevSelected.BgColor)
		raylib.DrawCircle(g.Selected.X+pieceRadius, g.Selected.Y+pieceRadius, float32(pieceRadius), g.Selected.PieceColor)
		raylib.DrawEllipseLines(g.Selected.X+pieceRadius, g.Selected.Y+pieceRadius, float32(pieceRadius), float32(pieceRadius), g.Selected.BandColor)

		g.MovePhase = 0
	}

	raylib.EndDrawing()
}

func main() {
	game := Game{}
	game.Init()

	raylib.InitWindow(int32(game.ScreenWidth), int32(game.ScreenHeight), "Hnefatafl")
	raylib.SetTargetFPS(60)

	for !raylib.WindowShouldClose() {
		game.Update()
		game.Draw()
	}

	raylib.CloseWindow()
}
