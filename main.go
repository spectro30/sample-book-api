package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// ---------------Authentication--------------

var jwtkey = []byte("chandlerbing")

var admins = map[string]string{
	"alex" : "123456",
	"ben"  : "654321",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}



func Login(w http.ResponseWriter, r *http.Request){
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	correctpassword, ok := admins[creds.Username]
	if !ok || correctpassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expiretime := time.Now().Add(15 * time.Minute)

	claims := jwt.StandardClaims{
		Id: creds.Username,
		ExpiresAt:  expiretime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name : "token",
		Value: tokenString,
		Expires: expiretime,
	})
}

func Accesscheck(w http.ResponseWriter, r *http.Request) bool{
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
	}

	tokenstring := c.Value
	claims := &jwt.StandardClaims{}
	tkn, err := jwt.ParseWithClaims(tokenstring, claims, func(token *jwt.Token)(interface{}, error){
		return jwtkey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid{
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	return true
}

// --------------------------------



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

func Assignid () int {
	if len(unusedid) == 0 {
		return len(books) + 1
	}
	var res = unusedid[0]
	unusedid = unusedid[1:]
	return res
}

func Appendbooks () {
	books = append(books, Book{ID: "1", Title: "The Kite Runner",
		Author: Author{Firstname:"Khaled", Lastname: "Hosseini", AuthorID: "40"}, Genre: "Drama" })
	books = append(books, Book{ID: "2", Title: "Inception Point",
		Author: Author{Firstname:"Dan", Lastname: "Brown", AuthorID: "53"}, Genre: "Thriller" })
	books = append(books, Book{ID: "3", Title: "Lost Symbol",
		Author: Author{Firstname:"Dan", Lastname: "Brown", AuthorID: "53"}, Genre: "Thriller" })
}

func Getallbooks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Context-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func Getbookbyid(w http.ResponseWriter, r *http.Request){
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

func Getbookbyauthorid(w http.ResponseWriter, r *http.Request){
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

func Getbookbygenre(w http.ResponseWriter, r *http.Request){
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

func Addbook(w http.ResponseWriter, r *http.Request){
	if !Accesscheck(w, r){
		return
	}
	w.Header().Set("Context-Type", "application/json")
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	book.ID = strconv.Itoa(Assignid())
	books = append(books, book)

	json.NewEncoder(w).Encode(book)
}

func Updatebook(w http.ResponseWriter, r *http.Request){
	if !Accesscheck(w, r){
		return
	}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books[index] = books[len(books)-1]
			books = books[:len(books)-1]
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

func Deletebook(w http.ResponseWriter, r *http.Request) {
	if !Accesscheck(w, r){
		return
	}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books[index] = books[len(books)-1]
			books = books[:len(books)-1]
			break
		}
	}
	json.NewEncoder(w).Encode(books)
	intid, _ := strconv.ParseInt(params["id"], 10, 64)
	unusedid = append(unusedid, int(intid))
}


func main(){
	r := mux.NewRouter()
	Appendbooks()

	r.HandleFunc("/login", Login)
	r.HandleFunc("/books", Getallbooks).Methods("GET")
	r.HandleFunc("/books/bookid/{id}", Getbookbyid).Methods("GET")
	r.HandleFunc("/books/authorid/{authorid}", Getbookbyauthorid).Methods("GET")
	r.HandleFunc("/books/genre/{genre}", Getbookbygenre).Methods("GET")
	r.HandleFunc("/books", Addbook).Methods("POST")
	r.HandleFunc("/books/{id}", Updatebook).Methods("PUT")
	r.HandleFunc("/books/{id}", Deletebook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9003", r))

}


/*

{
    "id": "11",
    "title": "HARRY POTTER AND THE PHILOSOPHER’S STONE",
    "author": {
        "firstname": "J.K.",
        "lastname": "Rowling",
        "authorid": "25"
    },
    "genre": "Fantasy"
}
{
    "id": "12",
    "title": "Harry Potter and the Prisoner of Azkaban",
    "author": {
        "firstname": "J.K.",
        "lastname": "Rowling",
        "authorid": "25"
    },
    "genre": "Fantasy"
}

{
    "username": "ben",
    "password": "654321"
}

{
    "username": "alex",
    "password": "123456"

}


 */