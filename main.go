package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Person struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var db *gorm.DB
var err error

func main() {
	db, err = gorm.Open("sqlite3", "gorm.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	db.AutoMigrate(&Person{})

	r := gin.Default()
	r.GET("/", GetPeople)
	r.POST("/users", CreatePerson)
	r.GET("/users/:id", GetPerson)
	r.PATCH("/users/:id", UpdatePerson)
	r.DELETE("/users/:id", DeletePerson)

	r.Run()
}

func GetPeople(c *gin.Context) {
	var people []Person
	if err := db.Find(&people).Error; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, people)
}

func GetPerson(c *gin.Context) {
	id := c.Params.ByName("id")

	var person Person
	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	} else {
		c.JSON(http.StatusOK, person)
	}
}

func CreatePerson(c *gin.Context) {
	var person Person
	c.BindJSON(&person)

	db.Create(&person)
	c.JSON(http.StatusCreated, person)
}

func UpdatePerson(c *gin.Context) {
	var person Person
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).Find(&person); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		fmt.Println(err)
	}
	c.BindJSON(&person)

	db.Save(&person)
	c.JSON(http.StatusOK, person)
}

func DeletePerson(c *gin.Context) {
	var person Person
	id := c.Params.ByName("id")

	d := db.Where("id = ?", id).Delete(&person)
	fmt.Println(d)
	c.JSON(http.StatusOK, gin.H{"id #" + id: "deleted"})
}

/*
# GET /
GET http://localhost:8080

# GET /users/:id

GET http://localhost:8080/users/1

# POST /users
POST http://localhost:8080/users
Content-Type: application/json
{"first_name": "Nam", "last_name": "Tran Xuan"}

# PATCH /users/:id
PATCH http://localhost:8080/users/1
Context-Type: application/json
{"first_name: "Nam Handsome", "last_name: "Tran Xuan"}

# DELETE /users/:id
DELETE http://localhost:8080/users/1
*/
