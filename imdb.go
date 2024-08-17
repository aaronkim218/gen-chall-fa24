package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func getAndStoreIMDBData(ids []string) []*Movie {
	var movies []*Movie
	for _, id := range ids {
		resp, err := http.Get(os.Getenv("OMDB_BASE") + "&i=" + id)
		if err != nil {
			log.Fatalf("Failed to get data from id: %s with error: %v", id, err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading body: %v", err)
		}
		if err = resp.Body.Close(); err != nil {
			log.Fatalf("Error closing body: %v", err)
		}

		var imdbResp IMDBResponse
		if err = json.Unmarshal(body, &imdbResp); err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
		}

		val, err := strconv.Atoi(imdbResp.Year)
		if err != nil {
			log.Fatalf("Error converting year to int: %v", err)
		} else if val < 0 {
			log.Fatalf("Year is negative")
		}
		year := uint(val)

		val, err = strconv.Atoi(strings.TrimSuffix(imdbResp.Runtime, " min"))
		if err != nil {
			log.Fatalf("Error converting runtime to string: %v", err)
		} else if val < 0 {
			log.Fatalf("Runtime is negative")
		}
		runtime := uint(val)

		actors := strings.Split(imdbResp.Actors, ", ")
		genres := strings.Split(imdbResp.Genre, ", ")

		var rottenTomatoes *uint
		for _, rating := range imdbResp.Ratings {
			if rating.Source == "Rotten Tomatoes" {
				val, err = strconv.Atoi(strings.TrimSuffix(rating.Value, "%"))
				if err != nil {
					log.Fatalf("Error converting value to integer: %v", err)
				} else if val < 0 {
					log.Fatalf("Rotten tomatoes score is negative")
				}

				uval := uint(val)
				rottenTomatoes = &uval
				break
			}
		}

		if rottenTomatoes == nil {
			log.Fatalf("Rotten tomatoes value not found")
		}

		movie := Movie{
			ID:             id,
			Title:          imdbResp.Title,
			Year:           year,
			Rated:          imdbResp.Rated,
			Runtime:        runtime,
			Genres:         genres,
			Director:       imdbResp.Director,
			Actors:         actors,
			Plot:           imdbResp.Plot,
			RottenTomatoes: *rottenTomatoes,
		}
		movies = append(movies, &movie)
	}

	data, err := json.MarshalIndent(movies, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	file, err := os.Create("movies.json")
	if err != nil {
		log.Fatalf("Error creating movies.json: %v", err)
	}

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("Error writing data to json file: %v", err)
	}

	if err = file.Close(); err != nil {
		log.Fatalf("Error closing file: %v", err)
	}

	return movies
}
