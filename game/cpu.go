package game

import (
	"log"
	"math/rand/v2"
)

func newCpuPlayer(n string, ft FieldType, m CpuMethod) CpuPlayer {
	return CpuPlayer{Player: Player{name: n, fieldType: ft}, method: m}
}

func (c CpuPlayer) StartTurn(t *TTTBoard) (int, FieldType) {
	var mov int

	for {
		mov = c.GetMove(t)
		if t.isLegalMove(mov) {
			break
		} else {
			continue
		}
	}

	return mov, c.fieldType
}
func (c CpuPlayer) SingleTurn(t *TTTBoard) bool {
	t.displayBoard()
	i, ft := c.StartTurn(t)
	t.updateField(i, ft)
	isWon := t.checkWinCond(c.fieldType)
	if t.turnCounter > MAX_NUMBER {
		t.winner = EmptyField
		t.UpdateWinner()
		return true
	}
	if isWon {
		t.winner = c.fieldType
		t.displayBoard()
		t.UpdateWinner()
		return true
	}
	return false
}

func (c CpuPlayer) GetMove(t *TTTBoard) int {
	switch c.method {
	case RandomMove:
		return rand.IntN(MAX_NUMBER + 1)
	case StatisticalMove:
		num, err := t.bestMove()
		if err != nil {
			log.Fatal(err)
		}
		if num != -1 {
			return num
		}
	}
	return 0 // Should never be reached
}
