package dtos

import (
	"github.com/birukbelay/gocmn/src/resp_const"
)

type PResp[T any] struct {
	Status       int   `json:"status,omitempty"`
	RowsAffected int64 `json:"rows_affected"`
	Count        int64 `json:"count,omitempty"`
	HasMore      bool  `json:"has_more"`
	HasPrev      bool  `json:"has_prev"`

	//messages
	Code    resp_const.RespCode `json:"code,omitempty"`
	Message string              `json:"message,omitempty"`
	Error   string              `json:"error,omitempty"`
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

	Code         resp_const.RespCode `json:"code,omitempty"`
	Message      string              `json:"message,omitempty"`
	Error        string              `json:"error,omitempty"`
	RowsAffected int64               `json:"rows_affected"`
	NextCursor   string              `json:"next_cursor,omitempty"`
	PrevCursor   string              `json:"prev_cursor,omitempty"`
	Body         T                   `json:"body"`
	Ok           bool                `json:"ok,omitempty"`
}
type CursorStruct struct {
	NextCursor string `json:"next_cursor"`
	PrevCursor string `json:"prev_cursor"`
	HasPrev    bool   `json:"has_prev"`
	HasMore    bool   `json:"has_more"`
}

type PaginationInput struct {
	Limit  int    `default:"25" query:"_limit"`
	Page   int    `default:"1"  query:"_page"`
	Cursor string ` query:"_cursor"`
	//
	SortDir     string   ` query:"_sort_dir" enum:"asc,desc" default:"desc"`
	AllowedSort []string `query:"-"  `

	//fields that are overridden on api side, by default
	Select []string `query:"-"`
	SortBy string   `query:"-" ` //we dont use this because we override it with the models Sort

	//fields that are not allowed on api side, but used internally
	PrefixColLike string   `query:"-" doc:"used for prefix like query"`
	TxtSearchCols []string `query:"-"  doc:"used for text search"`
	//fields that need to be set on the api side to be used
	Like  string `query:"-"` //set `ColName`` on api side to use this
	Query string `query:"-"` //set `ColNames` on api side to use this
	//fields that need to be set on the api side to be used
	StartDate string `query:"-"`
	EndDate   string `query:"-"` //todo, these should also be allowed on api side

	//tag fields
	Tags    []string `query:"-"`
	AllTags []string `query:"all_tags"`
	AnyTags []string `query:"any_tags"`
}

type LocResp[T any] struct {
	Status  int                 `json:"status"`
	Code    resp_const.RespCode `json:"code"`
	Message string              `json:"message"`
	Ok      bool
	Val     T
}
