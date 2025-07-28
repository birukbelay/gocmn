package generic

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/birukbelay/gocmn/src/dtos"
)

func (uh *IGenericController[T, C, U, F, Q]) AuthCreateOne(ctx context.Context, inputs *dtos.HumaReqBody[C]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	k, ok := ctx.Value("validationKey").(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	v, ok := ctx.Value("validationValue").(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbCreateOne[T](uh.GormConn, ctx, inputs.Body, &Opt{AuthKey: &k, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}
func (uh *IGenericController[T, C, U, F, Q]) AuthGetOneByFilter(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbGetOne[T](uh.GormConn, ctx, filter, &Opt{})
	return dtos.HumaReturnG(resp, err)
}
func (uh *IGenericController[T, C, U, F, Q]) AuthGetOneById(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbGetOneByID[T](uh.GormConn, ctx, filter.ID, &Opt{})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericController[T, C, U, F, Q]) AuthUpdateOneById(ctx context.Context, dto *dtos.HumaReqBodyId[U]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbUpdateOneById[T](uh.GormConn, ctx, dto.ID, dto.Body, &Opt{})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericController[T, C, U, F, Q]) AuthDeleteOneByID(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	resp, err := DbDeleteOneById[T](uh.GormConn, ctx, filter.ID, &Opt{})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericController[T, C, U, F, Q]) AuthCountRecords(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[int64]], error) {
	resp, err := DbCount[T](uh.GormConn, ctx, filter, &Opt{})
	return dtos.HumaReturnG(resp, err)
}
