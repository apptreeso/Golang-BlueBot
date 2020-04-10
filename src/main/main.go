package main

import (
  "fmt"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"

	"github.com/labstack/echo"
)

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Dog struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Hamster struct {
	Name string `json:"name"`
	Type string `json:"type"`
}


func yallo(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world")
}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")

	dataType := c.Param("data")

	if (dataType == "string") {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is %s\nand type is: %s", catName, catType))	
	}

	if (dataType == "json") {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "you need to let us know if you want json or string data",
	})
}

func addCat(c echo.Context) error {
	cat := Cat{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for addCats: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("Failed unmarshaling in addCats: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("this is your cat: %#v", cat)
	return c.String(http.StatusOK, "we got your cat!")
}

func addDog(c echo.Context) error {
	dog := Dog{}
	
	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("Failed processing addDog request: $%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("this is your dog: %#v", dog)
	return c.String(http.StatusOK, "we got your dog")
}

func addHamster(c echo.Context) error {
	hamster := Hamster{}

	err := c.Bind(&hamster)
	if err != nil {
		log.Printf("Failed processing addHamster request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("this is your hamster: %#v", hamster)
	return c.String(http.StatusOK, "we got your hamster")
}

func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "horay you are on the secret admin main page!")
}

func main() {
  fmt.Println("Welcome")
  
	e := echo.New()

	g := e.Group("/admin")

	g.GET("/main", mainAdmin)

	e.GET("/", yallo)
	e.GET("/cats/:data", getCats) //http://localhost:8080/cats/string?name="myName"&type="myTytpe"

	e.POST("/cats", addCat)	//{"name": "fishmaster",	"type": "cat-fish"}
	e.POST("/dogs", addDog)	//{"name": "doggymaster",	"type": "dog-fish"}
	e.POST("/hamsters", addHamster) //{"name": "hamster",	"type": "hamster-fish"}

	e.Start(":8080")
}
