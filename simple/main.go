package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Product struct {
	ID    int
	Title string
	Price float32
}

var products []Product

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello turma!!!!!</h1>")
}

func listProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func main() {

	p1 := Product{
		ID:    1,
		Title: "Tenis azuis",
		Price: 10.0,
	}

	p2 := Product{
		ID:    2,
		Title: "T-shirt amarela",
		Price: 20.0,
	}

	products = make([]Product, 0)
	products = append(products, p1, p2)

	http.HandleFunc("/", handler)
	http.HandleFunc("/products", listProducts)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
