package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestLogin(t *testing.T) {
	testData := []struct{
		name      string
		pass      string
		errorText string
		status 	  int
	}{
		{name: "alex", pass: "123456", errorText: "200 StatusOK", status: 200},
		{name: "ben", pass: "654321", errorText: "200 StatusOK", status: 200},
		{name: "pumba", pass: "111111", errorText: "401 Unauthorized", status: 401},
		{name: "simba", pass: "11111", errorText: "401 Unauthorized", status: 401},
		{name: "timon", pass: "1111", errorText: "401 Unauthorized", status: 401},
	}
	for _, item := range testData {
		cred := Credentials{Username: item.name, Password: item.pass}
		b, _ := json.Marshal(cred)
		iorData := bytes.NewReader(b)
		req, err := http.NewRequest("POST", "localhost:9003/login", iorData)
		if err != nil {
			t.Fatalf("Could not create Request: %v", err)
		}
		rr := httptest.NewRecorder()
		Login(rr, req)
		if res := rr.Result(); res.StatusCode != item.status {
			t.Errorf("expected %v: got %v", item.errorText ,res.Status)
		}
	}
}


func TestGetallbooks(t *testing.T) {
	req, err := http.NewRequest("GET", "/books", nil)
	req.Header.Set("","")
	if err != nil {
		t.Fatalf("Could not create Request: %v", err)
	}
	rr := httptest.NewRecorder()
	Getallbooks(rr, req)
	if res := rr.Result() ; res.StatusCode != http.StatusOK{
		t.Errorf("expected 200 StatusOK: got %v" ,res.Status)
	}
}

func TestGetbookbyauthorid(t *testing.T) {
	testData := []struct{
		qry      string
		errorText string
		status 	  int
	}{
		{qry: "53",  errorText: "200 StatusOK", status: 200},
		{qry: "40",  errorText: "200 StatusOK", status: 200},
		{qry: "5",  errorText: "200 StatusOK", status: 200},
		{qry: "1", errorText: "200 StatusOK", status: 200},
	}
	for _, item := range testData {
		req, err := http.NewRequest("GET", "localhost:9003/books/authorid/"+item.qry, nil)
		if err != nil {
			t.Fatalf("Could not create Request: %v", err)
		}
		rr := httptest.NewRecorder()
		Getbookbyauthorid(rr, req)
		if res := rr.Result(); res.StatusCode != item.status {
			t.Errorf("expected %v: got %v", item.errorText ,res.Status)
		}
	}
}

func TestGetbookbygenre(t *testing.T) {
	testData := []struct{
		qry      string
		errorText string
		status 	  int
	}{
		{qry: "Fantasy",  errorText: "200 StatusOK", status: 200},
		{qry: "Genre",  errorText: "200 StatusOK", status: 200},
		{qry: "Thriller",  errorText: "200 StatusOK", status: 200},
		{qry: "Drama", errorText: "200 StatusOK", status: 200},
	}
	for _, item := range testData {
		req, err := http.NewRequest("GET", "localhost:9003/books/genre/"+item.qry, nil)
		if err != nil {
			t.Fatalf("Could not create Request: %v", err)
		}
		rr := httptest.NewRecorder()
		Getbookbyauthorid(rr, req)
		if res := rr.Result(); res.StatusCode != item.status {
			t.Errorf("expected %v: got %v", item.errorText ,res.Status)
		}
	}
}

func TestGetbookbyid(t *testing.T) {
	testData := []struct{
		qry      string
		errorText string
		status 	  int
	}{
		{qry: "1",  errorText: "200 StatusOK", status: 200},
		{qry: "2",  errorText: "200 StatusOK", status: 200},
		{qry: "3",  errorText: "200 StatusOK", status: 200},
		{qry: "4", errorText: "200 StatusOK", status: 200},
	}
	for _, item := range testData {
		req, err := http.NewRequest("GET", "localhost:9003/books/bookid/"+item.qry, nil)
		if err != nil {
			t.Fatalf("Could not create Request: %v", err)
		}
		rr := httptest.NewRecorder()
		Getbookbyauthorid(rr, req)
		if res := rr.Result(); res.StatusCode != item.status {
			t.Errorf("expected %v: got %v", item.errorText ,res.Status)
		}
	}
}

func TestAddbook(t *testing.T) {
	var tt []Book
	tt = append(tt, Book{ID: "1", Title: "The Kite Runner",
		Author: Author{Firstname:"Khaled", Lastname: "Hosseini", AuthorID: "40"}, Genre: "Drama" })
	tt = append(tt, Book{ID: "2", Title: "Inception Point",
		Author: Author{Firstname:"Dan", Lastname: "Brown", AuthorID: "53"}, Genre: "Thriller" })
	tt = append(tt, Book{ID: "3", Title: "Lost Symbol",
		Author: Author{Firstname:"Dan", Lastname: "Brown", AuthorID: "53"}, Genre: "Thriller" })

	for _, item := range tt {
		b, _ := json.Marshal(item)
		iorData := bytes.NewReader(b)
		req, err := http.NewRequest("POST", "localhost:9003/books", iorData)
		if err != nil {
			t.Fatalf("Could not create Request: %v", err)
		}
		rr := httptest.NewRecorder()
		Addbook(rr, req)
		if res := rr.Result(); res.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 200 StatusOK: got %v",res.Status)
		}
	}
}


func TestDeletebook(t *testing.T) {
	testData := []struct{
		qry      string
		errorText string
		status 	  int
	}{
		{qry: "1",  errorText: "200 StatusOK", status: 401},
		{qry: "2",  errorText: "200 StatusOK", status: 401},
		{qry: "3",  errorText: "200 StatusOK", status: 401},
		{qry: "4", errorText: "200 StatusOK", status: 401},
	}
	for _, item := range testData {
		req, err := http.NewRequest("DELETE", "localhost:9003/books/"+item.qry, nil)
		if err != nil {
			t.Fatalf("Could not create Request: %v", err)
		}
		rr := httptest.NewRecorder()
		Deletebook(rr, req)
		if res := rr.Result(); res.StatusCode != item.status {
			t.Errorf("expected %v: got %v", item.errorText ,res.Status)
		}
	}
}









