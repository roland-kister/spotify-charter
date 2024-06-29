package model

import "time"

type Country struct {
	CountryCode   string
	Name          string
	TopPlaylistID string
}

type Artist struct {
	SpotifyID string
	Name      string
}

type Album struct {
	SpotifyID string
	Name      string
	Images    []Image
}

type Image struct {
	URL   string
	Width uint
}

type Track struct {
	SpotifyID string
	Name      string
	Album     Album
	Artists   []Artist
}

type ChartType string

const (
	DailyTopTrack ChartType = "DAILY_TOP_TRACK"
)

type ChartTrack struct {
	Country   *Country
	Track     *Track
	ChartType ChartType
	Date      int64
	Position  int
}

func TimeToDatestamp(t time.Time) int64 {
	t = t.UTC()

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
}
