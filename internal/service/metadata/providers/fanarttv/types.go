package fanarttv

// MovieResponse is the Fanart.tv API response for movie images.
type MovieResponse struct {
	Name           string        `json:"name"`
	TMDbID         string        `json:"tmdb_id"`
	IMDbID         string        `json:"imdb_id"`
	HDMovieLogos   []FanartImage `json:"hdmovielogo"`
	MovieLogos     []FanartImage `json:"movielogo"`
	HDClearArt     []FanartImage `json:"hdmovieclearart"`
	MovieArt       []FanartImage `json:"movieart"`
	MovieDiscs     []FanartImage `json:"moviedisc"`
	MoviePosters   []FanartImage `json:"movieposter"`
	MovieBackdrops []FanartImage `json:"moviebackground"`
	MovieBanners   []FanartImage `json:"moviebanner"`
	MovieThumbs    []FanartImage `json:"moviethumb"`
}

// TVShowResponse is the Fanart.tv API response for TV show images.
type TVShowResponse struct {
	Name          string        `json:"name"`
	TVDbID        string        `json:"thetvdb_id"`
	HDTVLogos     []FanartImage `json:"hdtvlogo"`
	ClearLogos    []FanartImage `json:"clearlogo"`
	HDClearArt    []FanartImage `json:"hdclearart"`
	ClearArt      []FanartImage `json:"clearart"`
	ShowBackdrops []FanartImage `json:"showbackground"`
	TVThumbs      []FanartImage `json:"tvthumb"`
	TVBanners     []FanartImage `json:"tvbanner"`
	CharacterArt  []FanartImage `json:"characterart"`
	TVPosters     []FanartImage `json:"tvposter"`
	SeasonPosters []SeasonImage `json:"seasonposter"`
	SeasonThumbs  []SeasonImage `json:"seasonthumb"`
	SeasonBanners []SeasonImage `json:"seasonbanner"`
}

// FanartImage represents a single image from Fanart.tv.
type FanartImage struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Lang  string `json:"lang"`
	Likes string `json:"likes"`
}

// SeasonImage is a season-specific image with season number.
type SeasonImage struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Lang   string `json:"lang"`
	Likes  string `json:"likes"`
	Season string `json:"season"`
}
