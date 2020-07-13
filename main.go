package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)


type Book struct{
	ID string `json:"id"`
	Title string `json:"title"`
	Author Author `json:"author"`
	Genre string `json:"genre"`
}

type Author struct{
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	AuthorID string `json:"authorid"`
}

var books []Book
func appendbooks () {
	books = append(books, Book{ID: "101", Title: "The Kite Runner",
		Author: Author{Firstname:"Khaled", Lastname: "Hosseini", AuthorID: "40"}, Genre: "Drama" })
	books = append(books, Book{ID: "102", Title: "Inception Point",
		Author: Author{Firstname:"Dan", Lastname: "Brown", AuthorID: "53"}, Genre: "Thriller" })
	books = append(books, Book{ID: "103", Title: "Lost Symbol",
		Author: Author{Firstname:"Dan", Lastname: "Brown", AuthorID: "53"}, Genre: "Thriller" })
}

func getbooks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Context-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getbookbyid(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books{
		if item.ID == params["id"]{
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func getbookbyauthorid(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var bookbyauth []Book
	for _, item := range books{
		if item.Author.AuthorID == params["authorid"]{
			bookbyauth = append(bookbyauth, item)
		}
	}
	json.NewEncoder(w).Encode(bookbyauth)
}



func main(){
	r := mux.NewRouter()
	appendbooks()
	r.HandleFunc("/books", getbooks).Methods("GET")
	r.HandleFunc("/books/id/{id}", getbookbyid).Methods("GET")
	r.HandleFunc("/books/author/{authorid}", getbookbyauthorid).Methods("GET")

	log.Fatal(http.ListenAndServe(":9002", r))

}