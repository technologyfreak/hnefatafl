package piece

type PieceKind int8

const (
	None PieceKind = 1 << iota
	BlackPawn
	WhitePawn
	King
)
