package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	ID       int
	Title    string
	Author   string
	ISBN     string
	Quantity int
}

var db *sql.DB
var tpl *template.Template

func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/library")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database.")
}

func main() {
	initDB()
	defer db.Close()

	tpl = template.Must(template.ParseGlob("templates/*.html"))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/delete", deleteHandler)

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Quantity)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, book)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.ExecuteTemplate(w, "index.html", struct{ Books []Book }{books})
	if err != nil {
		log.Fatal(err)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	book := Book{
		Title:    r.FormValue("title"),
		Author:   r.FormValue("author"),
		ISBN:     r.FormValue("isbn"),
		Quantity: parseInt(r.FormValue("quantity")),
	}

	_, err := db.Exec("INSERT INTO books(title, author, isbn, quantity) VALUES(?, ?, ?, ?)",
		book.Title, book.Author, book.ISBN, book.Quantity)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	id := parseInt(r.URL.Query().Get("id"))

	var book Book
	err := db.QueryRow("SELECT * FROM books WHERE id = ?", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Quantity)
	if err != nil {
		log.Fatal(err)
	}

	err = tpl.ExecuteTemplate(w, "edit.html", book)
	if err != nil {
		log.Fatal(err)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	book := Book{
		ID:       parseInt(r.FormValue("id")),
		Title:    r.FormValue("title"),
		Author:   r.FormValue("author"),
		ISBN:     r.FormValue("isbn"),
		Quantity: parseInt(r.FormValue("quantity")),
	}

	_, err := db.Exec("UPDATE books SET title = ?, author = ?, isbn = ?, quantity = ? WHERE id = ?",
		book.Title, book.Author, book.ISBN, book.Quantity, book.ID)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	id := parseInt(r.FormValue("id"))

	_, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func parseInt(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}
