package main

import (
	"cmp"
	"log"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func getCmpMoviesFunc(people []*Person) func(*Movie, *Movie) int {
	return func(movie1 *Movie, movie2 *Movie) int {
		med1 := calcMedProp(movie1, people)
		med2 := calcMedProp(movie2, people)

		return -1 * cmp.Compare(med1, med2)
	}
}

// calculate median satisfaction proportion
func calcMedProp(movie *Movie, people []*Person) float32 {
	var props []float32
	for _, person := range people {
		props = append(props, calcProp(movie, person))
	}

	slices.SortFunc(props, func(a, b float32) int { return cmp.Compare(a, b) })

	length := len(props)

	if length == 0 {
		return 0
	} else if length%2 == 1 {
		return props[length/2]
	} else {
		return (props[(length/2)-1] + props[length/2]) / 2
	}
}

// calculate satisfaction proportion
func calcProp(movie *Movie, person *Person) float32 {
	totalSat := float32(0)
	totalWeight := uint(0)

	t := reflect.TypeOf(person.Preferences).Elem()
	v := reflect.ValueOf(person.Preferences).Elem()

	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i)

		if val.IsNil() {
			continue
		}

		pref := val.Interface()

		var sat float32
		var weight uint

		switch t.Field(i).Name {
		case "AfterYearInclusive":
			p := pref.(*Preference[uint])
			if isAfterYearInclusive(movie.Year, p.Value) {
				sat = float32(p.Weight)
			} else {
				sat = -1 * float32(p.Weight)
			}
			weight = p.Weight
		case "BeforeYearExclusive":
			p := pref.(*Preference[uint])
			if isBeforeYearExclusive(movie.Year, p.Value) {
				sat = float32(p.Weight)
			} else {
				sat = -1 * float32(p.Weight)
			}
			weight = p.Weight
		case "MaximumAgeRatingInclusive":
			p := pref.(*Preference[string])
			if isMaximumAgeRatingInclusive(movie.Rated, p.Value) {
				sat = float32(p.Weight)
			} else {
				sat = -1 * float32(p.Weight)
			}
			weight = p.Weight
		case "ShorterThanExclusive":
			p := pref.(*Preference[string])
			if isShorterThanExclusive(movie.Runtime, p.Value) {
				sat = float32(p.Weight)
			} else {
				sat = -1 * float32(p.Weight)
			}
			weight = p.Weight
		case "FavoriteGenre":
			p := pref.(*Preference[string])
			if isFavoriteGenre(movie.Genres, p.Value) {
				sat = float32(p.Weight)
			} else {
				sat = -1 * float32(p.Weight)
			}
			weight = p.Weight
		case "LeastFavoriteDirector":
			// maybe just penalize if present and not reward if not present -- consider this method for others
			p := pref.(*Preference[string])
			if isLeastFavoriteDirector(movie.Director, p.Value) {
				sat = -1 * float32(p.Weight)
			} else {
				sat = float32(p.Weight)
			}
			weight = p.Weight
		case "FavoriteActors":
			p := pref.(*Preference[[]string])
			ratio := ratioFavoriteActors(movie.Actors, p.Value)
			if ratio > 0 {
				sat = ratio * float32(p.Weight)
			} else {
				sat = -1 * float32(p.Weight)
			}
			weight = p.Weight
		case "FavoritePlotElements":
			p := pref.(*Preference[[]string])
			ratio := ratioFavoritePlotElements(movie.Plot, p.Value)
			if ratio > 0 {
				sat = ratio * float32(p.Weight)
			} else {
				sat = -1 * float32(p.Weight)
			}
			weight = p.Weight
		case "MinimumRottenTomatoesScoreInclusive":
			p := pref.(*Preference[uint])
			if isMinimumRottenTomatoesScoreInclusive(movie.RottenTomatoes, p.Value) {
				sat = float32(p.Weight)
			} else {
				sat = -1 * float32(p.Weight)
			}
			weight = p.Weight
		}

		totalSat += sat
		totalWeight += weight
	}

	return (float32(totalSat) + float32(totalWeight)) / (2 * float32(totalWeight))
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
