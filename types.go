package main

type Movie struct {
	ID             string
	Title          string
	Year           uint
	Rated          string
	Runtime        uint // minutes
	Genres         []string
	Director       string
	Actors         []string
	Plot           string
	RottenTomatoes uint // percent
	MedProp        float64
	AvgProp        float64
}

type IMDBResponse struct {
	Title    string
	Year     string
	Rated    string
	Runtime  string
	Genre    string
	Director string
	Actors   string
	Plot     string
	Ratings  []*Rating
}

type Rating struct {
	Source string
	Value  string
}

type ChallengeResponse struct {
	Token  string  `json:"token"`
	Prompt *Prompt `json:"prompt"`
}

type Prompt struct {
	Movies []string  `json:"movies"`
	People []*Person `json:"people"`
}

type Person struct {
	Name        string       `json:"name"`
	Preferences *Preferences `json:"preferences"`
}

type Preferences struct {
	AfterYearInclusive                  *Preference[uint]     `json:"afterYear(inclusive)"`
	BeforeYearExclusive                 *Preference[uint]     `json:"beforeYear(exclusive)"`
	MaximumAgeRatingInclusive           *Preference[string]   `json:"maximumAgeRating(inclusive)"`
	ShorterThanExclusive                *Preference[string]   `json:"shorterThan(exclusive)"`
	FavoriteGenre                       *Preference[string]   `json:"favoriteGenre"`
	LeastFavoriteDirector               *Preference[string]   `json:"leastFavoriteDirector"`
	FavoriteActors                      *Preference[[]string] `json:"favoriteActors"`
	FavoritePlotElements                *Preference[[]string] `json:"favoritePlotElements"`
	MinimumRottenTomatoesScoreInclusive *Preference[uint]     `json:"minimumRottenTomatoesScore(inclusive)"`
}

type Preference[T any] struct {
	Value  T    `json:"value"`
	Weight uint `json:"weight"`
}

type B struct {
	MAX_YEAR           uint
	MIN_YEAR           uint
	MIN_RUNTIME        uint
	MIN_ROTTENTOMATOES uint
	MIN_RATING         uint
}

var Bounds = B{
	MAX_YEAR:           2024,
	MIN_YEAR:           1888,
	MIN_RUNTIME:        0,
	MIN_ROTTENTOMATOES: 0,
	MIN_RATING:         0,
}

var RATINGS = map[string]uint{
	"G":     0,
	"PG":    1,
	"PG-13": 2,
	"R":     3,
	"NC-17": 4,
}
