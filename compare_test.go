package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAfterYearInclusive(t *testing.T) {
	tests := []struct {
		movieYear uint
		prefYear  uint
		expected  bool
	}{
		{1999, 2000, false},
		{2000, 2000, true},
		{2001, 2000, true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, isAfterYearInclusive(test.movieYear, test.prefYear))
	}
}

func TestIsBeforeYearExclusive(t *testing.T) {
	tests := []struct {
		movieYear uint
		prefYear  uint
		expected  bool
	}{
		{1999, 2000, true},
		{2000, 2000, false},
		{2001, 2000, false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, isBeforeYearExclusive(test.movieYear, test.prefYear))
	}
}

func TestIsMaximumAgeRatingInclusive(t *testing.T) {
	tests := []struct {
		prefRating  string
		movieRating string
		expected    bool
	}{
		{"G", "G", true},
		{"G", "PG", false},
		{"G", "PG-13", false},
		{"G", "R", false},
		{"G", "NC-17", false},
		{"PG", "G", true},
		{"PG", "PG", true},
		{"PG", "PG-13", false},
		{"PG", "R", false},
		{"PG", "NC-17", false},
		{"PG-13", "G", true},
		{"PG-13", "PG", true},
		{"PG-13", "PG-13", true},
		{"PG-13", "R", false},
		{"PG-13", "NC-17", false},
		{"R", "G", true},
		{"R", "PG", true},
		{"R", "PG-13", true},
		{"R", "R", true},
		{"R", "NC-17", false},
		{"NC-17", "G", true},
		{"NC-17", "PG", true},
		{"NC-17", "PG-13", true},
		{"NC-17", "R", true},
		{"NC-17", "NC-17", true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, isMaximumAgeRatingInclusive(test.movieRating, test.prefRating))
	}
}

func TestIsShorterThanExclusive(t *testing.T) {
	tests := []struct {
		movieRuntime uint
		prefRuntime  string
		expected     bool
	}{
		{131, "2h12m0s", true},
		{132, "2h12m0s", false},
		{133, "2h12m0s", false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, isShorterThanExclusive(test.movieRuntime, test.prefRuntime))
	}
}

func TestCalcRuntime(t *testing.T) {
	tests := []struct {
		runtime  string
		expected float32
	}{
		{"0h0m0s", 0},
		{"0h0m1s", float32(1) / 60},
		{"0h1m0s", 1},
		{"0h1m1s", 1 + (float32(1) / 60)},
		{"1h0m0s", 60},
		{"1h0m1s", 60 + (float32(1) / 60)},
		{"1h1m1s", 60 + 1 + (float32(1) / 60)},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, calcRuntime(test.runtime))
	}
}

func TestIsFavoriteGenre(t *testing.T) {
	tests := []struct {
		prefGenre   string
		movieGenres []string
		expected    bool
	}{
		{"action", []string{"comedy", "drama"}, false},
		{"action", []string{"action", "drama"}, true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, isFavoriteGenre(test.movieGenres, test.prefGenre))
	}
}

func TestIsLeastFavoriteDirector(t *testing.T) {
	tests := []struct {
		prefDirector  string
		movieDirector string
		expected      bool
	}{
		{"Christopher Nolan", "Ridley Scott", false},
		{"Christopher Nolan", "Christopher Nolan", true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, isLeastFavoriteDirector(test.movieDirector, test.prefDirector))
	}
}

func TestIsMinimumRottenTomatoesScoreInclusive(t *testing.T) {
	tests := []struct {
		movieScore uint
		prefScore  uint
		expected   bool
	}{
		{86, 87, false},
		{87, 87, true},
		{88, 87, true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, isMinimumRottenTomatoesScoreInclusive(test.movieScore, test.prefScore))
	}
}

func TestRatioFavoriteActors(t *testing.T) {
	tests := []struct {
		movieActors []string
		prefActors  []string
		expected    float32
	}{
		{[]string{"Ryan Gosling", "Jamie Foxx", "Chris Evans"}, []string{"Ryan Reynolds", "Josh Brolin", "Robert Pattinson"}, 0},
		{[]string{"Ryan Reynolds", "Jamie Foxx", "Chris Evans"}, []string{"Ryan Reynolds", "Josh Brolin", "Robert Pattinson"}, float32(1) / 3},
		{[]string{"Ryan Reynolds", "Josh Brolin", "Chris Evans"}, []string{"Ryan Reynolds", "Josh Brolin", "Robert Pattinson"}, float32(2) / 3},
		{[]string{"Ryan Reynolds", "Josh Brolin", "Robert Pattinson"}, []string{"Ryan Reynolds", "Josh Brolin", "Robert Pattinson"}, 1},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, ratioFavoriteActors(test.movieActors, test.prefActors))
	}
}

func TestRatioFavoritePlotElements(t *testing.T) {
	tests := []struct {
		moviePlot     string
		prefPlotElems []string
		expected      float32
	}{
		{"this is a movie about a serial killer", []string{"family", "war", "love"}, 0},
		{"this is a movie about war", []string{"family", "war", "love"}, float32(1) / 3},
		{"this is a movie about war and love", []string{"family", "war", "love"}, float32(2) / 3},
		{"this is a movie about family, war, and love", []string{"family", "war", "love"}, 1},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, ratioFavoritePlotElements(test.moviePlot, test.prefPlotElems))
	}
}

func TestCalcSat(t *testing.T) {
	goodMovie := Movie{
		ID:             "good",
		Title:          "Good Movie",
		Year:           2005,
		Rated:          "PG-13",
		Runtime:        100,
		Genres:         []string{"Action", "Comedy", "Drama"},
		Director:       "Aaron Kim",
		Actors:         []string{"Robert Pattinson", "Josh Brolin", "Chris Evans"},
		Plot:           "this is a movie about family, war, and love",
		RottenTomatoes: 85,
	}

	badMovie := Movie{
		ID:             "bad",
		Title:          "Bad Movie",
		Year:           1990,
		Rated:          "R",
		Runtime:        120,
		Genres:         []string{"Thriller", "Horror"},
		Director:       "Ridley Scott",
		Actors:         []string{"Kevin Durant", "Jon Jones", "Dana White"},
		Plot:           "this is a documentary about history",
		RottenTomatoes: 60,
	}

	badMovieBeforeYearExclusive := Movie{Year: 2011}

	partialMovie := Movie{
		Actors: []string{"Chris Evans"},
		Plot:   "this is a movie about love",
	}

	prefs := Preferences{
		AfterYearInclusive:                  &Preference[uint]{Value: 2000, Weight: 10},
		BeforeYearExclusive:                 &Preference[uint]{Value: 2010, Weight: 10},
		MaximumAgeRatingInclusive:           &Preference[string]{Value: "PG-13", Weight: 10},
		ShorterThanExclusive:                &Preference[string]{Value: "1h45m0s", Weight: 10},
		FavoriteGenre:                       &Preference[string]{Value: "Action", Weight: 10},
		LeastFavoriteDirector:               &Preference[string]{Value: "Ridley Scott", Weight: 10},
		FavoriteActors:                      &Preference[[]string]{Value: []string{"Chris Evans", "Josh Brolin"}, Weight: 10},
		FavoritePlotElements:                &Preference[[]string]{Value: []string{"love", "family"}, Weight: 10},
		MinimumRottenTomatoesScoreInclusive: &Preference[uint]{Value: 70, Weight: 10},
	}

	tests := []struct {
		prefName       string
		prefVal        reflect.Value
		movie          *Movie
		expectedSat    float32
		expectedWeight uint
	}{
		// good movie
		{prefName: "AfterYearInclusive", prefVal: reflect.ValueOf(prefs.AfterYearInclusive), movie: &goodMovie, expectedSat: 10, expectedWeight: 10},
		{prefName: "BeforeYearExclusive", prefVal: reflect.ValueOf(prefs.BeforeYearExclusive), movie: &goodMovie, expectedSat: 10, expectedWeight: 10},
		{prefName: "MaximumAgeRatingInclusive", prefVal: reflect.ValueOf(prefs.MaximumAgeRatingInclusive), movie: &goodMovie, expectedSat: 10, expectedWeight: 10},
		{prefName: "ShorterThanExclusive", prefVal: reflect.ValueOf(prefs.ShorterThanExclusive), movie: &goodMovie, expectedSat: 10, expectedWeight: 10},
		{prefName: "FavoriteGenre", prefVal: reflect.ValueOf(prefs.FavoriteGenre), movie: &goodMovie, expectedSat: 10, expectedWeight: 10},
		{prefName: "LeastFavoriteDirector", prefVal: reflect.ValueOf(prefs.LeastFavoriteDirector), movie: &goodMovie, expectedSat: 10, expectedWeight: 10},
		{prefName: "FavoriteActors", prefVal: reflect.ValueOf(prefs.FavoriteActors), movie: &goodMovie, expectedSat: 10, expectedWeight: 10},
		{prefName: "FavoritePlotElements", prefVal: reflect.ValueOf(prefs.FavoritePlotElements), movie: &goodMovie, expectedSat: 10, expectedWeight: 10},
		{prefName: "MinimumRottenTomatoesScoreInclusive", prefVal: reflect.ValueOf(prefs.MinimumRottenTomatoesScoreInclusive), movie: &goodMovie, expectedSat: 10, expectedWeight: 10},
		// bad movie
		{prefName: "AfterYearInclusive", prefVal: reflect.ValueOf(prefs.AfterYearInclusive), movie: &badMovie, expectedSat: -10, expectedWeight: 10},
		{prefName: "BeforeYearExclusive", prefVal: reflect.ValueOf(prefs.BeforeYearExclusive), movie: &badMovieBeforeYearExclusive, expectedSat: -10, expectedWeight: 10},
		{prefName: "MaximumAgeRatingInclusive", prefVal: reflect.ValueOf(prefs.MaximumAgeRatingInclusive), movie: &badMovie, expectedSat: -10, expectedWeight: 10},
		{prefName: "ShorterThanExclusive", prefVal: reflect.ValueOf(prefs.ShorterThanExclusive), movie: &badMovie, expectedSat: -10, expectedWeight: 10},
		{prefName: "FavoriteGenre", prefVal: reflect.ValueOf(prefs.FavoriteGenre), movie: &badMovie, expectedSat: -10, expectedWeight: 10},
		{prefName: "LeastFavoriteDirector", prefVal: reflect.ValueOf(prefs.LeastFavoriteDirector), movie: &badMovie, expectedSat: -10, expectedWeight: 10},
		{prefName: "FavoriteActors", prefVal: reflect.ValueOf(prefs.FavoriteActors), movie: &badMovie, expectedSat: -10, expectedWeight: 10},
		{prefName: "FavoritePlotElements", prefVal: reflect.ValueOf(prefs.FavoritePlotElements), movie: &badMovie, expectedSat: -10, expectedWeight: 10},
		{prefName: "MinimumRottenTomatoesScoreInclusive", prefVal: reflect.ValueOf(prefs.MinimumRottenTomatoesScoreInclusive), movie: &badMovie, expectedSat: -10, expectedWeight: 10},
		// partial movie
		{prefName: "FavoriteActors", prefVal: reflect.ValueOf(prefs.FavoriteActors), movie: &partialMovie, expectedSat: 5, expectedWeight: 10},
		{prefName: "FavoritePlotElements", prefVal: reflect.ValueOf(prefs.FavoritePlotElements), movie: &partialMovie, expectedSat: 5, expectedWeight: 10},
	}

	for _, test := range tests {
		sat, weight := calcSat(test.prefName, &test.prefVal, test.movie)
		assert.Equal(t, test.expectedSat, sat)
		assert.Equal(t, test.expectedWeight, weight)
	}
}
