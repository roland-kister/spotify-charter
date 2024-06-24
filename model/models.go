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

type TopTracks struct {
	Country  *Country
	Track    *Track
	Date     time.Time
	Position int
}
