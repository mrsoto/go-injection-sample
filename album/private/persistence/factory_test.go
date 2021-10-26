package persistence

import (
	"math"
	"testing"

	"github.com/go-playground/assert/v2"
)

func Test_factory(t *testing.T) {
	a := NewAlbum("id1", "foo", "Jhon", 100)
	assert.Equal(t, a.ID, "id1")
	assert.Equal(t, a.Title, "foo")
	assert.Equal(t, a.Artist, "Jhon")
	assert.Equal(t, true, math.Abs(a.Price-100.0) < 1e-6)
}
