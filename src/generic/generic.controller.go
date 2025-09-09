package generic

import (
	"context"

	"gorm.io/gorm"

	"github.com/birukbelay/gocmn/src/dtos"
)

type IGenericController[T, C, U, F any, Q Queryable[F]] struct {
	GormConn      *gorm.DB
	AuthKey       string
	AuthValGetter string
	//Conf     *conf.EnvConfig
	//Service  IGenericGormServT[T, C, U, F]
}

type Input[T any, S any] struct {
	Filter T
	Query  S
}

func NewGenericController[T, C, U, F any, Q Queryable[F]](db *gorm.DB) *IGenericController[T, C, U, F, Q] {
	return &IGenericController[T, C, U, F, Q]{GormConn: db}
}

func (uh *IGenericController[T, C, U, F, Q]) OffsetPaginated(ctx context.Context, query *Q) (*dtos.HumaResponse[dtos.PResp[[]T]], error) {
	val := *query
	filter, pagi,opt := val.GetFilter()
	// filter.PaginationInput.Select = selectedFields
	// filter.PaginationInput.SortBy = sort
	resp, err := DbFetchManyWithOffset[T](uh.GormConn, ctx, filter, pagi, opt)
	return dtos.PHumaReturn(resp, err)
}
func (uh *IGenericController[T, C, U, F, Q]) CursorPaginated(ctx context.Context, query *Q) (*dtos.HumaResponse[dtos.PResp[[]T]], error) {
	val := *query
	filter, pagi,opt := val.GetFilter()
	// filter.PaginationInput.Select = selectedFields
	// filter.PaginationInput.SortBy = sort
	resp, err := DbFetchManyWithCursor[T](uh.GormConn, ctx, filter, pagi, opt)
	return dtos.PHumaReturn(resp, err)
}

//=======================  Single Operations

func (uh *IGenericController[T, C, U, F, Q]) CreateOne(ctx context.Context, inputs *dtos.HumaReqBody[C]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbCreateOne[T](uh.GormConn, ctx, inputs.Body, &Opt{})
	return dtos.HumaReturnG(resp, err)
}
func (uh *IGenericController[T, C, U, F, Q]) GetOneByFilter(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbGetOne[T](uh.GormConn, ctx, filter, &Opt{})
	return dtos.HumaReturnG(resp, err)
}
func (uh *IGenericController[T, C, U, F, Q]) GetOneById(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbGetOneByID[T](uh.GormConn, ctx, filter.ID, &Opt{})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericController[T, C, U, F, Q]) UpdateOneById(ctx context.Context, dto *dtos.HumaReqBodyId[U]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbUpdateOneById[T](uh.GormConn, ctx, dto.ID, dto.Body, &Opt{})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericController[T, C, U, F, Q]) DeleteOneByID(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbDeleteOneById[T](uh.GormConn, ctx, filter.ID, &Opt{})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericController[T, C, U, F, Q]) CountRecords(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[int64]], error) {
	resp, err := DbCount[T](uh.GormConn, ctx, filter, &Opt{})
	return dtos.HumaReturnG(resp, err)
}

// ===================== Batch operations ===========
func (uh *IGenericController[T, C, U, F, Q]) CreateMany(ctx context.Context, inputs *dtos.HumaReqBody[[]C]) (*dtos.HumaResponse[dtos.GResp[[]T]], error) {
	resp, err := DbCreateMany[T](uh.GormConn, ctx, inputs.Body, &Opt{})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericController[T, C, U, F, Q]) DeleteMany(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[int64]], error) {
	resp, err := DbDeleteMany[T](uh.GormConn, ctx, filter, &Opt{})
	return dtos.HumaReturnG(resp, err)
}

// func (uh *IGenericController[T, C, U, F, Q]) UpdateMany(ctx context.Context, inputs *struct {
// 	Q    F
// 	Body U
// }) (*dtos.HumaResponse[dtos.GResp[int64]], error) {
// 	resp, err := DbUpdateMany[T](uh.GormConn, ctx, inputs.Q, inputs.Body, &Opt{})
// 	return dtos.HumaReturnG(resp, err)
// }
