package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"ttt_the_game/game"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


type gameResult struct {
	score1 int
	score2 int
}

func runParallelGames(db *sql.DB, mode game.GameMode, gameNum int, workers int) (int, int) {
	jobs := make(chan struct{}, gameNum)
	results := make(chan gameResult, gameNum)

	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Go(func() {
			for range jobs {
				switch game.StartGame(db, mode) {
				case 1:
					results <- gameResult{score1: 1}
				case -1:
					results <- gameResult{score2: 1}
				default:
					results <- gameResult{}
				}
			}
		})
	}

	for i := 0; i < gameNum; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var score1, score2 int
	for r := range results {
		score1 += r.score1
		score2 += r.score2
	}

	return score1, score2
}

func main() {
	var score1 int
	var score2 int
	gameNum := 1

	var w1 string
	var num int

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Choose game mode: 1-PVP 2-PVE 3-EVP 4-EVE 5-TrainDB(X) 6-TrainDB(O) 7-TrainDB(Random)")

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
				gameNum = 10000
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

	db.SetMaxOpenConns(32)
	db.SetMaxIdleConns(32)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	game.CreateDbTable(db)

	switch gameMode {
	case game.PVP, game.PVE, game.EVP:
		for i := 0; i < gameNum; i++ {
			switch game.StartGame(db, gameMode) {
			case 1:
				score1++
			case -1:
				score2++
			}
		}
	default:
		score1, score2 = runParallelGames(db, gameMode, gameNum, game.WorkerCount)
	}

	fmt.Printf("\n------ Score %d:%d -------\n", score1, score2)
}
