package dtos

import (
	response_const "github.com/birukbelay/gocmn/resp_codes"
	"net/http"
)

func ReturnOffsetPaginatedS[T any](item T, count int64, hasMore bool, rowsAffected int64) PResp[T] {
	return PResp[T]{
		Status:       http.StatusOK,
		Code:         response_const.Success,
		Message:      response_const.Success.ToStr(),
		Body:         item,
		Count:        count,
		HasMore:      hasMore,
		RowsAffected: rowsAffected,
	}
}
func ReturnCursorPaginatedS[T any](item T, cursor CursorStruct, count int64, rowsAffected int64) PResp[T] {
	return PResp[T]{
		Status:       http.StatusOK,
		Code:         response_const.Success,
		Message:      response_const.Success.ToStr(),
		Body:         item,
		Count:        count,
		NextCursor:   cursor.NextCursor,
		PrevCursor:   cursor.PrevCursor,
		HasMore:      cursor.HasMore,
		HasPrev:      cursor.HasPrev,
		RowsAffected: rowsAffected,
	}
}

// PInternalErr is internal error for paginated responses
func PInternalErrS[T any](message string) PResp[T] {
	return PResp[T]{
		Status: http.StatusInternalServerError,
		Code:   response_const.FAIL,
		Error:  message,
		//Body:         item,
	}
}

// ====================================   Non Paginated messages =======================
func FetchSuccessAS[T any](item T, message string, rowsAffected int64) GResp[T] {
	return GResp[T]{
		Status:       http.StatusOK,
		Code:         response_const.Success,
		Message:      message,
		Body:         item,
		RowsAffected: rowsAffected,
	}
}

func NotFoundErrS[T any](message string) GResp[T] {
	return GResp[T]{
		Status: http.StatusNotFound,
		Code:   response_const.FAIL,
		Error:  message,
	}
}

// ===============  Bad Request ======
func BadReqErrS[T any](message string) GResp[T] {
	return GResp[T]{
		Status: http.StatusBadRequest,
		Code:   response_const.FAIL,
		Error:  message,
		Ok:     false,
	}
}
func BadRequestMsgS[T any](errorMessage string) GResp[T] {
	return GResp[T]{
		Status: http.StatusBadRequest,
		Code:   response_const.BadRequest,
		Error:  errorMessage,
		Ok:     false,
	}
}

func FetchSuccessWithCountS[T any](item T, message string, rowsAffected, count int64) GResp[T] {
	return GResp[T]{
		Status:       http.StatusOK,
		Code:         response_const.Success,
		Message:      message,
		Body:         item,
		RowsAffected: rowsAffected,
		Count:        count,
	}

}

//====================================   Non Paginated messages =======================

// SuccessCS used for  getOne, create, update & delete success
func SuccessCS[T any](item T, code response_const.RespCode, rowsAffected int64) GResp[T] {
	return GResp[T]{
		Status:       http.StatusOK,
		Code:         code,
		Message:      code.Msg(),
		Body:         item,
		RowsAffected: rowsAffected,
	}
}
func SuccessS[T any](item T, rowsAffected int64) GResp[T] {
	return GResp[T]{
		Status:       http.StatusCreated,
		Code:         response_const.Success,
		Message:      "success",
		Body:         item,
		RowsAffected: rowsAffected,
	}
}

func InternalErrS[T any](message string) GResp[T] {
	return GResp[T]{
		Status: http.StatusInternalServerError,
		Code:   response_const.FAIL,
		Error:  message,
		//Item:         item,
	}
}
