package game

type PlayerLogic interface {
	GetMove(t *TTTBoard) int
	StartTurn(t *TTTBoard) (int, FieldType)
	SingleTurn(t *TTTBoard) bool
	FieldType() FieldType
}

type Player struct {
	name      string
	fieldType FieldType
}

type HumanPlayer struct {
	Player
}

type CpuPlayer struct {
	Player
	method CpuMethod
}

func (p Player) FieldType() FieldType {
	return p.fieldType
}
