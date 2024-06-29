package spotify

import (
	"net/http"
	"spotify-charter/model"
)

type Album struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Images []Image `json:"images"`
}

type Artists struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Image struct {
	URL   string `json:"url"`
	Width uint   `json:"width"`
}

type Track struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Album   Album     `json:"album"`
	Artists []Artists `json:"artists"`
}

type Item struct {
	Track Track `json:"track"`
}

type GetPlaylistResp struct {
	Items []Item `json:"items"`
}

func (c ApiClient) GetPlaylist(id string) ([]*model.Track, error) {
	req, err := http.NewRequest("GET", baseURL+"/v1/playlists/"+id+"/tracks", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("fields", "items(track(album(id,name,images(url,width)),artists(id,name),id,name))")
	query.Add("limit", "5")

	req.URL.RawQuery = query.Encode()

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, regErrRespToErr(&res.Body)
	}

	resp, err := decodeResp[GetPlaylistResp](&res.Body)
	if err != nil {
		return nil, err
	}

	tracks := make([]*model.Track, 0)

	for _, spotifyTrack := range resp.Items {
		tracks = append(tracks, spotifyTrackToTrack(&spotifyTrack.Track))
	}

	return tracks, nil
}

func spotifyTrackToTrack(track *Track) *model.Track {
	album := model.Album{
		SpotifyID: track.Album.ID,
		Name:      track.Name,
		Images:    make([]model.Image, 0),
	}

	for _, spotifyImage := range track.Album.Images {
		image := model.Image{
			URL:   spotifyImage.URL,
			Width: spotifyImage.Width,
		}

		album.Images = append(album.Images, image)
	}

	artists := make([]model.Artist, 0)

	for _, spotifyArtist := range track.Artists {
		artist := model.Artist{
			SpotifyID: spotifyArtist.ID,
			Name:      spotifyArtist.Name,
		}

		artists = append(artists, artist)
	}

	return &model.Track{
		SpotifyID: track.ID,
		Name:      track.Name,
		Album:     album,
		Artists:   artists,
	}
}
