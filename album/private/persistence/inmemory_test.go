package persistence

import (
	"context"
	"testing"

	"github.com/go-playground/assert/v2"
)

func Test_GetAlbums(t *testing.T) {

	r := NewInMemoryRepository()
	a, err := r.GetAlbumByID("1", context.Background())
	assert.Equal(t, err, nil)
	assert.Equal(t, a.ID, "1")
}

func Test_GetAlbums_whenFail(t *testing.T) {

	r := NewInMemoryRepository()
	a, err := r.GetAlbumByID("12", context.Background())
	assert.NotEqual(t, err, nil)
	assert.Equal(t, a.ID, "")
}

func Test_GetAlbums_whenTimeOut(t *testing.T) {

	c, cancel := context.WithCancel(context.Background())
	cancel()
	r := NewInMemoryRepository()
	a, err := r.GetAlbumByID("1", c)
	assert.Equal(t, err, context.Canceled)
	assert.Equal(t, a.ID, "")
}

func Test_AddAlbums(t *testing.T) {

	r := NewInMemoryRepository()
	a := Album{ID: "newId"}
	newA, err := r.AddAlbum(a, context.Background())
	assert.Equal(t, err, nil)
	assert.Equal(t, newA.ID, "4")
}
