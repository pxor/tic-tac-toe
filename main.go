package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"ttt_the_game/game"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	var score1 int
	var score2 int
	gameNum := 1 // Change to train statistics

	var w1 string
	var num int

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Choose game mode: 1-PVP 2-PVE 3-EVP 4-EVE 5-TrainDB(X) 6-TrainDB(O)")

	for {
		_, err := fmt.Scan(&w1)
		if err != nil {
			log.Fatal(err)
		}

		num, err = strconv.Atoi(w1)
		if err != nil {
			log.Fatal(err)
		}
		if game.PVP <= num && num <= game.EVE {
			if num == game.EVE {
				gameNum = 1000
			}
			break
		} else if num == game.TrainMode1 || num == game.TrainMode2 || num == game.TrainMode3 {
			gameNum = 10000
			break
		} else {
			fmt.Println("I said choose a game mode!")
		}
	}

	gameMode := game.GameMode(num)

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSL := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, dbSSL,
	)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	game.CreateDbTable(db)

	for range gameNum {
		if check := game.StartGame(db, gameMode); check == 1 {
			score1++
		} else if check == -1 {
			score2++
		}
	}

	fmt.Printf("\n------ Score %d:%d -------\n", score1, score2)
}
