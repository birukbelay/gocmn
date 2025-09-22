package dtos

import (
	"net/http"

	"github.com/birukbelay/gocmn/src/resp_const"
)

func ReturnOffsetPaginatedS[T any](item T, count int64, hasMore bool, rowsAffected int64) PResp[T] {
	return PResp[T]{
		Status:       http.StatusOK,
		Code:         resp_const.Success,
		Message:      resp_const.Success.ToStr(),
		Body:         item,
		Count:        count,
		HasMore:      hasMore,
		RowsAffected: rowsAffected,
	}
}
func ReturnCursorPaginatedS[T any](item T, cursor CursorStruct, count int64, rowsAffected int64) PResp[T] {
	return PResp[T]{
		Status:       http.StatusOK,
		Code:         resp_const.Success,
		Message:      resp_const.Success.ToStr(),
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
		Code:   resp_const.FAIL,
		Error:  message,
		//Body:         item,
	}
}

// ====================================   Non Paginated messages =======================
// RespStatusMsg is GENERIC user supply msg and Status code
func RespStatusMsg[T any](message string, status int) *GResp[T] {
	return &GResp[T]{
		Status: status,
		Error:  message,
	}
}

// RespStatusMsg is GENERIC user supply msg and Status code
func RespStatusMsgS[T any](message string, status int) GResp[T] {
	return GResp[T]{
		Status: status,
		Error:  message,
	}
}

//====================================   Success Messages =======================

// SuccessS `PARAM: item, rowsaffected`   , mostly we dont need code for sucess `P!`
func SuccessS[T any](item T, rowsAffected int64) GResp[T] {
	return GResp[T]{
		Status:       http.StatusCreated,
		Code:         resp_const.Success,
		Message:      "success",
		Body:         item,
		RowsAffected: rowsAffected,
	}
}

// SuccessCS used for  getOne, create, update & delete success, `Param: item,code, rowsaffected`
func SuccessCS[T any](item T, code resp_const.RespCode, rowsAffected int64) GResp[T] {
	return GResp[T]{
		Status:       http.StatusOK,
		Code:         code,
		Message:      code.Msg(),
		Body:         item,
		RowsAffected: rowsAffected,
	}
}

// FetchSuccessWithCountS
func FetchSuccessWithCountS[T any](item T, message string, rowsAffected, count int64) GResp[T] {
	return GResp[T]{
		Status:       http.StatusOK,
		Code:         resp_const.Success,
		Message:      message,
		Body:         item,
		RowsAffected: rowsAffected,
		Count:        count,
	}

}

//====================================   ERROR Messages =======================

func NotFoundErrS[T any](message string) GResp[T] {
	return GResp[T]{
		Status: http.StatusNotFound,
		Code:   resp_const.FAIL,
		Error:  message,
	}
}
func NotModifiedErr[T any](message string) GResp[T] {
	return GResp[T]{
		Status: http.StatusExpectationFailed,
		Code:   resp_const.NotModified,
		Error:  message,
	}
}

// ===============  Bad Request ======

func BadReqM[T any](errorMessage string) GResp[T] {
	return GResp[T]{
		Status: http.StatusBadRequest,
		Code:   resp_const.BadRequest,
		Error:  errorMessage,
		Ok:     false,
	}
}
func BadReqC[T any](code resp_const.RespCode) GResp[T] {
	return GResp[T]{
		Status:  http.StatusBadRequest,
		Code:    code,
		Error:   code.Msg(),
		Message: code.Msg(),
	}
}

// BadReqRespMsgCode accept message & response code,
func BadReqRespMsgCode[T any](message string, code resp_const.RespCode) GResp[T] {
	return GResp[T]{
		Status: http.StatusBadRequest,
		Code:   code,
		Error:  message,
	}
}

//=================  other errors

// InternalErrMS : m shows accepts Message
func InternalErrMS[T any](message string) GResp[T] {
	return GResp[T]{
		Status: http.StatusInternalServerError,
		Code:   resp_const.FAIL,
		Error:  message,
		//Item:         item,
	}
}

func NotModifiedErrM[T any](message string) *GResp[T] {
	return &GResp[T]{
		Status: http.StatusExpectationFailed,
		Code:   resp_const.NotModified,
		Error:  message,
	}
}
