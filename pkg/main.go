package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

const (
	API_PATH = "/apis/v1/books"
)

type library struct {
	dbHost, dbPass, dbName string
}
type Book struct {
	Id, Name, Isbn string
}

func main() {
	// Set envoirenment variable for SQL_DB
	dbHost := os.Getenv("DB_HOST") // Set DB_HOST
	if dbHost == "" {
		dbHost = "localhost:8080"
	}
	dbPass := os.Getenv("DB_PASS") // Set DB_PASS
	if dbPass == "" {
		dbPass = "ashishjain"
	}
	apiPath := os.Getenv("API_PATH") // Set API PATH
	if apiPath == "" {
		apiPath = API_PATH
	}
	dbName := os.Getenv("DB_NAME") // Set DB Name
	if dbName == "" {
		dbName = "library"
	}

	l := library{
		dbHost: dbHost,
		dbPass: dbPass,
		dbName: dbName,
	}
	// Go router
	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods(http.MethodGet) // Get Books data
	r.HandleFunc(apiPath, l.postBooks).Methods(http.MethodPost)
	http.ListenAndServe(":8080", r)
}

func (l library) getBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("Get Book api was called")
	db := l.openConnection()                     // Open Connection
	rows, err := db.Query("select * from books") // Read all the Books
	if err != nil {
		log.Fatalf("querying the books found error %s\n", err.Error())
	}
	fmt.Print(rows)
	books := []Book{}
	for rows.Next() {
		var id, isbn, name string
		err := rows.Scan(&id, &name, &isbn)
		if err != nil {
			log.Fatalf("found error %s\n", err.Error())
		}
		aBook := Book{
			Id:   id,
			Name: name,
			Isbn: isbn,
		}
		books = append(books, aBook)
	}
	fmt.Print(books)
	json.NewEncoder(w).Encode(books)
	l.closeConnection(db) // Close Connections
}
func (l library) postBooks(w http.ResponseWriter, r *http.Request) {
	// read the request into an instance of book
	book := Book{}
	json.NewDecoder(r.Body).Decode(&book)
	// open connection
	db := l.openConnection()
	// write the data
	insertQuery, err := db.Prepare("insert into books values (?, ?, ?)")
	if err != nil {
		log.Fatalf("preparing the db query %s\n", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("while begining the transaction %s\n", err.Error())
	}
	_, err = tx.Stmt(insertQuery).Exec(book.Id, book.Name, book.Isbn)
	if err != nil {
		log.Fatalf("execing the insert command %s\n", err.Error())
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("while commint the transaction %s\n", err.Error())
	}
	l.closeConnection(db) // close the connection
}

func (l library) openConnection() *sql.DB {
	// db, err := sql.Open("mysql", "user:password@/dbname"
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", l.dbPass, l.dbHost, l.dbName))
	if err != nil {
		log.Fatalf("Opening the connection to the database %s\n", err.Error())
	}
	return db
}

func (l library) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("Error while closing connection %s\n", err.Error())
	}
}
