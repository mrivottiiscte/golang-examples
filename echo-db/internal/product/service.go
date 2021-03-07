package product

import (
	"database/sql"
	"go-example/echo-db/pagination"

	uuid "github.com/satori/go.uuid"
)

type service struct {
	repo repository
}

func NewService(db *sql.DB) service {
	return service{repo: repository{db: db}}
}

func (s service) saveProduct(req *productRequest) (*product, error) {

	//p, err := s.repo.getProductByName()

	p, err := s.repo.saveProduct(req.Name, req.Units, req.Weight)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s service) getProductByID(id uuid.UUID) (*product, error) {

	p, err := s.repo.getProductByID(id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s service) getProductByName(name string) (*product, error) {

	p, err := s.repo.getProductByName(name)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s service) listProducts(page int, pageSize int) ([]product, error) {

	limit, offset := pagination.LimitOffset(page, pageSize)

	products, err := s.repo.listProducts(limit, offset)
	if err != nil {
		return nil, err
	}

	return products, nil
}
