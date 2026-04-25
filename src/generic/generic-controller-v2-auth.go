package generic

import (
	"context"
	"net/http"

	models "github.com/birukbelay/gocmn/src/base"
	"github.com/birukbelay/gocmn/src/consts"
	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type IController[T any, C Settable, U, F any, Q Queryable[F], S Gettable[F]] struct {
	GormConn      *gorm.DB
	AuthKey       string
	AuthValGetter string
	//Conf     *conf.EnvConfig
	//Service  IGenericGormServT[T, C, U, F]
}

func NewIController[T any, C Settable, U, F any, Q Queryable[F], S Gettable[F]](db *gorm.DB, authKey consts.AUTH_FIELD, authValGetter consts.ContextKey) *IController[T, C, U, F, Q, S] {
	return &IController[T, C, U, F, Q, S]{GormConn: db, AuthKey: authKey.S(), AuthValGetter: authValGetter.Str()}
}

func (uh *IController[T, C, U, F, Q, S]) AuthOffsetPaginated(ctx context.Context, query *Q) (*dtos.HumaResponse[dtos.PResp[[]T]], error) {
	val := *query
	filter, pagi, opt := val.GetFilter()
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	if opt == nil {
		opt = &Opt{AuthKey: &uh.AuthKey, AuthVal: &v}
	} else {
		opt.AuthKey = &uh.AuthKey
		opt.AuthVal = &v
	}
	resp, err := DbFetchManyWithOffset[T](uh.GormConn, ctx, filter, pagi, opt)
	return dtos.PHumaReturn(resp, err)
}
func (uh *IController[T, C, U, F, Q, S]) AuthCursorPaginated(ctx context.Context, query *Q) (*dtos.HumaResponse[dtos.PResp[[]T]], error) {
	val := *query
	filter, pagi, opt := val.GetFilter()
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	if opt == nil {
		opt = &Opt{AuthKey: &uh.AuthKey, AuthVal: &v}
	} else {
		opt.AuthKey = &uh.AuthKey
		opt.AuthVal = &v
	}
	resp, err := DbFetchManyWithCursor[T](uh.GormConn, ctx, filter, pagi, opt)
	return dtos.PHumaReturn(resp, err)
}

// ===================  Single Ops =====================

func (uh *IController[T, C, U, F, Q, S]) AuthCountRecords(ctx context.Context, filter *F) (*dtos.HumaResponse[dtos.GResp[int64]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbCount[T](uh.GormConn, ctx, filter, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IController[T, C, U, F, Q, S]) AuthGetOneById(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbGetOne[T](uh.GormConn, ctx, models.Base{ID: filter.ID}, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v, Debug: true})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IController[T, C, U, F, Q, S]) AuthGetOneByIdOpt(ctx context.Context, q *S) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	val := *q
	id, _, opt := val.GetFilter()
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	if opt == nil {
		opt = &Opt{AuthKey: &uh.AuthKey, AuthVal: &v}
	} else {
		opt.AuthKey = &uh.AuthKey
		opt.AuthVal = &v
	}
	resp, err := DbGetOne[T](uh.GormConn, ctx, models.Base{ID: id.Val}, opt)
	return dtos.HumaReturnG(resp, err)
}
func (uh *IController[T, C, U, F, Q, S]) AuthGetOneByFilter(ctx context.Context, filter *F) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbGetOne[T](uh.GormConn, ctx, filter, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IController[T, C, U, F, Q, S]) AuthGetOneByQuery(ctx context.Context, query *Q) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	val := *query
	filter, _, opt := val.GetFilter()
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	if opt == nil {
		opt = &Opt{AuthKey: &uh.AuthKey, AuthVal: &v}
	} else {
		opt.AuthKey = &uh.AuthKey
		opt.AuthVal = &v
	}
	resp, err := DbGetOne[T](uh.GormConn, ctx, filter, opt)
	return dtos.HumaReturnG(resp, err)
}

func (uh *IController[T, C, U, F, Q, S]) AuthGetOneByPath(ctx context.Context, query *S) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	val := *query
	sval, filter, opt := val.GetFilter()
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	whereStr := []WhereStr{{Query: sval.Key, Args: []any{sval.Val}}}
	if opt == nil {
		opt = &Opt{AuthKey: &uh.AuthKey, AuthVal: &v, WhereQuery: whereStr}
	} else {
		opt.AuthKey = &uh.AuthKey
		opt.AuthVal = &v
		opt.WhereQuery = append(opt.WhereQuery, whereStr...)
	}
	resp, err := DbGetOne[T](uh.GormConn, ctx, filter, opt)
	return dtos.HumaReturnG(resp, err)
}

// ===================  Non Fetch Ops =====================

func (uh *IController[T, C, U, F, Q, S]) AuthCreateOne(ctx context.Context, dto *dtos.HumaReqBody[C]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	body := dto.Body
	body.SetOnCreate(v)
	resp, err := DbCreateOne[T](uh.GormConn, ctx, body, &Opt{Debug: true})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IController[T, C, U, F, Q, S]) AuthUpdateOneById(ctx context.Context, dto *dtos.HumaReqBodyId[U]) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbUpdateByFilter[T](uh.GormConn, ctx, models.Base{ID: dto.ID}, dto.Body, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}

func (uh *IController[T, C, U, F, Q, S]) AuthDeleteOneByID(ctx context.Context, filter *dtos.HumaInputId) (*dtos.HumaResponse[dtos.GResp[T]], error) {
	v, ok := ctx.Value(uh.AuthValGetter).(string)
	if !ok {
		return nil, huma.NewError(http.StatusUnauthorized, "The Token is Not Correct Form")
	}
	resp, err := DbDeleteByFilter[T](uh.GormConn, ctx, models.Base{ID: filter.ID}, &Opt{AuthKey: &uh.AuthKey, AuthVal: &v})
	return dtos.HumaReturnG(resp, err)
}

//DELETE & UPDATE by filter are very risky, so not implemented
