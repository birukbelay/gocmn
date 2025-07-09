package dtos

import (
	"net/http"
)

// HumaResponse is response To a request
type HumaResponse[T any] struct {
	Body      T `json:"body" doc:"response Body"`
	Status    int
	SetCookie http.Cookie `header:"Set-Cookie"`
}

// HumaReqBody is
type HumaReqBody[T any] struct {
	Body T      `json:"body" doc:"response Body"`
	Auth string `header:"Authorization"`
}

// HumaInputId used for input with id only,  get & delete requests
type HumaInputId struct {
	ID   string `path:"id"`
	Auth string `header:"Authorization"`
}

// HumaReqBodyId used for input with id & body,  patch requests
type HumaReqBodyId[T any] struct {
	ID   string `path:"id"`
	Body T      `json:"body" doc:"request Body"`
	Auth string `header:"Authorization"`
}

type AuthParam struct {
	Auth string `header:"Authorization"`
}
