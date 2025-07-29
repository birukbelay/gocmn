package generic

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"

	models "github.com/birukbelay/gocmn/src/base"
	"github.com/birukbelay/gocmn/src/dtos"
)

func NewGenericAuthController[T, C, U, F any, Q Queryable](db *gorm.DB, authKey, authValGetter string) *IGenericController[T, C, U, F, Q] {
	return &IGenericController[T, C, U, F, Q]{GormConn: db, AuthKey: authKey, AuthValGetter: authValGetter}
}

//	func (uh *IGenericController[T, C, U, F, Q]) AuthCreateOne(ctx context.Context, inputs *dtos.HumaReqBody[C]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
//		k, ok := ctx.Value("validationKey").(string)
//		if !ok {
//			return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
//		}
//		v, ok := ctx.Value("validationValue").(string)
//		if !ok {
//			return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
//		}
//		resp, err := DbCreateOne[T](uh.GormConn, ctx, inputs.Body, &Opt{AuthKey: &k, AuthVal: &v})
//		return dtos.HumaReturnG(resp, err)
//	}
func (uh *IGenericController[T, C, U, F, Q]) AuthGetOneByFilter(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbGetOne[T](uh.GormConn, ctx, filter, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}
func (uh *IGenericController[T, C, U, F, Q]) AuthGetOneById(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbGetOne[T](uh.GormConn, ctx, models.Base{ID: filter.ID}, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v, Debug: true})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericController[T, C, U, F, Q]) AuthUpdateOneById(ctx context.Context, dto *dtos.HumaReqBodyId[U]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbUpdateByFilter[T](uh.GormConn, ctx, models.Base{ID: dto.ID}, dto.Body, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericController[T, C, U, F, Q]) AuthDeleteOneByID(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbDeleteByFilter[T](uh.GormConn, ctx, models.Base{ID: filter.ID}, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericController[T, C, U, F, Q]) AuthCountRecords(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[int64]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbCount[T](uh.GormConn, ctx, filter, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}
