package dtos

import (
	cmn "github.com/birukbelay/gocmn/src/logger"
)

type Resp[T any] struct {
	Body    T `json:"body" doc:"response Body"`
	Status  int
	Message string
	Ok      bool
	Error   error
	HasResp bool
}

func SUCCEED[T any](body T, message string) *Resp[T] {
	cmn.LogTrace("succeed", message)

	return &Resp[T]{
		Message: message,
		Body:    body,
		Ok:      true,
		HasResp: true,
	}
}

func FAIL[T any](data T, err error, message string) *Resp[T] {
	cmn.LogTraceN("error happned", message, 3)

	return &Resp[T]{
		Message: message,
		Ok:      false,
		Error:   err,
		HasResp: true,
		Body:    data,
	}
}
