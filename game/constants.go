package game

const MIN_NUMBER = 0
const MAX_NUMBER = 8
const MIN_PRINT = MIN_NUMBER + 1
const MAX_PRINT = MAX_NUMBER + 1

type CpuMethod int

const (
	RandomMove = iota
	StatisticalMove
)

type FieldType int

const (
	EmptyField = iota
	XField
	OField
)

type GameMode int

const (
	PVP = iota + 1
	PVE
	EVP
	EVE
	TrainMode1
	TrainMode2
	TrainMode3
)

func (f FieldType) String() string {
	switch f {
	case EmptyField:
		return "EmptyField"
	case XField:
		return "XField"
	case OField:
		return "OField"
	default:
		return ""
	}
}
