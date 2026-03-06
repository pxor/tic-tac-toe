package game

import (
	"fmt"
	"log"
	"strconv"
)

func newHumanPlayer(n string, ft FieldType) HumanPlayer {
	return HumanPlayer{Player: Player{name: n, fieldType: ft}}
}

func (hp HumanPlayer) GetMove(t *TTTBoard) int {
	var w1 string
	var num int

	for {
		_, err := fmt.Scan(&w1)
		if err != nil {
			log.Fatal(err)
		}

		num, err = strconv.Atoi(w1)
		if err != nil {
			log.Fatal(err)
		}
		if MIN_PRINT <= num && num <= MAX_PRINT {
			break
		} else {
			fmt.Printf("YOU NEED TO USE ONLY %d-%d\n", MIN_PRINT, MAX_PRINT)
		}
	}

	return num - 1
}

func (hp HumanPlayer) StartTurn(t *TTTBoard) (int, FieldType) {
	var mov int
	fmt.Printf("%s input move (%d-%d): ", hp.name, MIN_PRINT, MAX_PRINT)
	for {
		mov = hp.GetMove(t)
		if t.isLegalMove(mov) {
			break
		} else {
			fmt.Printf("This is an invalid move, %s input again (%d-%d): ", hp.name, MIN_PRINT, MAX_PRINT)
		}
	}

	return mov, hp.fieldType
}

func (hp HumanPlayer) SingleTurn(t *TTTBoard) bool {
	t.displayBoard()
	i, ft := hp.StartTurn(t)
	t.updateField(i, ft)
	isWon := t.checkWinCond(hp.fieldType)
	if isWon {
		t.winner = hp.fieldType
		t.displayBoard()
		t.UpdateWinner()
		return true
	}
	return false
}
