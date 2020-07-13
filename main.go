package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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
var unusedid [] int

func assignid () int {
	if len(unusedid) == 0 {
		return len(books) + 1
	}
	var res = unusedid[0]
	unusedid = unusedid[1:len(unusedid)]
	return res
}

func appendbooks () {
	books = append(books, Book{ID: "1", Title: "The Kite Runner",
		Author: Author{Firstname:"Khaled", Lastname: "Hosseini", AuthorID: "40"}, Genre: "Drama" })
	books = append(books, Book{ID: "2", Title: "Inception Point",
		Author: Author{Firstname:"Dan", Lastname: "Brown", AuthorID: "53"}, Genre: "Thriller" })
	books = append(books, Book{ID: "3", Title: "Lost Symbol",
		Author: Author{Firstname:"Dan", Lastname: "Brown", AuthorID: "53"}, Genre: "Thriller" })
}

func getallbooks(w http.ResponseWriter, r *http.Request){
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

func getbookbygenre(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var bookbygen []Book
	for _, item := range books{
		if item.Genre == params["genre"]{
			bookbygen = append(bookbygen, item)
		}
	}
	json.NewEncoder(w).Encode(bookbygen)
}

func addbook(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Context-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(assignid())
	//fmt.Println(book.ID)
	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

func updatebook(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

func deletebook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
	intid, _ := strconv.ParseInt(params["id"], 10, 64)
	unusedid = append(unusedid, int(intid))
}


func main(){
	r := mux.NewRouter()
	appendbooks()

	r.HandleFunc("/books", getallbooks).Methods("GET")
	r.HandleFunc("/books/bookid/{id}", getbookbyid).Methods("GET")
	r.HandleFunc("/books/authorid/{authorid}", getbookbyauthorid).Methods("GET")
	r.HandleFunc("/books/genre/{genre}", getbookbygenre).Methods("GET")
	r.HandleFunc("/books", addbook).Methods("POST")
	r.HandleFunc("/books/{id}", updatebook).Methods("PUT")
	r.HandleFunc("/books/{id}", deletebook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9003", r))

}


/*
addbook body:

{
    "id": "104",
    "title": "HARRY POTTER AND THE PHILOSOPHER’S STONE",
    "author": {
        "firstname": "J.K.",
        "lastname": "Rowling",
        "authorid": "25"
    },
    "genre": "Fiction"
}


 */