package album

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"example/web-service-gin/album/private/persistence"
	"io/ioutil"
	"log"
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

func (r RepositoryStub) AddAlbum(a persistence.Album, c context.Context) (persistence.Album, error) {
	na := a
	na.ID = "3"
	return na, nil
}

func (r RepositoryStub) GetAlbumByID(id string, c context.Context) (persistence.Album, error) {
	return albums[0], nil
}

func (r RepositoryStub) GetAlbums(c context.Context) ([]persistence.Album, error) {
	return albums, nil
}

func (r RepositoryFailuresStub) GetAlbumByID(id string, c context.Context) (persistence.Album, error) {
	return persistence.Album{}, errors.New("not found")
}

func (r RepositoryFailuresStub) GetAlbums(c context.Context) ([]persistence.Album, error) {
	return nil, errors.New("failure")
}

func (r RepositoryFailuresStub) AddAlbum(a persistence.Album, c context.Context) (persistence.Album, error) {
	return persistence.Album{}, errors.New("DB Error")
}

type configStub struct{}

func (c configStub) GetBaseUrl() string {
	return "https://sample.com/albums"
}

func setGetAlbumsRouter(ctrl Controller, url string) (*http.Request, *httptest.ResponseRecorder) {
	r := gin.New()

	r.GET("/albums", ctrl.GetAlbums)
	r.GET("/albums/:id", ctrl.GetAlbumByID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

func setPostAlbumsRouter(ctrl Controller, a albumDto) (*http.Request, *httptest.ResponseRecorder) {
	r := gin.New()

	r.POST("/albums", ctrl.PostAlbums)

	b, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewBuffer(b)
	req, err := http.NewRequest(http.MethodPost, "/albums", body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

func Test_GetAlbumByID_whenFound(t *testing.T) {
	ctrl := NewController(RepositoryStub{}, configStub{})
	_, w := setGetAlbumsRouter(ctrl, "/albums/1")

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
	ctrl := NewController(RepositoryFailuresStub{}, configStub{})
	_, w := setGetAlbumsRouter(ctrl, "/albums/1")

	assert.Equal(t, http.StatusNotFound, w.Code)
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

func Test_PostAlbumByID_whenSuccess(t *testing.T) {
	ctrl := NewController(RepositoryStub{}, configStub{})
	a := albumDto{
		Title:      "Harry",
		Artist:     "Mary",
		Price:      300,
		FinalPrice: 150,
		Discount:   0.5,
	}
	_, w := setPostAlbumsRouter(ctrl, a)

	assert.Equal(t, http.StatusCreated, w.Code)
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error(err)
	}

	actual := albumDto{}
	if err := json.Unmarshal(body, &actual); err != nil {
		t.Error(err)
	}
	expected := a
	expected.ID = "3"
	expected.Link = oData{
		Url: "https://sample.com/albums/3",
		All: "https://sample.com/albums",
	}

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Error(diff)
	}

}

func Test_PostAlbumByID_whenDbFail(t *testing.T) {
	ctrl := NewController(RepositoryFailuresStub{}, configStub{})
	a := albumDto{
		Title:      "Harry",
		Artist:     "Mary",
		Price:      300,
		FinalPrice: 150,
		Discount:   0.5,
	}
	_, w := setPostAlbumsRouter(ctrl, a)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Error("error unexpected")
	}
	payload := string(body[:])
	expected := "{\"message\":\"album not accepted\"}"
	if payload != expected {
		t.Errorf("error actual=%s expected: %s", payload, expected)
	}
}

func setPostAlbumsJsonRouter(ctrl Controller, p string) (*http.Request, *httptest.ResponseRecorder) {
	r := gin.New()

	r.POST("/albums", ctrl.PostAlbums)

	body := bytes.NewBufferString(p)
	req, err := http.NewRequest(http.MethodPost, "/albums", body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return req, w
}

func Test_PostAlbumByID_whenJsonError(t *testing.T) {
	ctrl := NewController(RepositoryStub{}, configStub{})

	_, w := setPostAlbumsJsonRouter(ctrl, "}")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
