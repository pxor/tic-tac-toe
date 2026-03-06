package game

import "database/sql"

func StartGame(db *sql.DB, mode GameMode) int {
	tBoard := newBoard(db)
	var p1, p2 PlayerLogic
	switch mode {
	case PVP:
		p1 = newHumanPlayer("Player_1", XField)
		p2 = newHumanPlayer("Player_2", OField)
	case PVE:
		p1 = newHumanPlayer("Player_1", XField)
		p2 = newCpuPlayer("CPU_1", OField, StatisticalMove)
	case EVP:
		p1 = newCpuPlayer("CPU_1", XField, StatisticalMove)
		p2 = newHumanPlayer("Player_1", OField)
	case EVE:
		p1 = newCpuPlayer("CPU_1", XField, StatisticalMove)
		p2 = newCpuPlayer("CPU_2", OField, StatisticalMove)
	case TrainMode1:
		p1 = newCpuPlayer("CPU_1", XField, StatisticalMove)
		p2 = newCpuPlayer("CPU-2", OField, RandomMove)
	case TrainMode2:
		p1 = newCpuPlayer("CPU-2", XField, RandomMove)
		p2 = newCpuPlayer("CPU_1", OField, StatisticalMove)
	case TrainMode3:
		p1 = newCpuPlayer("CPU-2", XField, RandomMove)
		p2 = newCpuPlayer("CPU-2", OField, RandomMove)
	}

	for {
		tBoard.incTurnCounter()
		tBoard.changeActivePlayer(p1.FieldType())
		if p1.SingleTurn(&tBoard) {
			break
		}
		tBoard.incTurnCounter()
		tBoard.changeActivePlayer(p2.FieldType())
		if p2.SingleTurn(&tBoard) {
			break
		}
	}

	switch tBoard.getWinner() {
	case XField:
		return 1
	case OField:
		return -1
	}

	return 0
}
