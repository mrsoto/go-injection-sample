package album

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string
	Title  string
	Artist string
	Price  float64
}

// album represents data about a record album.
type albumDto struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	Artist     string  `json:"artist"`
	Price      float64 `json:"price"`
	Discount   float32 `json:"discount"`
	FinalPrice float32 `json:"final-price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func newAlbum(a albumDto) album {
	return album{ID: a.ID, Title: a.Title, Artist: a.Artist, Price: a.Price}
}

func newAlbumDto(a album) albumDto {
	return albumDto{
		ID:         a.ID,
		Title:      a.Title,
		Artist:     a.Artist,
		Price:      a.Price,
		Discount:   0.5,
		FinalPrice: float32(a.Price) * 0.5,
	}
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
func (s Services) GetAlbums(c *gin.Context) {
	toCtx, cancel := context.WithTimeout(c, 1*time.Millisecond)
	ch := make(chan []albumDto, 1)

	go func() {
		defer cancel()
		defer close(ch)

		albumsDto := make([]albumDto, 0, len(albums))

		sleepMs, sleepOk := getIntParam(c, "sleep")
		log.Printf("sleep: %v %d\n", sleepOk, sleepMs)

	outher:
		for _, a := range albums {
			select {
			case <-toCtx.Done():
				break outher
			default:
				albumsDto = append(albumsDto, newAlbumDto(a))
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
func (s Services) GetAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.JSON(http.StatusOK, newAlbumDto(a))
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// postAlbums adds an album from JSON received in the request body.
func (s Services) PostAlbums(c *gin.Context) {
	var nAlbum albumDto

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&nAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum(nAlbum))
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

type Services struct {
}

func NewServices() Services {
	return Services{}
}
