package persistence

import (
	"context"
	"testing"

	"github.com/go-playground/assert/v2"
)

func Test_GetAlbums(t *testing.T) {

	r := NewInMemoryRepository()
	a, err := r.GetAlbums(context.Background())
	assert.Equal(t, err, nil)
	assert.Equal(t, r.albums, a)
}

func Test_GetAlbum(t *testing.T) {

	r := NewInMemoryRepository()
	a, err := r.GetAlbumByID("1", context.Background())
	assert.Equal(t, err, nil)
	assert.Equal(t, a.ID, "1")
}

func Test_GetAlbum_whenFail(t *testing.T) {

	r := NewInMemoryRepository()
	a, err := r.GetAlbumByID("12", context.Background())
	assert.NotEqual(t, err, nil)
	assert.Equal(t, a.ID, "")
}

func Test_GetAlbum_whenTimeOut(t *testing.T) {

	c, cancel := context.WithCancel(context.Background())
	cancel()
	r := NewInMemoryRepository()

	a, err := r.GetAlbumByID("1", c)

	assert.Equal(t, err, context.Canceled)
	assert.Equal(t, a.ID, "")
}

func Test_AddAlbum(t *testing.T) {

	r := NewInMemoryRepository()
	a := Album{ID: "newId"}

	newA, err := r.AddAlbum(a, context.Background())

	assert.Equal(t, err, nil)
	assert.Equal(t, newA.ID, "4")
}

func Test_AddAlbums_whenCanceled(t *testing.T) {

	r := NewInMemoryRepository()
	c, cancel := context.WithCancel(context.Background())
	cancel()
	a := Album{ID: "newId"}

	newA, err := r.AddAlbum(a, c)

	assert.NotEqual(t, err, nil)
	assert.Equal(t, newA, Album{})
}

func Test_AddAlbums_whenInvalidData(t *testing.T) {

	r := NewInMemoryRepository()
	r.albums[0].ID = "a"
	a := Album{ID: "newId"}

	newA, err := r.AddAlbum(a, context.Background())

	assert.NotEqual(t, err, nil)
	assert.Equal(t, newA, Album{})
}
