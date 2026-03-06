package game

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"strings"
)

type TTTBoard struct {
	fields       [MAX_NUMBER + 1]FieldType
	turnCounter  int
	gamePid      int
	activePlayer FieldType
	winner       FieldType
	db           *sql.DB
}

func (t *TTTBoard) incTurnCounter() {
	t.turnCounter++
}

func (t *TTTBoard) changeActivePlayer(ft FieldType) {
	t.activePlayer = ft
}

func (t *TTTBoard) getWinner() FieldType {
	return t.winner
}

func (t *TTTBoard) initFields(db *sql.DB) {
	for i := range MAX_NUMBER + 1 {
		t.fields[i] = EmptyField
	}
	t.db = db
	t.gamePid = t.genPID()
}

func (t *TTTBoard) displayBoard() {
	for i := range MAX_NUMBER + 1 {
		if i%3 == 0 {
			fmt.Print("\n")
		}
		switch t.fields[i] {
		case EmptyField:
			fmt.Print(".")
		case XField:
			fmt.Print("X")
		case OField:
			fmt.Print("O")
		}
	}

	fmt.Print("\n")
}

func (t *TTTBoard) genPID() int {
	return rand.IntN(100000)
}

func (t *TTTBoard) isLegalMove(i int) bool {
	return t.fields[i] == EmptyField
}

func (t *TTTBoard) updateField(i int, ft FieldType) {
	t.fields[i] = ft
	err := insertMove(t, i)

	if err != nil {
		log.Fatal(err)
	}
}

func (t *TTTBoard) UpdateWinner() {
	updateDbField(t)
}

func (t *TTTBoard) checkWinCond(activeField FieldType) bool {
	// rows
	var firstField FieldType
	isWon := false
	for i := range 3 {
		for y := range 3 {
			if t.fields[i*3+y] != activeField {
				isWon = false
				break
			}
			if (i*3+y)%3 == 0 {
				firstField = t.fields[i*3+y]
				isWon = true
			}

			if t.fields[i*3+y] != firstField {
				isWon = false
				break
			}

		}
		if isWon {
			return true
		}
	}

	// col
	firstField = EmptyField
	isWon = false
	for y := range 3 {
		for i := range 3 {
			if t.fields[i*3+y] != activeField {
				isWon = false
				break
			}
			if (i*3+y)%3 == y {
				firstField = t.fields[i*3+y]
				isWon = true
			}

			if t.fields[i*3+y] != firstField {
				isWon = false
				break
			}

		}
		if isWon {
			return true
		}
	}

	// diagonal
	firstField = EmptyField
	isWon = false
	for i := 0; i <= MAX_NUMBER; i += 4 {
		if t.fields[i] != activeField {
			isWon = false
			break
		} else {
			firstField = t.fields[i]
		}

		if t.fields[i] == firstField {
			isWon = true
			continue
		}
	}

	if isWon {
		return true
	}

	firstField = EmptyField
	isWon = false
	for i := 2; i < MAX_NUMBER; i += 2 {
		if t.fields[i] != activeField {
			isWon = false
			break
		} else {
			firstField = t.fields[i]
		}

		if t.fields[i] == firstField {
			isWon = true
			continue
		}
	}

	return isWon || t.turnCounter == MAX_NUMBER+1
}

func (t *TTTBoard) fieldsAsString() string {
	var s strings.Builder
	for i := range MAX_NUMBER + 1 {
		s.WriteString(t.fields[i].String())
		if i != MAX_NUMBER {
			s.WriteString(",")
		}
	}

	return s.String()
}

func (t *TTTBoard) randomEmptyMove() int {
	var available []int

	for i, f := range t.fields {
		if f == EmptyField {
			available = append(available, i)
		}
	}

	if len(available) == 0 {
		return -1
	}

	return available[rand.IntN(len(available))]
}

func (t *TTTBoard) bestMove() (int, error) {
	n, err := calcMove(t)

	if err != nil {
		log.Fatal(err)
	} else if n != -1 {
		return n, nil
	}

	return t.randomEmptyMove(), nil
}
