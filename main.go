package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// user represents data about an artist.
type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Genre string `json:"genre"`
}

// album represents data about a record album.
type Album struct {
	Id       int64   `json:"id"`
	Title    string  `json:"title"`
	ArtistId int64   `json:"artistId"`
	Artist   *User   `json:"artist" pg:"rel:has-one"`
	Price    float64 `json:"price"`
}

var db = new(pg.DB)

func main() {
	//DB Connection
	db = db_connection()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	//Request Routing
	router := gin.Default()

	router.GET("/albums", getAlbums)
	router.GET("/users", getUsers)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbum)
	router.POST("/users", postUsers)

	router.Run("localhost:8080")
}

func db_connection() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "rootuser",
		Password: "rootpassword",
	})

	return db
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
		(*Album)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func getAlbums(c *gin.Context) {
	var albums []Album

	err := db.Model(&albums).Select()
	if err != nil {
		return
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	album := new(Album)

	err := db.Model(album).
		Relation("Artist").
		Where("album.id = ?", id).
		Select()

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found!"})
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}

func postAlbum(c *gin.Context) {
	var newAlbum Album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not create album"})
		return
	}

	_, err := db.Model(&newAlbum).Insert()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not create album"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getUsers(c *gin.Context) {
	var users []User

	err := db.Model(&users).Select()
	if err != nil {
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func postUsers(c *gin.Context) {
	var newUser User

	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not create user"})
		return
	}

	_, err := db.Model(&newUser).Insert()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not create user"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newUser)
}

// curl http://localhost:8080/users \
// >     --include \
// >     --header "Content-Type: application/json" \
// >     --request "POST" \
// >     --data '{"name": "Betty Carter","genre": "Jazz"}'

// curl http://localhost:8080/albums \
//     --include \
//     --header "Content-Type: application/json" \
//     --request "POST" \
//     --data '{"id": 1,"title": "The Modern Sound of Betty Carter","artistId":1,"price": 49.99}'

// curl http://localhost:8080/users \
//     --include \
//     --header "Content-Type: application/json" \
//     --request "POST" \
//     --data '{"id":1,"name": "Betty Carter","genre": "Jazz"}'
