package lib

import (
	"api/config"
	"strconv"

	"github.com/labstack/echo"
)

// swagger:model Pagination
type Pagination struct {
	Total        *uint64     `json:"total"`
	PerPage      uint64      `json:"per_page"`
	CurrentPage  uint64      `json:"current_page"`
	LastPage     *uint64     `json:"last_page"`
	FirstPageURL *string     `json:"first_page_url"`
	PrevPageURL  *string     `json:"prev_page_url"`
	NextPageURL  *string     `json:"next_page_url"`
	LastPageURL  *string     `json:"last_page_url"`
	From         uint64      `json:"from"`
	To           uint64      `json:"to"`
	Data         interface{} `json:"data,omitempty"`
}

func GenerateMeta(c echo.Context, total uint64, limit uint64, page uint64, offset uint64, pagination bool) Pagination {
	var meta Pagination

	if pagination == true {
		meta.Total = &total
	}
	if limit > 0 {
		meta.PerPage = limit
		if pagination == true {
			lastPage := total / limit
			meta.LastPage = &lastPage
			if (total % limit) > 0 {
				*meta.LastPage++
			}
		}
		if (offset + limit) > total {
			meta.To = total
		} else {
			meta.To = offset + limit
		}
	} else {
		meta.PerPage = total
		if pagination == true {
			*meta.LastPage = 1
		}
		meta.To = total
	}

	if page == 0 {
		meta.CurrentPage = 1
		meta.From = 1
		meta.To = total
	} else {
		meta.CurrentPage = page
	}
	if meta.CurrentPage > 1 {
		prevPageURL := config.Host + c.Path() + "?limit=" + strconv.FormatUint(limit, 10) + "&page=" + strconv.FormatUint(meta.CurrentPage-1, 10)
		meta.PrevPageURL = &prevPageURL
		firstPageURL := config.Host + c.Path() + "?limit=" + strconv.FormatUint(limit, 10) + "&page=1"
		meta.FirstPageURL = &firstPageURL
	}
	if pagination == true {
		if meta.CurrentPage < *meta.LastPage {
			nextPageURL := config.Host + c.Path() + "?limit=" + strconv.FormatUint(limit, 10) + "&page=" + strconv.FormatUint(meta.CurrentPage+1, 10)
			meta.NextPageURL = &nextPageURL
			lastPageURL := config.Host + c.Path() + "?limit=" + strconv.FormatUint(limit, 10) + "&page=" + strconv.FormatUint(*meta.LastPage, 10)
			meta.LastPageURL = &lastPageURL
		}
	} else {
		nextPageURL := config.Host + c.Path() + "?limit=" + strconv.FormatUint(limit, 10) + "&page=" + strconv.FormatUint(meta.CurrentPage+1, 10)
		meta.NextPageURL = &nextPageURL
	}
	meta.From = offset + 1

	return meta
}

func Paginate(c echo.Context, data interface{}, total uint64, limit uint64, page uint64, offset uint64, pagination bool) Pagination {
	meta := GenerateMeta(c, total, limit, page, offset, pagination)
	meta.Data = data
	return meta
}
