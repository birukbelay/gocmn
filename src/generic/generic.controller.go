package generic

import (
	"context"

	"gorm.io/gorm"

	"github.com/birukbelay/gocmn/src/dtos"
)

type IGenericController[T any, F any, D any, Q Queryable] struct {
	GormConn *gorm.DB
	//Conf     *conf.EnvConfig
	//Service  IGenericGormServT[T, F, D]
}

func (uh *IGenericController[T, F, D, Q]) OffsetPaginated(ctx context.Context, filter *struct {
	Filter F
	Query  Q
	dtos.PaginationInput
}) (*dtos.HumaResponse[dtos.PResp[[]T]], error) {
	sort, selectedFields := filter.Query.GetQueries()
	filter.PaginationInput.Select = selectedFields
	filter.PaginationInput.SortBy = sort
	resp, err := DbFetchManyWithOffset[T](uh.GormConn, ctx, filter.Filter, filter.PaginationInput, &Opt{})
	return dtos.PHumaReturn(resp, err)
}
func (uh *IGenericController[T, F, D, Q]) CursorPaginated(ctx context.Context, filter *struct {
	Filter F
	Query  Q
	dtos.PaginationInput
}) (*dtos.HumaResponse[dtos.PResp[[]T]], error) {
	sort, selectedFields := filter.Query.GetQueries()
	filter.PaginationInput.Select = selectedFields
	filter.PaginationInput.SortBy = sort
	resp, err := DbFetchManyWithCursor[T](uh.GormConn, ctx, filter.Filter, filter.PaginationInput, &Opt{})
	return dtos.PHumaReturn(resp, err)
}

//=======================  Single Operations

func (uh *IGenericController[T, F, D, Q]) CreateOne(ctx context.Context, inputs *dtos.HumaReqBody[D]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbCreateOne[T](uh.GormConn, ctx, inputs.Body, &Opt{})
	return dtos.HumaReturnG(resp, err)
}
func (uh *IGenericController[T, F, D, Q]) GetOneDomainByFilter(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbGetOne[T](uh.GormConn, ctx, filter, &Opt{})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericController[T, F, D, Q]) UpdateOneById(ctx context.Context, filter *dtos.HumaReqBodyId[D]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbUpdateOneById[T](uh.GormConn, ctx, filter.ID, filter.Body, &Opt{})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericController[T, F, D, Q]) DeleteOneByID(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbDeleteOneById[T](uh.GormConn, ctx, filter.ID, &Opt{})
	return dtos.HumaReturnG[T](resp, err)
}

//===================== Batch operations ===========
