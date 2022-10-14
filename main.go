package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)

type Person struct {
	ID		string	`json:"id"`
	Name	string	`json:"name"`
	City	string	`json:"city"`
	Year	int		`json:"year"`
}

var persons = []Person{
	{ID: "1", Name: "Deyvid Santos", City: "Recife", Year: 2002},
	{ID: "2", Name: "Kento Yamazaki", City: "Toquio", Year: 1994},
	{ID: "3", Name: "Uncle Bob", City: "Califórnia", Year: 1952},
}
func getPersons(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, persons)
}

func personById(c *gin.Context){
	id := c.Param("id")
	person, err := getPersonById(id)
	
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Person not found!"})
		return 
	}

	c.IndentedJSON(http.StatusOK, person)
}

func checkoutPerson(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
	}

	person, err := getPersonById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Person not found!"})
		return 
	}

	if person.Year <= 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Person not avaiable to sell"})
		return 
	}

	person.Year -= 1
	c.IndentedJSON(http.StatusOK, person)

}

func returnPerson(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	person, err := getPersonById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Person not found."})
		return
	}

	person.Year += 1
	c.IndentedJSON(http.StatusOK, person)
}

func getPersonById(id string) (*Person, error){
	for i, p := range persons {
		if p.ID == id {
			return &persons[i], nil
		}
	}

	return nil, errors.New("Person not found")
}

func createPerson(c *gin.Context){
	var newPerson Person
	
	if err := c.BindJSON(&newPerson); err != nil {
		return 
	}

	// Raw input variable - don't have database yet
	persons = append(persons, newPerson)
	c.IndentedJSON(http.StatusCreated, newPerson)
}

func main(){
	router := gin.Default()
	router.GET("/persons", getPersons)
	router.GET("/persons/:id", personById)
	router.POST("/persons", createPerson)
	router.PATCH("/checkout", checkoutPerson)
	router.PATCH("/return", returnPerson)
	router.Run("127.0.0.1:9000")
}