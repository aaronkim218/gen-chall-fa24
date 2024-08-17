package main

import (
	"cmp"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func getCompMoviesFunc(people []*Person) func(*Movie, *Movie) int {
	return func(movie1 *Movie, movie2 *Movie) int {
		prop1 := calcAvgProp(movie1, people)
		prop2 := calcAvgProp(movie2, people)

		return cmp.Compare(prop1, prop2)
	}
}

// calculate average satisfaction proportion
func calcAvgProp(movie *Movie, people []*Person) float32 {
	total := float32(0)
	count := uint(0)

	for _, person := range people {
		total += calcProp(movie, person)
		count++
	}

	return total / float32(count)
}

// calculate satisfaction proportion
func calcProp(movie *Movie, person *Person) float32 {
	total := float32(0)
	count := uint(0)

	t := reflect.TypeOf(person.Preferences).Elem()
	v := reflect.ValueOf(person.Preferences).Elem()

	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i)

		if val.IsNil() {
			continue
		}

		pref := val.Interface()

		sat := float32(0)

		switch t.Field(i).Name {
		case "AfterYearInclusive":
			p := pref.(*Preference[uint])
			if isAfterYearInclusive(movie.Year, p.Value) {
				sat = float32(p.Weight)
			}
		case "BeforeYearExclusive":
			p := pref.(*Preference[uint])
			if isBeforeYearExclusive(movie.Year, p.Value) {
				sat = float32(p.Weight)
			}
		case "MaximumAgeRatingInclusive":
			p := pref.(*Preference[string])
			if isMaximumAgeRatingInclusive(movie.Rated, p.Value) {
				sat = float32(p.Weight)
			}
		case "ShorterThanExclusive":
			p := pref.(*Preference[string])
			if isShorterThanExclusive(movie.Runtime, p.Value) {
				sat = float32(p.Weight)
			}
		case "FavoriteGenre":
			p := pref.(*Preference[string])
			if isFavoriteGenre(movie.Genres, p.Value) {
				sat = float32(p.Weight)
			}
		case "LeastFavoriteDirector":
			p := pref.(*Preference[string])
			if !isLeastFavoriteDirector(movie.Director, p.Value) {
				sat = float32(p.Weight)
			}
		case "FavoriteActors":
			p := pref.(*Preference[[]string])
			ratio := ratioFavoriteActors(movie.Actors, p.Value)
			sat = ratio * float32(p.Weight)
		case "FavoritePlotElements":
			p := pref.(*Preference[[]string])
			ratio := ratioFavoritePlotElements(movie.Plot, p.Value)
			sat = ratio * float32(p.Weight)
		case "MinimumRottenTomatoesScoreInclusive":
			p := pref.(*Preference[uint])
			if isMinimumRottenTomatoesScoreInclusive(movie.RottenTomatoes, p.Value) {
				sat = float32(p.Weight)
			}
		}

		total += sat
		count++
	}

	return float32(total) / float32(count)
}

func isAfterYearInclusive(movieYear, prefYear uint) bool {
	return movieYear >= prefYear
}

func isBeforeYearExclusive(movieYear, prefYear uint) bool {
	return movieYear < prefYear
}

func isMaximumAgeRatingInclusive(movieRating, prefRating string) bool {
	var ratings = map[string]int{
		"G":     0,
		"PG":    1,
		"PG-13": 2,
		"R":     3,
		"NC-17": 4,
	}

	return ratings[movieRating] <= ratings[prefRating]
}

func isShorterThanExclusive(movieRuntime uint, prefRuntime string) bool {
	return float32(movieRuntime) < calcRuntime(prefRuntime)
}

// calculate runtime in minutes
func calcRuntime(runtime string) float32 {
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

	return (60 * float32(hrs)) + float32(mins) + (float32(secs) / 60)
}

func isFavoriteGenre(movieGenres []string, prefGenre string) bool {
	for _, genre := range movieGenres {
		if genre == prefGenre {
			return true
		}
	}

	return false
}

func isLeastFavoriteDirector(movieDirector, prefDirector string) bool {
	return movieDirector == prefDirector
}

func ratioFavoriteActors(movieActors, prefActors []string) float32 {
	movieActorsMap := make(map[string]struct{})

	for _, actor := range movieActors {
		movieActorsMap[actor] = struct{}{}
	}

	total := uint(0)
	count := uint(0)

	for _, actor := range prefActors {
		if _, ok := movieActorsMap[actor]; ok {
			total++
		}
		count++
	}

	return float32(total) / float32(count)
}

func ratioFavoritePlotElements(moviePlot string, prefPlotElems []string) float32 {
	total := uint(0)
	count := uint(0)

	for _, elem := range prefPlotElems {
		if strings.Contains(moviePlot, elem) {
			total++
		}
		count++
	}

	return float32(total) / float32(count)
}

func isMinimumRottenTomatoesScoreInclusive(movieScore, prefScore uint) bool {
	return movieScore >= prefScore
}
