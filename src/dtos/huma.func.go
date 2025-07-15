package dtos

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	cmn "github.com/birukbelay/gocmn/src/logger"
)

//===========================   Paginated Responses  ================

// PHumaReturn returns a paginated huma response
func PHumaReturn[T any](resp PResp[T], err error) (*HumaResponse[PResp[T]], error) {
	if err != nil {
		cmn.LogTraceN("error happned", err.Error(), 3)

		return &HumaResponse[PResp[T]]{
			Status: resp.Status,
			Body:   resp,
		}, huma.NewError(resp.Status, err.Error())
	}
	return &HumaResponse[PResp[T]]{
		Status: resp.Status,
		Body:   resp,
	}, err
}

// PHumaErr returns ERROR for paginated huma resp
func PHumaErr[T any](err string, status int) (*HumaResponse[PResp[T]], error) {

	return &HumaResponse[PResp[T]]{
		Status: status,
		// Body:   *resp,
	}, huma.NewError(status, err)

}

//================  Non Paginated responses ========
// =================================================

//=============  Responses with the `Item, status & err param`

// HumaReturnTS takes `the Item, status & error`
func HumaReturnTS[T any](item T, status int, err error) (*HumaResponse[GResp[T]], error) {

	if err != nil {
		return &HumaResponse[GResp[T]]{
			Status: status,
			Body: GResp[T]{
				Status: 400,
				Error:  err.Error(),
			},
		}, huma.NewError(status, err.Error())
	}
	return &HumaResponse[GResp[T]]{
		Status: status,
		Body: GResp[T]{
			Body: item,
		},
	}, nil
}

//========================  Response With `Gresp param` ============

// HumaReturnG accepts Gresp, err
func HumaReturnG[T any](resp GResp[T], err error) (*HumaResponse[GResp[T]], error) {
	if err != nil {
		cmn.LogTraceN("error happned", err.Error(), 3)

		return &HumaResponse[GResp[T]]{
			Status: resp.Status,
			Body:   resp,
		}, huma.NewError(resp.Status, err.Error())
	}
	return &HumaResponse[GResp[T]]{
		Status: resp.Status,
		Body:   resp,
	}, err
}

func HumaReturnGWithCookie[T any](resp GResp[T], err error, cookie []http.Cookie) (*HumaResponse[GResp[T]], error) {
	if err != nil {
		cmn.LogTraceN("error happned", err.Error(), 3)

		return &HumaResponse[GResp[T]]{
			Status: resp.Status,
			Body:   resp,
		}, huma.NewError(resp.Status, err.Error())
	}
	return &HumaResponse[GResp[T]]{
		Status:    resp.Status,
		Body:      resp,
		SetCookie: cookie,
	}, err
}

// HumaReturnSE return huma error form the status & error message
func HumaReturnSE[T any](status int, errMsg string) (*HumaResponse[GResp[T]], error) {
	return &HumaResponse[GResp[T]]{
		Status: status,
		Body: GResp[T]{
			Status: status,
			Error:  errMsg,
		},
	}, huma.NewError(status, errMsg)
}

// =====================    unused   ==========================
// ============================================================

func ReturnGErrorB[T any](status int, errMsg string) *HumaResponse[GResp[T]] {
	return &HumaResponse[GResp[T]]{
		Status: status,
		Body: GResp[T]{
			Status: status,
			Error:  errMsg,
		},
	}
}

func ReturnBody[T any](item T, status int) *HumaResponse[T] {
	return &HumaResponse[T]{
		Status: status,
		Body:   item,
	}
}

//=====================  TO depricate

// GHumaReturn accepts pointer, and return huma with GResp: TO DEPRICATE
func GHumaReturn[T any](resp *GResp[T], err error) (*HumaResponse[GResp[T]], error) {
	if err != nil {
		cmn.LogTraceN("error happned", err.Error(), 3)

		return &HumaResponse[GResp[T]]{
			Status: resp.Status,
			Body:   *resp,
		}, huma.NewError(resp.Status, err.Error())
	}
	return &HumaResponse[GResp[T]]{
		Status: resp.Status,
		Body:   *resp,
	}, err
}
