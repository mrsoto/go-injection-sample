package persistence

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

func (r *Repository) fill() {
	// albums slice to seed record album data.
	var albums = []Album{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
	r.albums = albums
}

// getAlbums responds with the list of all albums.
func (r *Repository) GetAlbums(c context.Context) ([]Album, error) {
	return r.albums, nil
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func (r *Repository) GetAlbumByID(id string, c context.Context) (Album, error) {

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range r.albums {
		select {
		case <-c.Done():
			return Album{}, c.Err()
		default:
			if a.ID == id {
				return a, nil
			}
		}
	}
	return Album{}, errors.New("not found")
}

// AddAlbum adds an album.
func (r *Repository) AddAlbum(a Album, c context.Context) (Album, error) {

	lId := r.albums[len(r.albums)-1].ID
	nId, err := strconv.Atoi(lId)
	if err != nil {
		return Album{}, err
	}
	a.ID = fmt.Sprintf("%d", nId+1)
	// Add the new album to the slice.
	r.albums = append(r.albums, a)
	return a, nil
}

type Repository struct {
	albums []Album
}

// NewInMemoryRepository create a repository instance
func NewInMemoryRepository() *Repository {
	r := &Repository{}
	r.fill()
	return r
}
