package piece

type PieceKind int8

const (
	PieceRadius = 16
)

const (
	None PieceKind = iota
	Pawn
	King
)
