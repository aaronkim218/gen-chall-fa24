package main

type Movie struct {
	Title          string
	Year           uint
	Rated          string
	Runtime        uint // minutes
	Genres         []string
	Director       string
	Actors         []string
	Plot           string
	RottenTomatoes uint // percent
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
