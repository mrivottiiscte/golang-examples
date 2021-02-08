package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type Product struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Price float32 `json:"price"`
}

var products []Product

func getProduct(c echo.Context) error {
	id := c.Param("id")

	for _, p := range products {
		if p.ID == id {
			return c.JSON(http.StatusOK, p)
		}
	}

	return c.NoContent(http.StatusNotFound)
}

func listProducts(c echo.Context) error {
	return c.JSON(http.StatusOK, products)
}

func createProduct(c echo.Context) error {
	p := new(Product)
	if err := c.Bind(p); err != nil {
		return err
	}

	p.ID = uuid.NewV4().String()
	products = append(products, *p)

	return c.JSON(http.StatusCreated, p)
}

func main() {

	p1 := Product{
		ID:    "054e5ebd-62c9-4e02-8c19-d96b01ab7aee",
		Title: "Tenis azuis",
		Price: 10.0,
	}

	p2 := Product{
		ID:    "054e5ebd-62c9-4e02-8c19-d96b01ab7aee",
		Title: "T-shirt amarela",
		Price: 20.0,
	}

	products = make([]Product, 0)
	products = append(products, p1, p2)

	e := echo.New()
	e.GET("/products/:id", getProduct)
	e.GET("/products", listProducts)
	e.POST("/products", createProduct)
	e.Logger.Fatal(e.Start(":8080"))
}
