package persistence

func NewAlbum(id string, title string, artist string, price float64) Album {
	return Album{ID: id, Title: title, Artist: artist, Price: price}
}
