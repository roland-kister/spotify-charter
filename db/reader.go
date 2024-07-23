package db

import (
	"database/sql"
	"spotify-charter/model"
)

const (
	selPlaylistCountries = iota
	selChartTracks
	selArtistsByTrack
	setImagesByAlbum
)

var readerSqls = map[int]string{
	selPlaylistCountries: `
		SELECT code, name, top_playlist_id
			FROM countries
		WHERE top_playlist_id IS NOT NULL;`,

	selChartTracks: `
		SELECT ct.country_code, ct.position, ct.track_id, t.name AS track_name, t.album_id, a.name AS album_name
			FROM chart_tracks ct
			RIGHT JOIN tracks t ON t.spotify_id = ct.track_id 
			RIGHT JOIN albums a ON a.spotify_id = t.album_id 
		WHERE ct.chart_type = :chart_type AND ct.date = :date;`,

	selArtistsByTrack: `
		SELECT at.artist_id, a.name 
			FROM artists_tracks at
			INNER JOIN artists a ON a.spotify_id = at.artist_id 
		WHERE at.track_id = :track_id;`,

	setImagesByAlbum: `
		SELECT i.url, i.width FROM images i
			WHERE i.album_id = :album_id;`,
}

type Reader struct {
	db    *sql.DB
	stmts map[int]*sql.Stmt
}

func NewReader(db *sql.DB) *Reader {
	var err error

	reader := &Reader{
		db:    db,
		stmts: make(map[int]*sql.Stmt),
	}

	for index, sql := range readerSqls {
		if reader.stmts[index], err = reader.db.Prepare(sql); err != nil {
			panic(err)
		}
	}

	return reader
}

func (reader *Reader) Close() {
	reader.db = nil

	for index := range readerSqls {
		if err := reader.stmts[index].Close(); err != nil {
			panic(err)
		}
	}

}

func (reader *Reader) GetCountriesWithPlaylist() []*model.Country {
	rows, err := reader.stmts[selPlaylistCountries].Query()
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	countries := make([]*model.Country, 0)

	for rows.Next() {
		country := model.Country{}
		if err := rows.Scan(&country.Code, &country.Name, &country.TopPlaylistID); err != nil {
			panic(err)
		}

		countries = append(countries, &country)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return countries
}

func (reader *Reader) GetChartTracksExt(date int64) *model.ChartTracksExt {
	rows, err := reader.stmts[selChartTracks].Query(
		sql.Named("chart_type", model.DailyTopTrack),
		sql.Named("date", date))

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	chartTracks := make(model.ChartTracksExt)

	for rows.Next() {
		track := model.TrackExt{}

		var countryCode string
		var position int

		if err := rows.Scan(&countryCode, &position, &track.ID, &track.Name, &track.Album.ID, &track.Album.Name); err != nil {
			panic(err)
		}

		if chartTracks[countryCode] == nil {
			chartTracks[countryCode] = make([]*model.TrackExt, 5)
		}

		track.Artists = reader.getArtistsForTrack(track.ID)

		track.Album.Images = reader.getImagesForAlbum(track.Album.ID)

		chartTracks[countryCode][position] = &track
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}

	return &chartTracks
}

func (reader *Reader) getArtistsForTrack(trackID string) []model.ArtistExt {
	rows, err := reader.stmts[selArtistsByTrack].Query(sql.Named("track_id", trackID))
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	artists := make([]model.ArtistExt, 0)

	for rows.Next() {
		artist := model.ArtistExt{}

		if err := rows.Scan(&artist.ID, &artist.Name); err != nil {
			panic(err)
		}

		artists = append(artists, artist)
	}

	return artists
}

func (reader *Reader) getImagesForAlbum(albumID string) []model.ImageExt {
	rows, err := reader.stmts[setImagesByAlbum].Query(sql.Named("album_id", albumID))
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	images := make([]model.ImageExt, 0)

	for rows.Next() {
		image := model.ImageExt{}

		if err := rows.Scan(&image.URL, &image.Width); err != nil {
			panic(err)
		}

		images = append(images, image)
	}

	return images
}
