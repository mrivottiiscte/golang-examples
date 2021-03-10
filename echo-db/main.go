package main

import (
	"database/sql"
	"fmt"
	"go-example/echo-db/internal/product"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

func main() {

	host := os.Getenv("DATABASE_HOST")
	if host == "" {
		panic("no DATABASE_HOST")
	}

	port := os.Getenv("DATABASE_PORT")
	if host == "" {
		panic("no DATABASE_HOST")
	}

	user := os.Getenv("DATABASE_USER")
	if user == "" {
		panic("no DATABASE_HOST")
	}

	password := os.Getenv("DATABASE_PASSWORD")
	if password == "" {
		panic("no DATABASE_HOST")
	}

	dbname := os.Getenv("DATABASE_NAME")
	if dbname == "" {
		panic("no DATABASE_HOST")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS products (
		id uuid DEFAULT uuid_generate_v4 (),
		name VARCHAR NOT NULL,
		units INT NOT NULL,
		weight FLOAT NOT NULL,
		created_at TIME DEFAULT NOW(),
		updated_at TIME DEFAULT NOW(),
		PRIMARY KEY (id)
	);`)
	if err != nil {
		panic(err)
	}

	productHandler := product.NewHandler(db)

	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/product", productHandler.Post)
	e.GET("/product/:id", productHandler.Get)
	e.GET("/product/", productHandler.List)

	e.GET("/health", productHandler.Health)

	e.Logger.Fatal(e.Start(":8080"))
}
