package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"

	"github.com/joho/godotenv"
)

func main() {
	// load env vars
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}

	file, err := os.Open("response.json")
	if err != nil {
		log.Fatalf("Error opening response.json: %v", err)
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading response.json: %v", err)
	}

	if err = file.Close(); err != nil {
		log.Fatalf("Error closing file: %v", err)
	}

	var chall ChallengeResponse
	if err = json.Unmarshal(bytes, &chall); err != nil {
		log.Fatalf("Error unmarshalling json: %v", err)
	}

	var movies []*Movie
	// check if movies.json exists
	if _, err := os.Stat("movies.json"); err != nil {
		// if error is not exist error then get data and store in movies.json
		if os.IsNotExist(err) {
			movies = getAndStoreIMDBData(chall.Prompt.Movies)
		} else {
			log.Fatalf("Error checking if movies.json exists: %v", err)
		}
	} else {
		// read from cached movies data and unmarshal into movies slice
		file, err := os.Open("movies.json")
		if err != nil {
			log.Fatalf("Error opening movies.json: %v", err)
		}

		bytes, err := io.ReadAll(file)
		if err != nil {
			log.Fatalf("Error reading from movies.json: %v", err)
		}

		if err = json.Unmarshal(bytes, &movies); err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
		}
	}

	calcScore(movies, chall.Prompt.People)

	slices.SortFunc(movies, func(a, b *Movie) int { return int(b.Score - a.Score) })
	for _, movie := range movies {
		fmt.Printf("\"%s\",\n", movie.ID)
	}
}
