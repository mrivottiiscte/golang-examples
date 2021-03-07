package product

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	uuid "github.com/satori/go.uuid"
)

type repository struct {
	db *sql.DB
}

type product struct {
	ID        uuid.UUID
	Name      string
	Units     int
	Weight    float32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (repo repository) saveProduct(name string, units int, weight float32) (*product, error) {
	sqlStatement := `
	INSERT INTO products (name, units, weight)
	VALUES ($1, $2, $3)
	RETURNING id, name, units, weight, created_at, updated_at`

	p := product{}
	err := repo.db.QueryRow(sqlStatement, name, units, weight).Scan(&p.ID, &p.Name, &p.Units, &p.Weight, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &p, nil
}

func (repo repository) getProductByID(id uuid.UUID) (*product, error) {
	sqlStatement := `
	SELECT id, name, units, weight, created_at, updated_at
	FROM products
	WHERE id = $1`

	p := product{}
	err := repo.db.QueryRow(sqlStatement, id).Scan(&p.ID, &p.Name, &p.Units, &p.Weight, &p.CreatedAt, &p.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		log.Printf("get ok")
	}

	return &p, nil
}

func (repo repository) getProductByName(name string) (*product, error) {
	sqlStatement := `
	SELECT id, name, units, weight, created_at, updated_at
	FROM products
	WHERE name = $1`

	p := product{}
	err := repo.db.QueryRow(sqlStatement, name).Scan(&p.ID, &p.Name, &p.Units, &p.Weight, &p.CreatedAt, &p.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	default:
		log.Printf("get ok")
	}

	return &p, nil
}

func (repo repository) listProducts(limit int, offset int) ([]product, error) {
	sqlStatement := `
	SELECT id, name, units, weight, created_at, updated_at
	FROM products
	LIMIT $1 OFFSET $2`

	rows, err := repo.db.Query(sqlStatement, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]product, 0)
	for rows.Next() {
		p := product{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Units, &p.Weight, &p.CreatedAt, &p.UpdatedAt); err != nil {
			log.Fatal(err)
		}
		products = append(products, p)
	}

	return products, nil
}
