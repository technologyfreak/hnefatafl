package game

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	board "github.com/technologyfreak/hnefatafl/board"
	piece "github.com/technologyfreak/hnefatafl/piece"
	resources "github.com/technologyfreak/hnefatafl/resources"
	square "github.com/technologyfreak/hnefatafl/square"
)

const (
	screenWidth     = square.SquareSize * square.SquaresPerRow
	screenHeight    = screenWidth + square.SquareSize
	totalBlackPawns = 24
	totalWhitePawns = 12
	fontSize        = 25
	targetFPS       = 60
)

type NeigborKind uint8

const (
	Unopposed NeigborKind = iota
	Edge
	KingsSquare
	Opposed
)

const (
	BlacksTurnMsg   = "Black's Turn"
	WhitesTurnMsg   = "Whites's Turn"
	BlackWinsMsg    = "Black Wins!"
	WhiteWinsMsg    = "White Wins!"
	RestartBtnValue = "Click Here To Restart"
)

type CoordPair struct {
	X int32
	Y int32
}

type Game struct {
	ScreenWidth     int32
	ScreenHeight    int32
	BoardHeight     int32
	TurnMsgX        int32
	MsgY            int32
	WinMsgX         int32
	RestartBtnWidth int32
	RestartBtnX     int32
	RestartBtnY     int32

	BlackPawns uint8
	WhitePawns uint8
	MovePhase  uint8

	BlacksTurn              bool
	ShouldHighlightSelected bool
	Win                     bool

	Board board.Board

	Selected     *square.Square
	PrevSelected *square.Square

	BoardBackground raylib.Texture2D
	BlackPawnSprite raylib.Texture2D
	WhitePawnSprite raylib.Texture2D
	KingSprite      raylib.Texture2D
}

func (g *Game) Init() {
	g.ScreenWidth = screenWidth
	g.ScreenHeight = screenHeight
	g.BoardHeight = screenWidth

	g.Restart()

	raylib.InitWindow(int32(g.ScreenWidth), int32(g.ScreenHeight), "Hnefatafl")
	defer raylib.CloseWindow()

	raylib.ClearBackground(raylib.Beige)

	g.TurnMsgX = g.ScreenWidth/2 - raylib.MeasureText("XXXXX's Turn", fontSize)/2
	g.MsgY = g.ScreenHeight - fontSize

	g.WinMsgX = g.ScreenWidth/2 - raylib.MeasureText("XXXXX Wins!", fontSize)/2

	g.RestartBtnWidth = raylib.MeasureText(RestartBtnValue, fontSize)
	g.RestartBtnX = g.ScreenWidth/2 - g.RestartBtnWidth/2
	g.RestartBtnY = g.ScreenHeight/2 - g.RestartBtnWidth/2

	img := raylib.LoadImageFromMemory(".png", resources.BoardBackground, int32(len(resources.BoardBackground)))
	g.BoardBackground = raylib.LoadTextureFromImage(img)
	raylib.UnloadImage(img)
	defer raylib.UnloadTexture(g.BoardBackground)

	img = raylib.LoadImageFromMemory(".png", resources.BlackPawnSprite, int32(len(resources.BlackPawnSprite)))
	g.BlackPawnSprite = raylib.LoadTextureFromImage(img)
	raylib.UnloadImage(img)
	defer raylib.UnloadTexture(g.BlackPawnSprite)

	img = raylib.LoadImageFromMemory(".png", resources.WhitePawnSprite, int32(len(resources.WhitePawnSprite)))
	g.WhitePawnSprite = raylib.LoadTextureFromImage(img)
	raylib.UnloadImage(img)
	defer raylib.UnloadTexture(g.WhitePawnSprite)

	img = raylib.LoadImageFromMemory(".png", resources.KingSprite, int32(len(resources.KingSprite)))
	g.KingSprite = raylib.LoadTextureFromImage(img)
	raylib.UnloadImage(img)
	defer raylib.UnloadTexture(g.KingSprite)

	raylib.SetTargetFPS(targetFPS)

	for !raylib.WindowShouldClose() {
		g.Update()
		g.Draw()
	}

}

func (g *Game) SelectSquare() {
	if g.MovePhase == 0 {
		g.Selected = nil
		g.PrevSelected = nil
	} else {
		g.ShouldHighlightSelected = false
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
			((g.Selected.Piece&piece.BlackPawn == piece.BlackPawn && g.BlacksTurn) || (g.Selected.Piece&piece.WhitePawn == piece.WhitePawn && !g.BlacksTurn)) {
			g.MovePhase++
			g.ShouldHighlightSelected = true
		} else {
			g.MovePhase = 0
		}
	}

	if g.MovePhase == 3 {
		if !g.Selected.HasPiece() {
			if (g.Selected.IsKingsCorner() || g.Selected.IsCenter()) &&
				g.PrevSelected.Piece&piece.King != piece.King {
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
	return (g.Board.Squares[0][0].Piece&piece.King == piece.King) ||
		(g.Board.Squares[0][10].Piece&piece.King == piece.King) ||
		(g.Board.Squares[10][0].Piece&piece.King == piece.King) ||
		(g.Board.Squares[10][10].Piece&piece.King == piece.King)
}

func (g *Game) GetWesternNeighbor(wPiece *square.Square) (NeigborKind, *square.Square) {
	if wPiece == nil {
		return Edge, nil
	}

	west := CoordPair{X: square.ToRowOrCol(wPiece.X - square.SquareSize), Y: square.ToRowOrCol(wPiece.Y)}

	if square.InRowRange(west.X) && square.InRowRange(west.Y) {
		neigbor := &g.Board.Squares[west.X][west.Y]
		if neigbor.HasPiece() {
			if (neigbor.Piece & piece.BlackPawn) != (wPiece.Piece & piece.BlackPawn) {
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
			if (neigbor.Piece & piece.BlackPawn) != (wPiece.Piece & piece.BlackPawn) {
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
			if (neigbor.Piece & piece.BlackPawn) != (wPiece.Piece & piece.BlackPawn) {
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
			if (neigbor.Piece & piece.BlackPawn) != (wPiece.Piece & piece.BlackPawn) {
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

	if wPiece.Piece&piece.King == piece.King {
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
		if westSquare.Piece&piece.King == piece.King {
			g.Win = true
		} else if westSquare.Piece&piece.BlackPawn == piece.BlackPawn {
			g.BlackPawns--
		} else {
			g.WhitePawns--
		}

		westSquare.RemovePiece()
	}
}

func (g *Game) UpdateEasternNeighbor() {
	eastKind, eastSquare := g.GetEasternNeighbor(g.Selected)

	if eastKind > 0 && g.IsSandwiched(eastSquare) {
		if eastSquare.Piece&piece.King == piece.King {
			g.Win = true
		} else if eastSquare.Piece&piece.BlackPawn == piece.BlackPawn {
			g.BlackPawns--
		} else {
			g.WhitePawns--
		}

		eastSquare.RemovePiece()
	}
}

func (g *Game) UpdateNorthernNeighbor() {
	northKind, northSquare := g.GetNortherNeighbor(g.Selected)

	if northKind > 0 && g.IsSandwiched(northSquare) {
		if northSquare.Piece&piece.King == piece.King {
			g.Win = true
		} else if northSquare.Piece == piece.BlackPawn {
			g.BlackPawns--
		} else {
			g.WhitePawns--
		}

		northSquare.RemovePiece()
	}
}

func (g *Game) UpdateSouthernNeighbor() {
	southKind, southSquare := g.GetSouternNeighbor(g.Selected)

	if southKind > 0 && g.IsSandwiched(southSquare) {
		if southSquare.Piece&piece.King == piece.King {
			g.Win = true
		} else if southSquare.Piece&piece.BlackPawn == piece.BlackPawn {
			g.BlackPawns--
		} else {
			g.WhitePawns--
		}

		southSquare.RemovePiece()
	}
}

func (g *Game) Restart() {
	g.BlackPawns = totalBlackPawns
	g.WhitePawns = totalWhitePawns
	g.MovePhase = 0

	g.BlacksTurn = true
	g.ShouldHighlightSelected = false
	g.Win = false

	g.Board = board.NewBoard()
}

func (g *Game) Update() {
	if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
		if g.Win {
			x := raylib.GetMouseX()
			y := raylib.GetMouseY()

			if (x >= g.RestartBtnX && x < (g.RestartBtnX+g.RestartBtnWidth-1)) &&
				(y >= g.RestartBtnY && y < (g.RestartBtnY+fontSize-1)) {
				g.Restart()
			} else {
				return
			}
		}

		g.SelectSquare()
		g.ValidateMove()

		if g.MovePhase == 5 {
			// Swap piece to new location
			g.Selected.AddPiece(g.PrevSelected.Piece)
			g.PrevSelected.RemovePiece()

			g.UpdateWesternNeighbor()
			g.UpdateEasternNeighbor()
			g.UpdateNorthernNeighbor()
			g.UpdateSouthernNeighbor()

			g.BlacksTurn = !g.BlacksTurn // toggle turn order
			g.MovePhase = 0
		}
	}

	if g.KingHasReachedACorner() || g.BlackPawns == 0 || g.WhitePawns == 0 {
		g.Win = true
	}
}

func (g *Game) DrawPieces(wPiece *square.Square) {
	switch {
	case wPiece.Piece&piece.BlackPawn == piece.BlackPawn:
		raylib.DrawTexture(g.BlackPawnSprite, wPiece.X, wPiece.Y, raylib.RayWhite)
	case wPiece.Piece&piece.King == piece.King:
		raylib.DrawTexture(g.KingSprite, wPiece.X, wPiece.Y, raylib.RayWhite)
	case wPiece.Piece&piece.WhitePawn == piece.WhitePawn:
		raylib.DrawTexture(g.WhitePawnSprite, wPiece.X, wPiece.Y, raylib.RayWhite)
	}
}

func (g *Game) DrawBoard() {
	raylib.DrawTexture(g.BoardBackground, 0, 0, raylib.RayWhite)

	for i := 0; i < square.SquaresPerRow; i++ {
		for j := 0; j < square.SquaresPerRow; j++ {
			s := &g.Board.Squares[i][j]

			if s.HasPiece() {
				g.DrawPieces(s)
			}
		}
	}
}

func (g *Game) DrawTurnMsg() {
	raylib.DrawRectangle(g.TurnMsgX, g.MsgY, g.ScreenWidth, fontSize, raylib.Beige)

	turnMsg := BlacksTurnMsg
	turnColor := raylib.Black

	if !g.BlacksTurn {
		turnMsg = WhitesTurnMsg
		turnColor = raylib.White
	}

	raylib.DrawText(turnMsg, g.TurnMsgX+1, g.MsgY+1, fontSize, raylib.Gray)
	raylib.DrawText(turnMsg, g.TurnMsgX-1, g.MsgY-1, fontSize, raylib.Gray)
	raylib.DrawText(turnMsg, g.TurnMsgX, g.MsgY, fontSize, turnColor)
}

func (g *Game) DrawWinMsg() {
	raylib.DrawRectangle(g.TurnMsgX, g.MsgY, g.ScreenWidth, fontSize, raylib.Beige)

	winMsg := BlackWinsMsg
	winColor := raylib.Black

	if g.KingHasReachedACorner() || g.BlackPawns == 0 {
		winMsg = WhiteWinsMsg
		winColor = raylib.White
	}

	raylib.DrawText(winMsg, g.WinMsgX+1, g.MsgY+1, fontSize, raylib.Gray)
	raylib.DrawText(winMsg, g.WinMsgX-1, g.MsgY-1, fontSize, raylib.Gray)
	raylib.DrawText(winMsg, g.WinMsgX, g.MsgY, fontSize, winColor)
}

func (g *Game) DrawRestartBtn() {
	raylib.DrawRectangle(g.RestartBtnX, g.RestartBtnY, g.RestartBtnWidth, fontSize, raylib.DarkPurple)
	raylib.DrawText(RestartBtnValue, g.RestartBtnX+1, g.RestartBtnY+1, fontSize, raylib.Gray)
	raylib.DrawText(RestartBtnValue, g.RestartBtnX-1, g.RestartBtnY-1, fontSize, raylib.Gray)
	raylib.DrawText(RestartBtnValue, g.RestartBtnX, g.RestartBtnY, fontSize, raylib.Gold)
}

func (g *Game) Draw() {
	raylib.BeginDrawing()

	g.DrawBoard()
	if g.ShouldHighlightSelected {
		raylib.DrawRectangleLinesEx(raylib.NewRectangle(float32(g.Selected.X), float32(g.Selected.Y), square.SquareSize, square.SquareSize), 1.5, raylib.Green)
	}

	if g.Win {
		g.DrawWinMsg()
		g.DrawRestartBtn()
	} else {
		g.DrawTurnMsg()
	}
	raylib.EndDrawing()
}
