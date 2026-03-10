package game

import (
	"database/sql"
	"log"
)

func CreateDbTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS game (
		id SERIAL PRIMARY KEY,
		pid NUMERIC(6,0) NOT NULL,
		gameState VARCHAR(100) NOT NULL,
		turn NUMERIC(2,0) NOT NULL,
		move INTEGER NOT NULL,
		result VARCHAR(20) NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func newBoard(db *sql.DB) TTTBoard {
	t := TTTBoard{turnCounter: 0}
	t.initFields(db)
	return t
}

func insertMove(t *TTTBoard, i int) error {
	query := `INSERT INTO game (pid, gamestate, turn, move, result)
			  VALUES ($1, $2, $3, $4, $5)`

	_, err := t.db.Exec(query, t.gamePid, t.fieldsAsString(), t.turnCounter, i, EmptyField)
	return err
}

func updateDbField(t *TTTBoard) {
	query := `UPDATE game
          SET result = $1
          WHERE pid = $2`

	_, err := t.db.Exec(query, t.winner.String(), t.gamePid)
	if err != nil {
		log.Fatal(err)
	}
}

func calcMove(t *TTTBoard) (int, error) {
	query := `
		WITH filtered AS (
			SELECT
				move,
				result
			FROM game
			WHERE gamestate = $1
		),
		move_stats AS (
			SELECT
				move,
				COUNT(*) AS times_played,
				SUM(CASE WHEN result = $2 THEN 1 ELSE 0 END) AS wins,
				SUM(CASE WHEN result = 'EmptyField' THEN 1 ELSE 0 END) AS draws,
				SUM(CASE WHEN result <> $2 AND result <> 'EmptyField' THEN 1 ELSE 0 END) AS losses
			FROM filtered
			GROUP BY move
		)
		SELECT move
		FROM move_stats
		ORDER BY
			(
				wins::numeric * 1.0 +
				draws::numeric * 0.5 -
				losses::numeric * 1.0
			) DESC,
			(wins::numeric / NULLIF(times_played, 0)) DESC,
			(draws::numeric / NULLIF(times_played, 0)) DESC,
			(losses::numeric / NULLIF(times_played, 0)) ASC,
			times_played DESC,
			move ASC;
	`

	rows, err := t.db.Query(query, t.fieldsAsString(), t.activePlayer.String())
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	for rows.Next() {
		var move int
		if err := rows.Scan(&move); err != nil {
			return -1, err
		}

		if move >= 0 && move <= 8 && t.fields[move] == EmptyField {
			return move, nil
		}
	}

	if err := rows.Err(); err != nil {
		return -1, err
	}
	return -1, nil
}
