package dtos

import (
	response_const "github.com/birukbelay/gocmn/src/resp_codes"
)

type PResp[T any] struct {
	Status       int   `json:"status,omitempty"`
	RowsAffected int64 `json:"rows_affected"`
	Count        int64 `json:"count,omitempty"`
	HasMore      bool  `json:"has_more"`
	HasPrev      bool  `json:"has_prev"`

	//messages
	Code    response_const.RespCode `json:"code,omitempty"`
	Message string                  `json:"message,omitempty"`
	Error   string                  `json:"error,omitempty"`
	//Cursors
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
	Body       T      `json:"body"`
}
type GResp[T any] struct {
	Status  int   `json:"status,omitempty"`
	Count   int64 `json:"count"`
	HasMore bool  `json:"has_more"`
	HasPrev bool  `json:"has_prev"`

	Code         response_const.RespCode `json:"code,omitempty"`
	Message      string                  `json:"message,omitempty"`
	Error        string                  `json:"error,omitempty"`
	RowsAffected int64                   `json:"rows_affected"`
	NextCursor   string                  `json:"next_cursor,omitempty"`
	PrevCursor   string                  `json:"prev_cursor,omitempty"`
	Body         T                       `json:"body"`
	Ok           bool                    `json:"ok,omitempty"`
}
type CursorStruct struct {
	NextCursor string `json:"next_cursor"`
	PrevCursor string `json:"prev_cursor"`
	HasPrev    bool   `json:"has_prev"`
	HasMore    bool   `json:"has_more"`
}

type PaginationInput struct {
	Limit  int    `default:"25" form:"limit,omitempty" query:"limit"`
	Cursor string `form:"cursor" query:"cursor"`
	Page   int    `default:"1" form:"page" query:"page"`
	//sorting
	SortDir     string   `form:"dir" query:"sort_dir" enum:"asc,desc" default:"desc"`
	SortBy      string   `form:"ord" query:"-" ` //we dont use this because we override it with the models Sort
	AllowedSort []string `form:"-"  `
	Select      []string ``
	Tags        []string `query:"tags"`
	Query       string   `query:"q"`
	Like        string   `query:"like"`
	ColName     string   `query:"-"`
	ColNames    []string `query:"-"`
	StartDate   string   `query:"start_date"`
	EndDate     string   `query:"end_date"`
}

type LocResp[T any] struct {
	Status  int                     `json:"status"`
	Code    response_const.RespCode `json:"code"`
	Message string                  `json:"message"`
	Ok      bool
	Val     T
}
