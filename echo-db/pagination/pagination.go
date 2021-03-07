package pagination

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

const defaultPage = 0
const maxPageSize = 50

func LimitOffset(page int, pageSize int) (int, int) {

	limit := pageSize
	offset := page * pageSize

	return limit, offset
}

func PagePageSize(c echo.Context) (int, int) {

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = defaultPage
	}
	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil {
		pageSize = maxPageSize
	}

	if page < 0 {
		page = defaultPage
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return page, pageSize
}
