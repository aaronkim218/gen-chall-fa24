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

	slices.SortFunc(movies, getCmpMoviesFunc(chall.Prompt.People))

	for _, movie := range movies {
		fmt.Printf("\"%s\" --- med: %v, avg: %v\n", movie.ID, movie.MedProp, movie.AvgProp)
	}
}

// see how much each movie satisfies each person
// satisfying 1st pref is better than satisfying 2nd pref
// p1 - m1 satisfies 1st 2nd 4th pref
// p2 - m1 satisfies 4th pref
// p1 - m2 satisfies 2nd 3rd 4th pref
// p2 - m2 staisfies 4th pref
// movie 1 has higher satisfaction because satisfying 1st pref vs 3rd pref for p1 is better. everything else same
// but how to compare people that have different number of prefs? e.g. 4 pref vs 7

// nvm just go by weight
// comparator function should compare total weights of satisfied preferences
// could also calculate total weight of preference satisfaction per person and divide by total weight so i get proportion of satisfaction
// comparator function should compare average satisfaction
// get proportion of satisfaction for each person and then calculate average
// movie with higher average satisfaction is better

// TRY THIS ONE
// "tt0058150",
// "tt0432283",
// "tt22022452",
// "tt0112384",
// "tt0264464",
// "tt1074638",
// "tt1285016",
// "tt3783958",
// "tt2582802",
// "tt2084970",
