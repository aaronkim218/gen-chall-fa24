package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// calculate runtime in minutes
func calcRuntime(runtime string) float64 {
	re := regexp.MustCompile(`(\d+)h(\d+)m(\d+)s`)
	matches := re.FindStringSubmatch(runtime)

	hrs, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Fatalf("Error parsing hours: %v", err)
	}

	mins, err := strconv.Atoi(matches[2])
	if err != nil {
		log.Fatalf("Error parsing minutes: %v", err)
	}

	secs, err := strconv.Atoi(matches[3])
	if err != nil {
		log.Fatalf("Error parsing seconds: %v", err)
	}

	return (60 * float64(hrs)) + float64(mins) + (float64(secs) / 60)
}

func isFavoriteGenre(movieGenres []string, prefGenre string) bool {
	for _, genre := range movieGenres {
		if genre == prefGenre {
			return true
		}
	}

	return false
}

func matchesFavoriteActors(movieActors, prefActors []string) uint {
	movieActorsMap := make(map[string]struct{})

	for _, actor := range movieActors {
		movieActorsMap[actor] = struct{}{}
	}

	total := uint(0)

	for _, actor := range prefActors {
		if _, ok := movieActorsMap[actor]; ok {
			total++
		}
	}

	return total
}

func matchesFavoritePlotElements(moviePlot string, prefPlotElems []string) uint {
	total := uint(0)

	for _, elem := range prefPlotElems {
		if strings.Contains(moviePlot, elem) {
			total++
		}
	}

	return total
}

func calcPoints(movies []*Movie, people []*Person) {
	var wg sync.WaitGroup

	for _, movie := range movies {
		for _, person := range people {
			wg.Add(1)

			go func(m *Movie, p *Person) {
				defer wg.Done()
				atomic.AddInt64(&m.PointsV3, int64(score(m, p)))
			}(movie, person)
		}
	}

	wg.Wait()
}

func score(movie *Movie, person *Person) int {
	total := 0

	if person.Preferences.AfterYearInclusive != nil {
		if movie.Year >= person.Preferences.AfterYearInclusive.Value {
			total += int(person.Preferences.AfterYearInclusive.Weight)
		}
	}

	if person.Preferences.BeforeYearExclusive != nil {
		if movie.Year < person.Preferences.BeforeYearExclusive.Value {
			total += int(person.Preferences.BeforeYearExclusive.Weight)
		}
	}

	if person.Preferences.MaximumAgeRatingInclusive != nil {
		if RATINGS[movie.Rated] <= RATINGS[person.Preferences.MaximumAgeRatingInclusive.Value] {
			total += int(person.Preferences.MaximumAgeRatingInclusive.Weight)
		}
	}

	if person.Preferences.ShorterThanExclusive != nil {
		if movie.Runtime < uint(calcRuntime(person.Preferences.ShorterThanExclusive.Value)) {
			total += int(person.Preferences.ShorterThanExclusive.Weight)
		}
	}

	if person.Preferences.FavoriteGenre != nil {
		if isFavoriteGenre(movie.Genres, person.Preferences.FavoriteGenre.Value) {
			total += int(person.Preferences.FavoriteGenre.Weight)
		}
	}

	if person.Preferences.LeastFavoriteDirector != nil {
		if movie.Director == person.Preferences.LeastFavoriteDirector.Value {
			total -= int(person.Preferences.LeastFavoriteDirector.Weight)
		}
	}

	if person.Preferences.FavoriteActors != nil {
		matches := matchesFavoriteActors(movie.Actors, person.Preferences.FavoriteActors.Value)
		total += int(matches * person.Preferences.FavoriteActors.Weight)
	}

	if person.Preferences.FavoritePlotElements != nil {
		matches := matchesFavoritePlotElements(movie.Plot, person.Preferences.FavoritePlotElements.Value)
		total += int(matches * person.Preferences.FavoritePlotElements.Weight)
	}

	if person.Preferences.MinimumRottenTomatoesScoreInclusive != nil {
		if movie.RottenTomatoes >= person.Preferences.MinimumRottenTomatoesScoreInclusive.Value {
			total += int(person.Preferences.MinimumRottenTomatoesScoreInclusive.Weight)
		}
	}

	return total
}
