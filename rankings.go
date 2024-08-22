package main

import (
	"cmp"
	"log"
	"math"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

func getCmpMoviesFunc(people []*Person) func(*Movie, *Movie) int {
	return func(movie1, movie2 *Movie) int {
		med1 := calcMedProp(movie1, people)
		med2 := calcMedProp(movie2, people)
		avg1 := calcAvgProp(movie1, people)
		avg2 := calcAvgProp(movie1, people)

		movie1.MedProp = med1
		movie2.MedProp = med2
		movie1.AvgProp = avg1
		movie2.AvgProp = avg2

		return -1 * cmp.Compare(med1, med2)
		// return -1 * cmp.Compare(avg1, avg2)
	}
}

// calculate average satisfaction proportion
func calcAvgProp(movie *Movie, people []*Person) float64 {
	total := float64(0)
	count := uint(0)

	for _, person := range people {
		total += calcProp(movie, person)
		count++
	}

	return total / float64(count)
}

// calculate median satisfaction proportion
func calcMedProp(movie *Movie, people []*Person) float64 {
	var props []float64
	for _, person := range people {
		props = append(props, calcProp(movie, person))
	}

	slices.SortFunc(props, func(a, b float64) int { return cmp.Compare(a, b) })

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
func calcProp(movie *Movie, person *Person) float64 {
	totalSat := float64(0)
	totalWeight := uint(0)

	t := reflect.TypeOf(person.Preferences).Elem()
	v := reflect.ValueOf(person.Preferences).Elem()

	for i := 0; i < v.NumField(); i++ {
		val := v.Field(i)

		if val.IsNil() {
			continue
		}

		sat, weight := calcSat(t.Field(i).Name, &val, movie)
		totalSat += sat
		totalWeight += weight
	}

	// possible range [-totalWeight, totalWeight]
	// in order to keep prop [0, 1] add totalWeight to numerator and divide by 2 * totalWeight
	return (float64(totalSat) + float64(totalWeight)) / (2 * float64(totalWeight))
}

// calculate satisfaction given preference name, preference, and movie
func calcSat(name string, val *reflect.Value, movie *Movie) (float64, uint) {
	pref := val.Interface()

	var sat float64
	var weight uint

	switch name {
	case "AfterYearInclusive":
		p := pref.(*Preference[uint])
		if isAfterYearInclusive(movie.Year, p.Value) {
			sat = float64(p.Weight)
		} else {
			sat = calcDiffPenalty(float64(movie.Year), float64(p.Value), float64(Bounds.MIN_YEAR), float64(p.Weight))
		}
		weight = p.Weight
	case "BeforeYearExclusive":
		p := pref.(*Preference[uint])
		if isBeforeYearExclusive(movie.Year, p.Value) {
			sat = float64(p.Weight)
		} else {
			sat = calcDiffPenalty(float64(movie.Year), float64(p.Value-1), float64(Bounds.MAX_YEAR), float64(p.Weight))
		}
		weight = p.Weight
	case "MaximumAgeRatingInclusive":
		p := pref.(*Preference[string])
		movieRating := RATINGS[movie.Rated]
		prefRating := RATINGS[p.Value]
		if isMaximumAgeRatingInclusive(movieRating, prefRating) {
			sat = float64(p.Weight)
		} else {
			sat = calcDiffPenalty(float64(movieRating), float64(prefRating), float64(Bounds.MAX_RATING), float64(p.Weight))
		}
		weight = p.Weight
	case "ShorterThanExclusive":
		p := pref.(*Preference[string])
		prefRuntime := calcRuntime(p.Value)
		if isShorterThanExclusive(float64(movie.Runtime), prefRuntime) {
			sat = float64(p.Weight)
		} else {
			sat = calcDiffPenalty(float64(movie.Runtime), prefRuntime-1, float64(Bounds.MIN_RUNTIME), float64(p.Weight))
		}
		weight = p.Weight
	case "FavoriteGenre":
		p := pref.(*Preference[string])
		if isFavoriteGenre(movie.Genres, p.Value) {
			sat = float64(p.Weight)
		} else {
			sat = -1 * float64(p.Weight)
		}
		weight = p.Weight
	case "LeastFavoriteDirector":
		// maybe just penalize if present and not reward if not present -- consider this method for others
		p := pref.(*Preference[string])
		if isLeastFavoriteDirector(movie.Director, p.Value) {
			sat = -1 * float64(p.Weight)
		} else {
			sat = float64(p.Weight)
		}
		weight = p.Weight
	case "FavoriteActors":
		p := pref.(*Preference[[]string])
		ratio := ratioFavoriteActors(movie.Actors, p.Value)
		if ratio > 0 {
			sat = ratio * float64(p.Weight)
		} else {
			sat = -1 * float64(p.Weight)
		}
		weight = p.Weight
	case "FavoritePlotElements":
		p := pref.(*Preference[[]string])
		ratio := ratioFavoritePlotElements(movie.Plot, p.Value)
		if ratio > 0 {
			sat = ratio * float64(p.Weight)
		} else {
			sat = -1 * float64(p.Weight)
		}
		weight = p.Weight
	case "MinimumRottenTomatoesScoreInclusive":
		p := pref.(*Preference[uint])
		if isMinimumRottenTomatoesScoreInclusive(movie.RottenTomatoes, p.Value) {
			sat = float64(p.Weight)
		} else {
			sat = calcDiffPenalty(float64(movie.RottenTomatoes), float64(p.Value), float64(Bounds.MIN_ROTTENTOMATOES), float64(p.Weight))
		}
		weight = p.Weight
	}

	return sat, weight
}

func isAfterYearInclusive(movieYear, prefYear uint) bool {
	return movieYear >= prefYear
}

func isBeforeYearExclusive(movieYear, prefYear uint) bool {
	return movieYear < prefYear
}

func isMaximumAgeRatingInclusive(movieRating, prefRating uint) bool {
	return movieRating <= prefRating
}

func isShorterThanExclusive(movieRuntime, prefRuntime float64) bool {
	return movieRuntime < prefRuntime
}

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

func isLeastFavoriteDirector(movieDirector, prefDirector string) bool {
	return movieDirector == prefDirector
}

func ratioFavoriteActors(movieActors, prefActors []string) float64 {
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

	return float64(total) / float64(count)
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

func ratioFavoritePlotElements(moviePlot string, prefPlotElems []string) float64 {
	total := uint(0)
	count := uint(0)

	for _, elem := range prefPlotElems {
		if strings.Contains(moviePlot, elem) {
			total++
		}
		count++
	}

	return float64(total) / float64(count)
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

func isMinimumRottenTomatoesScoreInclusive(movieScore, prefScore uint) bool {
	return movieScore >= prefScore
}

func calcDiffPenalty(actual, pref, bound, weight float64) float64 {
	return -1 * float64(weight) * ((math.Abs(float64(actual) - float64(pref))) / (math.Abs(float64(bound) - float64(pref))))
}

//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////////////////////
// V2

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
