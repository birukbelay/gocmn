package generic

import (
	"context"

	"github.com/birukbelay/gocmn/src/dtos"
)

func (uh *IGenericAuthController[T, C, U, F, Q]) UnAuthOffsetPaginated(ctx context.Context, query *Q) (*dtos.HumaResponse[dtos.PResp[[]T]], error) {
	val := *query
	filter, pagi, opt := val.GetFilter()

	resp, err := DbFetchManyWithOffset[T](uh.GormConn, ctx, filter, pagi, opt)
	return dtos.PHumaReturn(resp, err)
}
func (uh *IGenericAuthController[T, C, U, F, Q]) UnAuthCursorPaginated(ctx context.Context, query *Q) (*dtos.HumaResponse[dtos.PResp[[]T]], error) {
	val := *query
	filter, pagi, opt := val.GetFilter()

	resp, err := DbFetchManyWithCursor[T](uh.GormConn, ctx, filter, pagi, opt)
	return dtos.PHumaReturn(resp, err)
}

//=======================  Single Operations

func (uh *IGenericAuthController[T, C, U, F, Q]) UnAuthCreateOne(ctx context.Context, inputs *dtos.HumaReqBody[C]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbCreateOne[T](uh.GormConn, ctx, inputs.Body, nil)
	return dtos.HumaReturnG(resp, err)
}
func (uh *IGenericAuthController[T, C, U, F, Q]) UnAuthGetOneByFilter(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbGetOne[T](uh.GormConn, ctx, filter, nil)
	return dtos.HumaReturnG(resp, err)
}
func (uh *IGenericAuthController[T, C, U, F, Q]) UnAuthGetOneById(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbGetOneByID[T](uh.GormConn, ctx, filter.ID, nil)
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericAuthController[T, C, U, F, Q]) UnAuthUpdateOneById(ctx context.Context, dto *dtos.HumaReqBodyId[U]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbUpdateOneById[T](uh.GormConn, ctx, dto.ID, dto.Body, nil)
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericAuthController[T, C, U, F, Q]) UnAuthDeleteOneByID(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbDeleteOneById[T](uh.GormConn, ctx, filter.ID, nil)
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericAuthController[T, C, U, F, Q]) UnAuthCountRecords(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[int64]], error) {
	resp, err := DbCount[T](uh.GormConn, ctx, filter, nil)
	return dtos.HumaReturnG(resp, err)
}

// ===================== Batch operations ===========
func (uh *IGenericAuthController[T, C, U, F, Q]) UnAuthCreateMany(ctx context.Context, inputs *dtos.HumaReqBody[[]C]) (*dtos.HumaResponse[dtos.GResp[[]T]], error) {
	resp, err := DbCreateMany[T](uh.GormConn, ctx, inputs.Body, nil)
	return dtos.HumaReturnG(resp, err)
}

// func (uh *IGenericAuthController[T, C, U, F, Q]) DeleteMany(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[int64]], error) {
// 	resp, err := DbDeleteMany[T](uh.GormConn, ctx, filter, &Opt{})
// 	return dtos.HumaReturnG(resp, err)
// }
