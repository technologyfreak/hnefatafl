package game

import (
	queue "github.com/eapache/queue/v2"
	raylib "github.com/gen2brain/raylib-go/raylib"
	board "github.com/technologyfreak/hnefatafl/board"
	piece "github.com/technologyfreak/hnefatafl/piece"
	square "github.com/technologyfreak/hnefatafl/square"
)

const (
	totalBlackPawns = 24
	totalWhitePawns = 12
)

type NeigborKind uint8

const (
	Unopposed NeigborKind = iota
	Edge
	KingsSquare
	Opposed
)

type CoordPair struct {
	X int32
	Y int32
}

type Game struct {
	FrameCounter int32
	FrameTarget  int32
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

	ReDraw *queue.Queue[*square.Square]
}

func (g *Game) Init() {
	g.FrameTarget = 60
	g.FrameCounter = 0
	g.ScreenWidth = square.SquareSize*square.SquaresPerRow + 1
	g.ScreenHeight = square.SquareSize*square.SquaresPerRow + 1

	g.BlackPawns = 0
	g.WhitePawns = 0
	g.MovePhase = 0

	g.Win = false
	g.FirstTurn = true

	g.Board = board.NewBoard()

	g.Turn = raylib.Black

	g.ReDraw = queue.New[*square.Square]()

	raylib.InitWindow(int32(g.ScreenWidth), int32(g.ScreenHeight), "Hnefatafl")
	raylib.SetTargetFPS(g.FrameTarget)

	for !raylib.WindowShouldClose() {
		g.Update()
		g.Draw()
	}

	raylib.CloseWindow()
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

		if !square.InRowRange(col) {
			return -1
		}

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

		if !square.InRowRange(col) {
			return -1
		}

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

		if !square.InRowRange(row) {
			return -1
		}

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

		if !square.InRowRange(row) {
			return -1
		}

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

func (g *Game) GetWesternNeighbor(wPiece *square.Square) (NeigborKind, *square.Square) {
	if wPiece == nil {
		return Edge, nil
	}

	west := CoordPair{X: square.ToRowOrCol(wPiece.X - square.SquareSize), Y: square.ToRowOrCol(wPiece.Y)}

	if square.InRowRange(west.X) && square.InRowRange(west.Y) {
		neigbor := &g.Board.Squares[west.X][west.Y]
		if neigbor.HasPiece() {
			if neigbor.PieceColor != wPiece.PieceColor {
				return Opposed, neigbor
			}
		} else if neigbor.IsKingsCorner() {
			return KingsSquare, neigbor
		}

		return Unopposed, neigbor
	}

	return Edge, nil
}

func (g *Game) GetEasternNeighbor(wPiece *square.Square) (NeigborKind, *square.Square) {
	if wPiece == nil {
		return Edge, nil
	}

	east := CoordPair{X: square.ToRowOrCol(wPiece.X + square.SquareSize), Y: square.ToRowOrCol(wPiece.Y)}

	if square.InRowRange(east.X) && square.InRowRange(east.Y) {
		neigbor := &g.Board.Squares[east.X][east.Y]
		if neigbor.HasPiece() {
			if neigbor.PieceColor != wPiece.PieceColor {
				return Opposed, neigbor
			}
		} else if neigbor.IsKingsCorner() {
			return KingsSquare, neigbor
		}

		return Unopposed, neigbor
	}

	return Edge, nil
}

func (g *Game) GetNortherNeighbor(wPiece *square.Square) (NeigborKind, *square.Square) {
	if wPiece == nil {
		return Edge, nil
	}

	north := CoordPair{X: square.ToRowOrCol(wPiece.X), Y: square.ToRowOrCol(wPiece.Y - square.SquareSize)}

	if square.InRowRange(north.X) && square.InRowRange(north.Y) {
		neigbor := &g.Board.Squares[north.X][north.Y]
		if neigbor.HasPiece() {
			if neigbor.PieceColor != wPiece.PieceColor {
				return Opposed, neigbor
			}
		} else if neigbor.IsKingsCorner() {
			return KingsSquare, neigbor
		}

		return Unopposed, neigbor
	}

	return Edge, nil
}

func (g *Game) GetSouternNeighbor(wPiece *square.Square) (NeigborKind, *square.Square) {
	if wPiece == nil {
		return Edge, nil
	}

	south := CoordPair{X: square.ToRowOrCol(wPiece.X), Y: square.ToRowOrCol(wPiece.Y + square.SquareSize)}

	if square.InRowRange(south.X) && square.InRowRange(south.Y) {
		neigbor := &g.Board.Squares[south.X][south.Y]
		if neigbor.HasPiece() {
			if neigbor.PieceColor != wPiece.PieceColor {
				return Opposed, neigbor
			}
		} else if neigbor.IsKingsCorner() {
			return KingsSquare, neigbor
		}

		return Unopposed, neigbor
	}

	return Edge, nil
}

func (g *Game) IsSandwiched(wPiece *square.Square) bool {
	if wPiece == nil {
		return false
	}

	westKind, westPiece := g.GetWesternNeighbor(wPiece)
	eastKind, eastPiece := g.GetEasternNeighbor(wPiece)
	northKind, northPiece := g.GetNortherNeighbor(wPiece)
	southKind, southPiece := g.GetSouternNeighbor(wPiece)

	TryIsCenter := func(s *square.Square) bool {
		return s != nil && s.IsCenter()
	}

	hasWest := westKind > 0
	hasEast := eastKind > 0
	hasNorth := northKind > 0
	hasSouth := southKind > 0

	if wPiece.Piece == piece.King {
		x := (hasWest || TryIsCenter(westPiece)) && (hasEast || TryIsCenter(eastPiece))
		y := (hasNorth || TryIsCenter(northPiece)) && (hasSouth || TryIsCenter(southPiece))

		return x && y
	}

	x := hasWest && hasEast
	y := hasNorth && hasSouth

	return x || y

}

func (g *Game) UpdateWesternNeighbor() {
	westKind, westSquare := g.GetWesternNeighbor(g.Selected)

	if westKind > 0 && g.IsSandwiched(westSquare) {
		if westSquare.Piece == piece.King {
			g.Win = true
		} else if westSquare.PieceColor == raylib.Black {
			g.BlackPawns--
		} else {
			g.WhitePawns--
		}

		westSquare.RemovePiece()
		g.ReDraw.Add(westSquare)
	}
}

func (g *Game) UpdateEasternNeighbor() {
	eastKind, eastSquare := g.GetEasternNeighbor(g.Selected)

	if eastKind > 0 && g.IsSandwiched(eastSquare) {
		if eastSquare.Piece == piece.King {
			g.Win = true
		} else if eastSquare.PieceColor == raylib.Black {
			g.BlackPawns--
		} else {
			g.WhitePawns--
		}

		eastSquare.RemovePiece()
		g.ReDraw.Add(eastSquare)
	}
}

func (g *Game) UpdateNorthernNeighbor() {
	northKind, northSquare := g.GetNortherNeighbor(g.Selected)

	if northKind > 0 && g.IsSandwiched(northSquare) {
		if northSquare.Piece == piece.King {
			g.Win = true
		} else if northSquare.PieceColor == raylib.Black {
			g.BlackPawns--
		} else {
			g.WhitePawns--
		}

		northSquare.RemovePiece()
		g.ReDraw.Add(northSquare)
	}
}

func (g *Game) UpdateSouthernNeighbor() {
	southKind, southSquare := g.GetSouternNeighbor(g.Selected)

	if southKind > 0 && g.IsSandwiched(southSquare) {
		if southSquare.Piece == piece.King {
			g.Win = true
		} else if southSquare.PieceColor == raylib.Black {
			g.BlackPawns--
		} else {
			g.WhitePawns--
		}

		southSquare.RemovePiece()
		g.ReDraw.Add(southSquare)
	}
}

func (g *Game) Update() {
	if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
		g.SelectSquare()
		g.ValidateMove()

		if g.MovePhase == 5 {
			// Swap piece to new location
			g.Selected.AddPiece(g.PrevSelected.Piece, g.PrevSelected.PieceColor, g.PrevSelected.BandColor)
			g.PrevSelected.RemovePiece()

			g.UpdateWesternNeighbor()
			g.UpdateEasternNeighbor()
			g.UpdateNorthernNeighbor()
			g.UpdateSouthernNeighbor()

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

func (g *Game) DrawWholeBoard() {
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
}

func (g *Game) DrawUpdatedSquares() {
	raylib.DrawRectangle(g.PrevSelected.X, g.PrevSelected.Y, square.SquareSize, square.SquareSize, g.PrevSelected.BgColor)
	raylib.DrawCircle(g.Selected.X+piece.PieceRadius, g.Selected.Y+piece.PieceRadius, float32(piece.PieceRadius), g.Selected.PieceColor)
	raylib.DrawEllipseLines(g.Selected.X+piece.PieceRadius, g.Selected.Y+piece.PieceRadius, float32(piece.PieceRadius), float32(piece.PieceRadius), g.Selected.BandColor)

	for g.ReDraw.Length() > 0 {
		tmp := g.ReDraw.Remove()
		if tmp != nil {
			raylib.DrawRectangle(tmp.X, tmp.Y, square.SquareSize, square.SquareSize, tmp.BgColor)
		}
	}

	if g.Selected.Piece == piece.King {
		raylib.DrawCircle(g.Selected.X+piece.PieceRadius, g.Selected.Y+piece.PieceRadius, piece.PieceRadius/2, raylib.Gold)
	}
}

func (g *Game) Draw() {
	raylib.BeginDrawing()
	if g.FirstTurn {
		g.DrawWholeBoard()
		g.FirstTurn = false
	}

	if g.MovePhase == 2 {
		raylib.DrawEllipseLines(g.Selected.X+piece.PieceRadius, g.Selected.Y+piece.PieceRadius, float32(piece.PieceRadius), float32(piece.PieceRadius), raylib.Green)
	}

	if g.MovePhase == 6 {
		g.DrawUpdatedSquares()
		g.MovePhase = 0
	}

	raylib.EndDrawing()
}
