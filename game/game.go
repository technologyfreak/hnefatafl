package game

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	board "github.com/technologyfreak/hnefatafl/board"
	piece "github.com/technologyfreak/hnefatafl/piece"
	square "github.com/technologyfreak/hnefatafl/square"
)

type Game struct {
	FrameCounter int32
	ScreenWidth  int32
	ScreenHeight int32

	BlackPawns uint8
	WhitePawns uint8
	MovePhase  uint8

	Win       bool
	FirstTurn bool

	Board board.Board

	Turn raylib.Color

	Selected     *square.Square
	PrevSelected *square.Square
}

func (g *Game) Init() {
	g.ScreenWidth = square.SquareSize*square.SquaresPerRow + 1
	g.ScreenHeight = square.SquareSize*square.SquaresPerRow + 1
	g.FrameCounter = 0

	g.BlackPawns = 0
	g.WhitePawns = 0
	g.MovePhase = 0

	g.Win = false
	g.FirstTurn = true

	g.Board = board.NewBoard()

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

	row = square.ToRowOrCol(row)
	col = square.ToRowOrCol(col)

	if square.InRowRange(row) && square.InRowRange(col) {
		g.Selected = &g.Board.Squares[row][col]
		g.MovePhase++
	}
}

func (g *Game) TryMoveWest(wPiece *square.Square, blank *square.Square) int32 {
	if blank.IsWestOf(wPiece) &&
		wPiece.Y == blank.Y {

		x := wPiece.X - square.SquareSize
		col := wPiece.Y

		col = square.ToRowOrCol(col)

		for x > blank.X && x > 0 {
			selected := &g.Board.Squares[square.ToRowOrCol(x)][col]

			if selected.HasPiece() {
				return -1
			}

			x -= square.SquareSize
		}

		g.MovePhase++
		return 0
	}

	return -1
}

func (g *Game) TryMoveEast(wPiece *square.Square, blank *square.Square) int32 {
	if blank.IsEastOf(wPiece) &&
		wPiece.Y == blank.Y {

		x := wPiece.X + square.SquareSize
		col := wPiece.Y

		col = square.ToRowOrCol(col)

		for x < blank.X && square.ToRowOrCol(x) < square.SquaresPerRow {
			selected := &g.Board.Squares[square.ToRowOrCol(x)][col]

			if selected.HasPiece() {
				return -1
			}

			x += square.SquareSize
		}

		g.MovePhase++
		return 0
	}

	return -1
}

func (g *Game) TryMoveNorth(wPiece *square.Square, blank *square.Square) int32 {
	if blank.IsNorthOf(wPiece) &&
		wPiece.X == blank.X {

		y := wPiece.Y - square.SquareSize
		row := wPiece.X

		row = square.ToRowOrCol(row)

		for y > blank.Y && y > 0 {
			selected := &g.Board.Squares[row][square.ToRowOrCol(y)]

			if selected.HasPiece() {
				return -1
			}

			y -= square.SquareSize
		}

		g.MovePhase++
		return 0
	}

	return -1
}

func (g *Game) TryMoveSouth(wPiece *square.Square, blank *square.Square) int32 {
	if blank.IsSouthOf(wPiece) &&
		wPiece.X == blank.X {

		y := wPiece.Y + square.SquareSize
		row := wPiece.X

		row = square.ToRowOrCol(row)

		for y < blank.Y && square.ToRowOrCol(y) < square.SquaresPerRow {
			selected := &g.Board.Squares[row][square.ToRowOrCol(y)]

			if selected.HasPiece() {
				return -1
			}

			y += square.SquareSize
		}

		g.MovePhase++
		return 0
	}

	return -1
}

func (g *Game) ValidateMove() {
	if g.MovePhase == 1 {
		if g.Selected.HasPiece() &&
			g.Selected.PieceColor == g.Turn {
			g.MovePhase++
		} else {
			g.MovePhase = 0
		}
	}

	if g.MovePhase == 3 {
		if !g.Selected.HasPiece() {
			if (g.Selected.IsKingsCorner() || g.Selected.IsCenter()) &&
				g.PrevSelected.Piece != piece.King {
				g.MovePhase = 0
			}

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

func (g *Game) KingHasReachedACorner() bool {
	return (g.Board.Squares[0][0].Piece == piece.King) ||
		(g.Board.Squares[0][10].Piece == piece.King) ||
		(g.Board.Squares[10][0].Piece == piece.King) ||
		(g.Board.Squares[10][10].Piece == piece.King)
}

func (g *Game) Update() {
	if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
		g.SelectSquare()
		g.ValidateMove()

		// Swap piece to new location
		if g.MovePhase == 5 {
			g.Selected.AddPiece(g.PrevSelected.Piece, g.PrevSelected.PieceColor, g.PrevSelected.BandColor)
			g.PrevSelected.RemovePiece()

			if g.Turn == raylib.Black {
				g.Turn = raylib.White
			} else {
				g.Turn = raylib.Black
			}

			g.MovePhase++
		}
	}

	if g.KingHasReachedACorner() {
		g.Win = true
	}

	g.FrameCounter++
}

func (g *Game) Draw() {
	raylib.BeginDrawing()
	if g.FirstTurn {
		for i := 0; i < square.SquaresPerRow; i++ {
			for j := 0; j < square.SquaresPerRow; j++ {
				s := &g.Board.Squares[i][j]

				raylib.DrawRectangle(s.X, s.Y, square.SquareSize, square.SquareSize, s.BgColor)

				if s.HasPiece() {
					raylib.DrawCircle(s.X+piece.PieceRadius, s.Y+piece.PieceRadius, piece.PieceRadius, s.PieceColor)
					raylib.DrawEllipseLines(s.X+piece.PieceRadius, s.Y+piece.PieceRadius, float32(piece.PieceRadius), float32(piece.PieceRadius), s.BandColor)

					if s.Piece == piece.King {
						raylib.DrawCircle(s.X+piece.PieceRadius, s.Y+piece.PieceRadius, piece.PieceRadius/2, raylib.Gold)
					}
				}
			}
		}

		g.FirstTurn = false
	}

	if g.MovePhase == 2 {
		// TODO - Add back selected piece highlighting
	}

	if g.MovePhase == 6 {
		raylib.DrawRectangle(g.PrevSelected.X, g.PrevSelected.Y, square.SquareSize, square.SquareSize, g.PrevSelected.BgColor)
		raylib.DrawCircle(g.Selected.X+piece.PieceRadius, g.Selected.Y+piece.PieceRadius, float32(piece.PieceRadius), g.Selected.PieceColor)
		raylib.DrawEllipseLines(g.Selected.X+piece.PieceRadius, g.Selected.Y+piece.PieceRadius, float32(piece.PieceRadius), float32(piece.PieceRadius), g.Selected.BandColor)

		if g.Selected.Piece == piece.King {
			raylib.DrawCircle(g.Selected.X+piece.PieceRadius, g.Selected.Y+piece.PieceRadius, piece.PieceRadius/2, raylib.Gold)
		}

		g.MovePhase = 0
	}

	raylib.EndDrawing()
}
