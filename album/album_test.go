package album

import (
	"context"
	"encoding/json"
	"errors"
	"example/web-service-gin/album/private/persistence"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/go-test/deep"
)

type RepositoryStub struct {
	Repository
}

type RepositoryFailuresStub struct {
	Repository
}

var albums []persistence.Album
var albumsDto []albumDto

func init() {
	albums = []persistence.Album{
		{ID: "1", Title: "Foo", Artist: "Jhon", Price: 100},
		{ID: "2", Title: "Bar", Artist: "Bill", Price: 200},
	}

	albumsDto = []albumDto{
		{ID: "1", Title: "Foo", Artist: "Jhon", Price: 100, Discount: 0.5, FinalPrice: 50,
			Link: oData{Url: "https://sample.com/albums/1", All: "https://sample.com/albums"},
		},
		{ID: "2", Title: "Bar", Artist: "Bill", Price: 200, Discount: 0.5, FinalPrice: 100,
			Link: oData{Url: "https://sample.com/albums/2", All: "https://sample.com/albums"},
		},
	}
}

func (r RepositoryStub) GetAlbumByID(id string, c context.Context) (persistence.Album, error) {
	if id == "0" {
		return persistence.Album{}, errors.New("not found")
	}
	return albums[0], nil
}

func (r RepositoryStub) GetAlbums(c context.Context) ([]persistence.Album, error) {
	return albums, nil
}

func (r RepositoryFailuresStub) GetAlbums(c context.Context) ([]persistence.Album, error) {
	return nil, errors.New("failure")
}

type configStub struct{}

func (c configStub) GetBaseUrl() string {
	return "https://sample.com/albums"
}

func setGetAlbumRouter(url string) (*http.Request, *httptest.ResponseRecorder) {
	r := gin.New()
	ctrl := NewController(RepositoryStub{}, configStub{})

	r.GET("/albums", ctrl.GetAlbums)
	r.GET("/albums/:id", ctrl.GetAlbumByID)
	r.POST("/albums", ctrl.PostAlbums)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

func setGetAlbumsRouter(ctrl Controller, url string) (*http.Request, *httptest.ResponseRecorder) {
	r := gin.New()

	r.GET("/albums", ctrl.GetAlbums)
	r.POST("/albums", ctrl.PostAlbums)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}
func Test_GetAlbumByID_whenFound(t *testing.T) {
	_, w := setGetAlbumRouter("/albums/1")

	assert.Equal(t, http.StatusOK, w.Code)
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
	}

	actual := albumDto{}
	if err := json.Unmarshal(body, &actual); err != nil {
		t.Error(err)
	}
	expected := albumsDto[0]

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Error(diff)
	}

}

func Test_GetAlbumByID_whenNotFound(t *testing.T) {
	_, w := setGetAlbumRouter("/albums/0")

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test_GetAlbumByID_whenInvalidRequest(t *testing.T) {
	_, w := setGetAlbumRouter("/albums/")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_GetAlbumsWhenOk(t *testing.T) {
	ctrl := NewController(RepositoryStub{}, configStub{})

	_, w := setGetAlbumsRouter(ctrl, "/albums")
	assert.Equal(t, http.StatusOK, w.Code)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
	}
	actual := []albumDto{}
	if err := json.Unmarshal(body, &actual); err != nil {
		t.Error(err)
	}
	expected := albumsDto

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Error(diff)
	}
}

func Test_GetAlbumsWhenTimeOut(t *testing.T) {
	ctrl := NewController(RepositoryStub{}, configStub{})
	_, w := setGetAlbumsRouter(ctrl, "/albums?sleep=10")
	assert.Equal(t, http.StatusRequestTimeout, w.Code)

}

func Test_GetAlbumsWhenDBFail(t *testing.T) {
	ctrl := NewController(RepositoryFailuresStub{}, configStub{})
	_, w := setGetAlbumsRouter(ctrl, "/albums")
	assert.Equal(t, http.StatusBadRequest, w.Code)

}
