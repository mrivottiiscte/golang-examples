package product

import (
	"database/sql"
	"go-example/echo-db/internal/apierrors"
	"go-example/echo-db/pagination"
	"net/http"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
)

type handler struct {
	service service
}

func NewHandler(db *sql.DB) handler {
	return handler{service: NewService(db)}
}

type productRequest struct {
	Name   string  `json:"name" validate:"required,max=50"`
	Units  int     `json:"units" validate:"required,min=1"`
	Weight float32 `json:"weight" validate:"required"`
}

type productResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Units  int       `json:"units"`
	Weight float32   `json:"weight"`
}

func toProductResponse(p *product) productResponse {
	return productResponse{
		ID:     p.ID,
		Name:   p.Name,
		Units:  p.Units,
		Weight: p.Weight,
	}
}

func (h handler) Post(c echo.Context) error {
	req := new(productRequest)

	err := c.Bind(req)
	if err != nil {
		return err
	}
	/*
		err = c.Validate(req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, apierrors.New(err.Error()))
		}
	*/
	p, err := h.service.getProductByName(req.Name)
	if err != nil {
		return err
	}
	if p != nil {
		return c.JSON(http.StatusBadRequest, apierrors.New("name already exists"))
	}

	p, err = h.service.saveProduct(req)
	if err != nil {
		return err
	}

	resp := toProductResponse(p)

	return c.JSON(http.StatusOK, resp)
}

func (h handler) Get(c echo.Context) error {

	id, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	p, err := h.service.getProductByID(id)
	if err != nil {
		return err
	}

	if p == nil {
		return c.NoContent(http.StatusNotFound)
	}

	resp := toProductResponse(p)

	return c.JSON(http.StatusOK, resp)
}

func (h handler) List(c echo.Context) error {

	page, pageSize := pagination.PagePageSize(c)

	products, err := h.service.listProducts(page, pageSize)
	if err != nil {
		return err
	}

	productsResponse := make([]productResponse, 0)
	for _, p := range products {
		productsResponse = append(productsResponse, toProductResponse(&p))
	}

	return c.JSON(http.StatusOK, productsResponse)
}

func (h handler) Health(c echo.Context) error {

	return c.JSON(http.StatusOK, "ok")
}
