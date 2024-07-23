package model

type ArtistExt struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AlbumExt struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Images []ImageExt `json:"images"`
}

type ImageExt struct {
	URL   string `json:"url"`
	Width uint   `json:"width"`
}

type TrackExt struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Album   AlbumExt    `json:"album"`
	Artists []ArtistExt `json:"artists"`
}

type ChartTracksExt = map[string][]*TrackExt
