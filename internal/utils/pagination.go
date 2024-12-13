package utils

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultSize = 10
)

// Pagination List Type Response Model
// @Description Pagination List Type Response Model
type PaginationResponse[T any] struct {
	TotalCount int  `json:"total_count"`
	TotalPages int  `json:"total_pages"`
	Page       int  `json:"page"`
	Size       int  `json:"size"`
	HasMore    bool `json:"has_more"`
	Values     []T  `json:"values"`
}

func PaginatedResponse[T any](count int, pq *PaginationQuery, list []T) *PaginationResponse[T] {
	return &PaginationResponse[T]{
		TotalCount: count,
		TotalPages: GetTotalPages(count, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    GetHasMore(pq.GetPage(), count, pq.GetSize()),
		Values:     list,
	}
}

func DefaultPaginationResponse[T any](pq *PaginationQuery) *PaginationResponse[T] {
	return &PaginationResponse[T]{
		TotalCount: 0,
		TotalPages: GetTotalPages(0, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    GetHasMore(pq.GetPage(), 0, pq.GetSize()),
		Values:     make([]T, 0),
	}
}

// Pagination query params
type PaginationQuery struct {
	Size     int    `json:"size,omitempty"`
	Page     int    `json:"page,omitempty"`
	OrderBy  string `json:"orderBy,omitempty"`
	OrderDir string `json:"orderDir,omitempty"`
}

// Set page size
func (q *PaginationQuery) SetSize(sizeQuery string) error {
	if sizeQuery == "" {
		q.Size = defaultSize
		return nil
	}
	n, err := strconv.Atoi(sizeQuery)
	if err != nil {
		return err
	}
	q.Size = n

	return nil
}

// Set page number
func (q *PaginationQuery) SetPage(pageQuery string) error {
	if pageQuery == "" {
		q.Size = 0
		return nil
	}
	n, err := strconv.Atoi(pageQuery)
	if err != nil {
		return err
	}
	q.Page = n

	return nil
}

// Set order by
func (q *PaginationQuery) SetOrderBy(orderByQuery string, orderDirQuery string) {
	if orderByQuery == "" {
		q.OrderBy = "id"
	} else {
		q.OrderBy = orderByQuery
	}

	if orderDirQuery == "" || orderDirQuery == strings.ToLower("asc") {
		q.OrderBy = fmt.Sprintf("%s %s", q.OrderBy, "ASC")
	} else {
		q.OrderBy = fmt.Sprintf("%s %s", q.OrderBy, "DESC")
	}
}

// Get offset
func (q *PaginationQuery) GetOffset() int {
	if q.Page == 0 {
		return 0
	}
	return (q.Page - 1) * q.Size
}

// Get limit
func (q *PaginationQuery) GetLimit() int {
	return q.Size
}

// Get OrderBy
func (q *PaginationQuery) GetOrderBy() string {
	return q.OrderBy
}

// Get OrderBy
func (q *PaginationQuery) GetPage() int {
	return q.Page
}

// Get OrderBy
func (q *PaginationQuery) GetSize() int {
	return q.Size
}

func (q *PaginationQuery) GetQueryString() string {
	return fmt.Sprintf("page=%v&size=%v&orderBy=%s", q.GetPage(), q.GetSize(), q.GetOrderBy())
}

// Get pagination query struct from
func GetPaginationFromRequest(r *http.Request) (*PaginationQuery, error) {
	q := &PaginationQuery{}
	if err := q.SetPage(r.URL.Query().Get("page")); err != nil {
		return nil, err
	}
	if err := q.SetSize(r.URL.Query().Get("size")); err != nil {
		return nil, err
	}
	q.SetOrderBy(r.URL.Query().Get("orderBy"), r.URL.Query().Get("orderDir"))

	return q, nil
}

// Get total pages int
func GetTotalPages(totalCount int, pageSize int) int {
	d := float64(totalCount) / float64(pageSize)
	return int(math.Ceil(d))
}

// Get has more
func GetHasMore(currentPage int, totalCount int, pageSize int) bool {
	return currentPage < totalCount/pageSize
}
