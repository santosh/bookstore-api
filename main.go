package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goji.io"
	"goji.io/pat"
)

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

type Book struct {
	ISBN    string   `json:"isbn" bson:"isbn,omitempty"`
	Title   string   `json:"title" bson:"title,omitempty"`
	Authors []string `json:"authors" bson:"authors,omitempty"`
	Price   string   `json:"price" bson:"price,omitempty"`
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://172.19.0.2"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/books"), allBooks(client))
	mux.HandleFunc(pat.Post("/books"), addBook(client))
	mux.HandleFunc(pat.Get("/books/:isbn"), bookByISBN(client))
	mux.HandleFunc(pat.Put("/books/:isbn"), updateBook(client))
	mux.HandleFunc(pat.Delete("/books/:isbn"), deleteBook(client))
	http.ListenAndServe("localhost:8080", mux)
}

func allBooks(c *mongo.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		coll := c.Database("store").Collection("books")

		cursor, err := coll.Find(context.TODO(), bson.M{})
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all books: ", err)
			return
		}

		var books []Book
		err = cursor.All(context.TODO(), &books)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all books: ", err)
			return
		}

		respBody, err := json.MarshalIndent(books, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func addBook(c *mongo.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var book Book
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&book)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		coll := c.Database("store").Collection("books")

		_, err = coll.InsertOne(context.TODO(), book)
		if err != nil {
			// FIXME: Check for duplicate document

			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed insert book: ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", r.URL.Path+"/"+book.ISBN)
		w.WriteHeader(http.StatusCreated)
	}
}

func bookByISBN(c *mongo.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		isbn := pat.Param(r, "isbn")

		coll := c.Database("store").Collection("books")

		var book Book
		err := coll.FindOne(context.TODO(), bson.M{"isbn": isbn}).Decode(&book)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed find book: ", err)
			return
		}

		if book.ISBN == "" {
			ErrorWithJSON(w, "Book not found", http.StatusNotFound)
			return
		}

		respBody, err := json.MarshalIndent(book, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func updateBook(c *mongo.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		isbn := pat.Param(r, "isbn")

		var book Book
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&book)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}

		coll := c.Database("store").Collection("books")
		filter := bson.M{"isbn": isbn}
		update := bson.D{{"$set", book}}

		_, err = coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed update book: ", err)
				return
			}
		}
	}
}

func deleteBook(c *mongo.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		isbn := pat.Param(r, "isbn")

		coll := c.Database("store").Collection("books")
		filter := bson.M{"isbn": isbn}

		_, err := coll.DeleteOne(context.TODO(), filter)
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
				log.Println("Failed delete book: ", err)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
