package model

import "time"

type Country struct {
	CountryCode   string
	Name          string
	TopPlaylistID string
}

type Album struct {
	SpotifyID string
	Name      string
	Images    []Image
}

type Image struct {
	URL    string
	Height int
	Width  int
}

type Track struct {
	SpotifyID string
	Name      string
	Album     Album
}

type TopTracks struct {
	Country  *Country
	Track    *Track
	Date     time.Time
	Position int
}
