package album

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"example/web-service-gin/album/private/persistence"

	"github.com/gin-gonic/gin"
)

type oData struct {
	Url string `json:"href"`
	All string `json:"collection"`
}

// album represents data about a record album.
type albumDto struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Artist     string  `json:"artist"`
	Price      float64 `json:"price"`
	Discount   float32 `json:"discount"`
	FinalPrice float64 `json:"final-price"`
	Link       oData   `json:"$links"`
}

func newAlbum(a albumDto) persistence.Album {
	return persistence.NewAlbum(a.ID, a.Title, a.Artist, a.Price)
}

func newAlbumDto(a persistence.Album) albumDto {
	return albumDto{
		ID:         a.ID,
		Title:      a.Title,
		Artist:     a.Artist,
		Price:      a.Price,
		Discount:   0.5,
		FinalPrice: float64(a.Price) * 0.5,
	}
}

func (s Controller) addOData(a albumDto) albumDto {
	a.Link = oData{
		Url: fmt.Sprintf("%s/%s", s.config.GetBaseUrl(), a.ID),
		All: s.config.GetBaseUrl(),
	}
	return a
}

func getIntParam(c *gin.Context, p string) (pv int64, ok bool) {
	if s, ok := c.GetQuery(p); ok {
		if v, err := strconv.Atoi(s); err == nil {
			return int64(v), true
		}
	}
	return 0, false
}

// getAlbums responds with the list of all albums as JSON.
func (s Controller) GetAlbums(c *gin.Context) {
	toCtx, cancel := context.WithTimeout(c, 1*time.Millisecond)
	ch := make(chan []albumDto, 1)

	go func() {
		defer cancel()
		defer close(ch)

		sleepMs, sleepOk := getIntParam(c, "sleep")
		log.Printf("sleep: %v %d\n", sleepOk, sleepMs)

		albums, err := s.r.GetAlbums(c)
		if err != nil {
			if err := c.AbortWithError(http.StatusBadRequest, toCtx.Err()); err != nil {
				log.Println(err)
			}

			return
		}

		albumsDto := make([]albumDto, 0, len(albums))
	outher:
		for _, a := range albums {
			select {
			case <-toCtx.Done():
				break outher
			default:
				albumsDto = append(albumsDto, s.addOData(newAlbumDto(a)))
				log.Printf("Album: %s\n", a.ID)
				if sleepOk {
					time.Sleep(time.Duration(sleepMs) * time.Millisecond)
				}
			}
		}
		ch <- albumsDto
	}()

	select {
	case <-toCtx.Done():
		err := c.AbortWithError(http.StatusRequestTimeout, toCtx.Err())
		log.Printf("Error: %v", err)
	case albumsDto := <-ch:
		c.JSON(http.StatusOK, albumsDto)
	}
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func (s Controller) GetAlbumByID(c *gin.Context) {
	if id := c.Param("id"); len(id) != 0 {
		a, err := s.r.GetAlbumByID(id, c)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
		}
		c.JSON(http.StatusOK, s.addOData(newAlbumDto(a)))
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

// postAlbums adds an album from JSON received in the request body.
func (s Controller) PostAlbums(c *gin.Context) {
	var nAlbumDto albumDto

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&nAlbumDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "album not parsed"})
		return
	}
	a := newAlbum(nAlbumDto)
	if err := s.r.AddAlbum(a, c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "album not accepted"})
		return
	}
	c.IndentedJSON(http.StatusCreated, s.addOData(nAlbumDto))
}

type Repository interface {
	GetAlbums(context.Context) ([]persistence.Album, error)
	GetAlbumByID(string, context.Context) (persistence.Album, error)
	AddAlbum(persistence.Album, context.Context) error
}

type Config interface {
	GetBaseUrl() string
}

type Controller struct {
	r      Repository
	config Config
}

func NewController(r Repository, config Config) Controller {
	return Controller{r: r, config: config}
}
