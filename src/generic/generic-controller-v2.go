package generic

import (
	"context"

	"github.com/birukbelay/gocmn/src/dtos"
)

func (uh *IController[T, C, U, F, Q, S]) UnAuthOffsetPaginated(ctx context.Context, query *Q) (*dtos.HumaResponse[dtos.PResp[[]T]], error) {
	val := *query
	filter, pagi, opt := val.GetFilter()

	resp, err := DbFetchManyWithOffset[T](uh.GormConn, ctx, filter, pagi, opt)
	return dtos.PHumaReturn(resp, err)
}
func (uh *IController[T, C, U, F, Q, S]) UnAuthCursorPaginated(ctx context.Context, query *Q) (*dtos.HumaResponse[dtos.PResp[[]T]], error) {
	val := *query
	filter, pagi, opt := val.GetFilter()

	resp, err := DbFetchManyWithCursor[T](uh.GormConn, ctx, filter, pagi, opt)
	return dtos.PHumaReturn(resp, err)
}

//=======================  Single Operations

//========== Fetch ops

func (uh *IController[T, C, U, F, Q, S]) UnAuthGetOneByFilter(ctx context.Context, filter *F) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbGetOne[T](uh.GormConn, ctx, filter, nil)
	return dtos.HumaReturnG(resp, err)
}
func (uh *IController[T, C, U, F, Q, S]) UnAuthGetOneByQuery(ctx context.Context, query *Q) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	val := *query
	filter, _, opt := val.GetFilter()
	resp, err := DbGetOne[T](uh.GormConn, ctx, filter, opt)
	return dtos.HumaReturnG(resp, err)
}
func (uh *IController[T, C, U, F, Q, S]) UnAuthGetOneById(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbGetOneByID[T](uh.GormConn, ctx, filter.ID, nil)
	return dtos.HumaReturnG(resp, err)
}
func (uh *IController[T, C, U, F, Q, S]) UnAuthGetOneByIdOpt(ctx context.Context, q *S) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	val := *q

	keyVal, _, opt := val.GetFilter()

	resp, err := DbGetOneByID[T](uh.GormConn, ctx, keyVal.Val, opt)
	return dtos.HumaReturnG(resp, err)
}
func (uh *IController[T, C, U, F, Q, S]) UnAuthGetOneByPath(ctx context.Context, q *S) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	val := *q

	keyVal, filter, opt := val.GetFilter()
	whereStr := []WhereStr{{Query: keyVal.Key, Args: []any{keyVal.Val}}}
	if opt == nil {
		opt = &Opt{WhereQuery: whereStr}
	} else {
		opt.WhereQuery = append(opt.WhereQuery, whereStr...)
	}
	resp, err := DbGetOne[T](uh.GormConn, ctx, filter, opt)
	return dtos.HumaReturnG(resp, err)
}

func (uh *IController[T, C, U, F, Q, S]) UnAuthCountRecords(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[int64]], error) {
	resp, err := DbCount[T](uh.GormConn, ctx, filter, nil)
	return dtos.HumaReturnG(resp, err)
}

//========== Non ops

func (uh *IController[T, C, U, F, Q, S]) UnAuthCreateOne(ctx context.Context, inputs *dtos.HumaReqBody[C]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbCreateOne[T](uh.GormConn, ctx, inputs.Body, nil)
	return dtos.HumaReturnG(resp, err)
}

func (uh *IController[T, C, U, F, Q, S]) UnAuthUpdateOneById(ctx context.Context, dto *dtos.HumaReqBodyId[U]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbUpdateOneById[T](uh.GormConn, ctx, dto.ID, dto.Body, nil)
	return dtos.HumaReturnG(resp, err)
}

func (uh *IController[T, C, U, F, Q, S]) UnAuthDeleteOneByID(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbDeleteOneById[T](uh.GormConn, ctx, filter.ID, nil)
	return dtos.HumaReturnG(resp, err)
}

// ===================== Batch operations ===========
func (uh *IController[T, C, U, F, Q, S]) UnAuthCreateMany(ctx context.Context, inputs *dtos.HumaReqBody[[]C]) (*dtos.HumaResponse[dtos.GResp[[]T]], error) {
	resp, err := DbCreateMany[T](uh.GormConn, ctx, inputs.Body, nil)
	return dtos.HumaReturnG(resp, err)
}
