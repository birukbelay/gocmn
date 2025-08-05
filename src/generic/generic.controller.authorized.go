package generic

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"

	models "github.com/birukbelay/gocmn/src/base"
	"github.com/birukbelay/gocmn/src/consts"
	"github.com/birukbelay/gocmn/src/dtos"
)

type IGenericAuthController[T any, C Settable, U, F any, Q Queryable] struct {
	GormConn      *gorm.DB
	AuthKey       string
	AuthValGetter string
	//Conf     *conf.EnvConfig
	//Service  IGenericGormServT[T, C, U, F]
}

func NewGenericAuthController[T any, C Settable, U, F any, Q Queryable](db *gorm.DB, authKey consts.AUTH_FIELD, authValGetter consts.ContextKey) *IGenericAuthController[T, C, U, F, Q] {
	return &IGenericAuthController[T, C, U, F, Q]{GormConn: db, AuthKey: authKey.S(), AuthValGetter: authValGetter.Str()}
}

func (uh *IGenericAuthController[T, C, U, F, Q]) AuthCreateOne(ctx context.Context, dto *dtos.HumaReqBody[C]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	dto.Body.SetOnCreate(v)
	resp, err := DbCreateOne[T](uh.GormConn, ctx, dto.Body, &Opt{Debug: true})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericAuthController[T, C, U, F, Q]) AuthGetOneByFilter(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbGetOne[T](uh.GormConn, ctx, filter, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}
func (uh *IGenericAuthController[T, C, U, F, Q]) AuthGetOneById(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbGetOne[T](uh.GormConn, ctx, models.Base{ID: filter.ID}, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v, Debug: true})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericAuthController[T, C, U, F, Q]) AuthUpdateOneById(ctx context.Context, dto *dtos.HumaReqBodyId[U]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbUpdateByFilter[T](uh.GormConn, ctx, models.Base{ID: dto.ID}, dto.Body, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericAuthController[T, C, U, F, Q]) AuthDeleteOneByID(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbDeleteByFilter[T](uh.GormConn, ctx, models.Base{ID: filter.ID}, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IGenericAuthController[T, C, U, F, Q]) AuthCountRecords(ctx context.Context, filter F) (*dtos.HumaResponse[dtos.GResp[int64]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbCount[T](uh.GormConn, ctx, filter, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}
