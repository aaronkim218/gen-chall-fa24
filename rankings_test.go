package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcRuntime(t *testing.T) {
	tests := []struct {
		runtime  string
		expected float64
	}{
		{"0h0m0s", 0},
		{"0h0m1s", float64(1) / 60},
		{"0h1m0s", 1},
		{"0h1m1s", 1 + (float64(1) / 60)},
		{"1h0m0s", 60},
		{"1h0m1s", 60 + (float64(1) / 60)},
		{"1h1m1s", 60 + 1 + (float64(1) / 60)},
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

func TestScore(t *testing.T) {
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

	prefs1 := Preferences{
		AfterYearInclusive:        &Preference[uint]{Value: 2000, Weight: 10},
		MaximumAgeRatingInclusive: &Preference[string]{Value: "PG-13", Weight: 10},
	}

	prefs2 := Preferences{
		ShorterThanExclusive: &Preference[string]{Value: "1h45m0s", Weight: 10},
		FavoriteGenre:        &Preference[string]{Value: "Action", Weight: 10},
	}

	prefs3 := Preferences{
		LeastFavoriteDirector:               &Preference[string]{Value: "Ridley Scott", Weight: 10},
		FavoriteActors:                      &Preference[[]string]{Value: []string{"Chris Evans", "Josh Brolin"}, Weight: 10},
		FavoritePlotElements:                &Preference[[]string]{Value: []string{"love", "family"}, Weight: 10},
		MinimumRottenTomatoesScoreInclusive: &Preference[uint]{Value: 70, Weight: 10},
	}

	prefsBadMovieBeforeYearExclusive := Preferences{
		BeforeYearExclusive: &Preference[uint]{Value: 2010, Weight: 10},
	}

	person1 := Person{
		Preferences: &prefs1,
	}

	person2 := Person{
		Preferences: &prefs2,
	}

	person3 := Person{
		Preferences: &prefs3,
	}

	personBadMovieBeforeYearExclusive := Person{
		Preferences: &prefsBadMovieBeforeYearExclusive,
	}

	tests := []struct {
		movie    *Movie
		person   *Person
		expected int
	}{
		// good movie
		{movie: &goodMovie, person: &person1, expected: 20},
		{movie: &goodMovie, person: &person2, expected: 20},
		{movie: &goodMovie, person: &person3, expected: 50},

		// bad movie
		{movie: &badMovie, person: &person1, expected: 0},
		{movie: &badMovie, person: &person2, expected: 0},
		{movie: &badMovie, person: &person3, expected: -10},
		{movie: &badMovieBeforeYearExclusive, person: &personBadMovieBeforeYearExclusive, expected: 0},

		// partial
		{movie: &partialMovie, person: &person3, expected: 20},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, score(test.movie, test.person))
	}
}

func TestCalcScores(t *testing.T) {
	movie1 := Movie{
		ID:             "first",
		Title:          "First Movie",
		Year:           2005,
		Rated:          "PG-13",
		Runtime:        100,
		Genres:         []string{"Action", "Comedy", "Drama"},
		Director:       "Aaron Kim",
		Actors:         []string{"Robert Pattinson", "Josh Brolin", "Chris Evans"},
		Plot:           "this is a movie about family, war, and love",
		RottenTomatoes: 85,
	}

	expectedMovie1 := Movie{
		ID:             "first",
		Title:          "First Movie",
		Year:           2005,
		Rated:          "PG-13",
		Runtime:        100,
		Genres:         []string{"Action", "Comedy", "Drama"},
		Director:       "Aaron Kim",
		Actors:         []string{"Robert Pattinson", "Josh Brolin", "Chris Evans"},
		Plot:           "this is a movie about family, war, and love",
		RottenTomatoes: 85,
		Score:          180,
	}

	movie2 := Movie{
		ID:             "second",
		Title:          "Second Movie",
		Year:           1990,
		Rated:          "R",
		Runtime:        200,
		Genres:         []string{"Thriller", "Horror"},
		Director:       "Ridley Scott",
		Actors:         []string{"Kevin Durant", "Jon Jones", "Dana White"},
		Plot:           "this is a documentary about history",
		RottenTomatoes: 60,
	}

	expectedMovie2 := Movie{
		ID:             "second",
		Title:          "Second Movie",
		Year:           1990,
		Rated:          "R",
		Runtime:        200,
		Genres:         []string{"Thriller", "Horror"},
		Director:       "Ridley Scott",
		Actors:         []string{"Kevin Durant", "Jon Jones", "Dana White"},
		Plot:           "this is a documentary about history",
		RottenTomatoes: 60,
		Score:          10,
	}

	movie3 := Movie{
		ID:             "third",
		Title:          "Third Movie",
		Year:           2020,
		Rated:          "PG",
		Runtime:        110,
		Genres:         []string{"Drama"},
		Director:       "Stanley Kubrick",
		Actors:         []string{"Jack Nicholson", "Jon Jones", "Dana White"},
		Plot:           "this is a movie about aliens and space",
		RottenTomatoes: 80,
	}

	expectedMovie3 := Movie{
		ID:             "third",
		Title:          "Third Movie",
		Year:           2020,
		Rated:          "PG",
		Runtime:        110,
		Genres:         []string{"Drama"},
		Director:       "Stanley Kubrick",
		Actors:         []string{"Jack Nicholson", "Jon Jones", "Dana White"},
		Plot:           "this is a movie about aliens and space",
		RottenTomatoes: 80,
		Score:          120,
	}

	person1 := Person{
		Preferences: &Preferences{
			AfterYearInclusive:                  &Preference[uint]{Value: 2000, Weight: 10},
			MaximumAgeRatingInclusive:           &Preference[string]{Value: "PG-13", Weight: 10},
			ShorterThanExclusive:                &Preference[string]{Value: "1h45m0s", Weight: 10},
			FavoriteGenre:                       &Preference[string]{Value: "Action", Weight: 10},
			LeastFavoriteDirector:               &Preference[string]{Value: "Ridley Scott", Weight: 10},
			FavoriteActors:                      &Preference[[]string]{Value: []string{"Chris Evans", "Josh Brolin"}, Weight: 10},
			FavoritePlotElements:                &Preference[[]string]{Value: []string{"love", "family"}, Weight: 10},
			MinimumRottenTomatoesScoreInclusive: &Preference[uint]{Value: 70, Weight: 10},
		},
	}

	person2 := Person{
		Preferences: &Preferences{
			AfterYearInclusive:                  &Preference[uint]{Value: 2000, Weight: 10},
			MaximumAgeRatingInclusive:           &Preference[string]{Value: "R", Weight: 10},
			ShorterThanExclusive:                &Preference[string]{Value: "2h15m0s", Weight: 10},
			FavoriteGenre:                       &Preference[string]{Value: "Comedy", Weight: 10},
			LeastFavoriteDirector:               &Preference[string]{Value: "Christopher Nolan", Weight: 10},
			FavoriteActors:                      &Preference[[]string]{Value: []string{"Jennifer Lawrence", "Amy Schumer"}, Weight: 10},
			FavoritePlotElements:                &Preference[[]string]{Value: []string{"honor", "war"}, Weight: 10},
			MinimumRottenTomatoesScoreInclusive: &Preference[uint]{Value: 60, Weight: 10},
		},
	}

	person3 := Person{
		Preferences: &Preferences{
			AfterYearInclusive:                  &Preference[uint]{Value: 2000, Weight: 10},
			MaximumAgeRatingInclusive:           &Preference[string]{Value: "G", Weight: 10},
			ShorterThanExclusive:                &Preference[string]{Value: "1h35m0s", Weight: 10},
			FavoriteGenre:                       &Preference[string]{Value: "Drama", Weight: 10},
			LeastFavoriteDirector:               &Preference[string]{Value: "Stanley Kubrick", Weight: 10},
			FavoriteActors:                      &Preference[[]string]{Value: []string{"Jack Nicholson", "Tom Hanks"}, Weight: 10},
			FavoritePlotElements:                &Preference[[]string]{Value: []string{"aliens", "space"}, Weight: 10},
			MinimumRottenTomatoesScoreInclusive: &Preference[uint]{Value: 75, Weight: 10},
		},
	}

	tests := []struct {
		movies   []*Movie
		people   []*Person
		expected []*Movie
	}{
		// good movie
		{[]*Movie{&movie1, &movie2, &movie3}, []*Person{&person1, &person2, &person3}, []*Movie{&expectedMovie1, &expectedMovie2, &expectedMovie3}},
	}

	for _, test := range tests {
		calcScores(test.movies, test.people)
		for i := 0; i < len(test.movies); i++ {
			assert.Equal(t, test.expected[i].Score, test.movies[i].Score)
		}
	}
}
