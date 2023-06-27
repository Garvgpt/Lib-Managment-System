package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type Book struct {
	ID        int
	Title     string
	Author    string
	Publisher string
}

func initDB() {
	db, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/librarygin")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initDB()
	defer db.Close()

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/books")
	})

	router.GET("/books", func(c *gin.Context) {
		books := getBooks()
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Library Management System",
			"books": books,
		})
	})

	router.GET("/books/new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new.html", gin.H{
			"title": "Add New Book",
		})
	})

	router.POST("/books", func(c *gin.Context) {
		title := c.PostForm("title")
		author := c.PostForm("author")
		publisher := c.PostForm("publisher")

		insertBook(title, author, publisher)

		c.Redirect(http.StatusMovedPermanently, "/books")
	})

	router.GET("/books/edit/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)

		book := getBookByID(id)

		c.HTML(http.StatusOK, "edit.html", gin.H{
			"title": "Edit Book",
			"book":  book,
		})
	})

	router.POST("/books/update/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)

		title := c.PostForm("title")
		author := c.PostForm("author")
		publisher := c.PostForm("publisher")

		updateBook(id, title, author, publisher)

		c.Redirect(http.StatusMovedPermanently, "/books")
	})

	router.GET("/books/delete/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.Atoi(idStr)

		deleteBook(id)

		c.Redirect(http.StatusMovedPermanently, "/books")
	})

	router.Run(":8080")
}

func getBooks() []Book {
	rows, err := db.Query("SELECT id, title, author, publisher FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Publisher)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}

	return books
}

func insertBook(title, author, publisher string) {
	_, err := db.Exec("INSERT INTO books (title, author, publisher) VALUES (?, ?, ?)", title, author, publisher)
	if err != nil {
		log.Fatal(err)
	}
}

func getBookByID(id int) Book {
	row := db.QueryRow("SELECT id, title, author, publisher FROM books WHERE id = ?", id)

	var book Book
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Publisher)
	if err != nil {
		log.Fatal(err)
	}

	return book
}

func updateBook(id int, title, author, publisher string) {
	_, err := db.Exec("UPDATE books SET title = ?, author = ?, publisher = ? WHERE id = ?", title, author, publisher, id)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteBook(id int) {
	_, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
}
