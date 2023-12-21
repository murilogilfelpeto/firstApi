package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

type errorDto struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Field     string    `json:"field"`
}

var books = []book{
	book{ID: "1", Title: "Golang pointers", Author: "Mr. Golang", Quantity: 10},
	book{ID: "2", Title: "Goroutines", Author: "Mr. Goroutine", Quantity: 20},
	book{ID: "3", Title: "Golang routers", Author: "Mr. Router", Quantity: 30},
	book{ID: "4", Title: "Golang concurrency", Author: "Mr. Currency", Quantity: 40},
	book{ID: "5", Title: "Golang good parts", Author: "Mr. Good", Quantity: 50},
}

func main() {
	startServer()
}

func startServer() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/books/:id", getById)
	router.PATCH("/books/checkout", checkoutBook)
	router.PATCH("/books/checkin", checkinBook)
	router.DELETE("/books/:id", deleteBook)
	err := router.Run("localhost:8080")

	if err != nil {
		return
	}
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	lastId, _ := strconv.Atoi(books[len(books)-1].ID)
	newBook.ID = strconv.Itoa(lastId + 1)

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getById(c *gin.Context) {
	id := c.Param("id")

	book, err := findBookById(id)
	if err != nil {
		errorDto := errorDto{
			Message:   "Book not found",
			Timestamp: time.Now(),
			Field:     "id",
		}
		c.IndentedJSON(http.StatusNotFound, errorDto)
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func findBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		errorDto := errorDto{
			Message:   "Parameter id is mandatory",
			Timestamp: time.Now(),
			Field:     "id",
		}
		c.IndentedJSON(http.StatusBadRequest, errorDto)
		return
	}

	book, err := findBookById(id)
	if err != nil {
		errorDto := errorDto{
			Message:   "Book not found",
			Timestamp: time.Now(),
			Field:     "id",
		}
		c.IndentedJSON(http.StatusNotFound, errorDto)
		return
	}

	if book.Quantity <= 0 {
		errorDto := errorDto{
			Message:   "Book not available",
			Timestamp: time.Now(),
			Field:     "quantity",
		}
		c.IndentedJSON(http.StatusBadRequest, errorDto)
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func checkinBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		errorDto := errorDto{
			Message:   "Parameter id is mandatory",
			Timestamp: time.Now(),
			Field:     "id",
		}
		c.IndentedJSON(http.StatusBadRequest, errorDto)
		return
	}

	book, err := findBookById(id)
	if err != nil {
		errorDto := errorDto{
			Message:   "Book not found",
			Timestamp: time.Now(),
			Field:     "id",
		}
		c.IndentedJSON(http.StatusNotFound, errorDto)
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func deleteBook(c *gin.Context) {
	id := c.Param("id")

	_, err := findBookById(id)
	if err != nil {
		errorDto := errorDto{
			Message:   "Book not found",
			Timestamp: time.Now(),
			Field:     "id",
		}
		c.IndentedJSON(http.StatusNotFound, errorDto)
		return
	}
	removeBookFromList(id)
	c.IndentedJSON(http.StatusOK, "")
}

func removeBookFromList(id string) {
	for i, b := range books {
		if b.ID == id {
			books = append(books[:i], books[i+1:]...)
		}
	}
}
