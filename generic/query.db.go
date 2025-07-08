package generic

import (
	"context"
	"errors"
	models "github.com/birukbelay/gocmn/base"
	"github.com/birukbelay/gocmn/dtos"
	cmn "github.com/birukbelay/gocmn/logger"
	respC "github.com/birukbelay/gocmn/resp_codes"
	"github.com/birukbelay/gocmn/util"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ==============================================   Fetch Many Operations ================
// ========================================================================================

func DbFetchManyWithOffset[T any](u *gorm.DB, ctx context.Context, filter any, pagi dtos.PaginationInput, options *Opt) (dtos.PResp[[]T], error) {
	orderedBy, sortDir, limit, page := util.GetPagiValues(pagi)
	orderString := orderedBy + " " + sortDir
	var lst []T
	query := u.WithContext(ctx)
	//select, preload and debug
	query = DebugPreloadSelect(query, options, pagi.Select)
	//tags
	_ = searchTags(query, pagi.Tags)
	//start & end date
	_, err := addStartEndDate(query, pagi.StartDate, pagi.EndDate)
	if err != nil {
		return dtos.PInternalErrS[[]T](err.Error()), err
	}
	//text search & like queries
	query = addSearchFilters(query, pagi)
	query = AddInQueries(query, options)

	//use a differnet session after this
	query = query.Session(&gorm.Session{})
	resp := query.Model(lst).Where(filter).Limit(limit + 1).Order(orderString).Offset(page * limit).Find(&lst)
	if resp.Error != nil {
		log.Err(resp.Error).Msg("fetch Error")
		return dtos.PInternalErrS[[]T](resp.Error.Error()), resp.Error
	}
	hasMore := false
	if len(lst) > limit {
		lst = lst[:len(lst)-1]
		hasMore = true
	}

	var count int64
	cnt := query.Model(lst).Where(filter).Count(&count)
	if cnt.Error != nil {
		return dtos.PInternalErrS[[]T](cnt.Error.Error()), cnt.Error
	}

	return dtos.ReturnOffsetPaginatedS[[]T](lst, count, hasMore, resp.RowsAffected), nil
}

func DbFetchManyWithCursor[T any](u *gorm.DB, ctx context.Context, filter any, pagi dtos.PaginationInput, options *Opt) (dtos.PResp[[]T], error) {

	var dataList []T
	query := u.WithContext(ctx).Where(filter)

	//select, preload and debug
	query = DebugPreloadSelect(query, options, pagi.Select)
	//tags
	_ = searchTags(query, pagi.Tags)

	// ----! start cursor based pagination
	orderedBy, sortDir, limit, _ := util.GetPagiValues(pagi)
	orderString := orderedBy + " " + sortDir

	query, orderString, forward, err := CreateCursorQuery(query, pagi.Cursor, orderString, sortDir)
	if err != nil {
		return dtos.PInternalErrS[[]T](err.Error()), err
	}
	//use a differnet session after this
	query = query.Session(&gorm.Session{})
	resp := query.Model(dataList).Limit(limit + 1).Order(orderString).Find(&dataList)
	if resp.Error != nil {
		return dtos.PInternalErrS[[]T](resp.Error.Error()), resp.Error
	}

	var count int64
	cnt := query.Model(dataList).Count(&count)
	if cnt.Error != nil {
		return dtos.PInternalErrS[[]T](cnt.Error.Error()), cnt.Error
	}
	//create Cursors For Pagination
	cursor, newDataLst, err := CreateCursor(limit, forward, dataList, orderedBy, pagi.Cursor)
	if err != nil {
		return dtos.PInternalErrS[[]T](err.Error()), err
	}

	return dtos.ReturnCursorPaginatedS(newDataLst, cursor, count, resp.RowsAffected), nil
}

// ==============================================   Single Operations ================
//========================================================================================

func DbGetOne[T any](u *gorm.DB, ctx context.Context, filter any, options *Opt) (dtos.GResp[T], error) {
	if isEmptyStruct(filter) {
		return dtos.BadRequestMsgS[T](respC.EmptyFilter.Msg()), errors.New(respC.EmptyFilter.Msg())
	}
	result := new(T)
	query := u.WithContext(ctx)
	query = DebugPreloadSelect(query, options, nil)

	query = query.Model(result).Where(filter)
	resp := query.First(result)
	if resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			return dtos.NotFoundErrS[T](resp.Error.Error()), resp.Error
		}
		//cmn.LogTrace("the error is", resp)
		return dtos.InternalErrS[T](resp.Error.Error()), resp.Error
	}
	return dtos.SuccessS[T](*result, resp.RowsAffected), nil
}

func DbGetOneByID[T any](u *gorm.DB, ctx context.Context, id string, options *Opt) (dtos.GResp[T], error) {
	return DbGetOne[T](u, ctx, models.Base{ID: id}, options)
}

func DbCreateOne[T any](u *gorm.DB, ctx context.Context, value any, options *Opt) (dtos.GResp[T], error) {
	result := new(T)
	if err := mapstructure.Decode(value, &result); err != nil {
		return dtos.BadRequestMsgS[T](err.Error()), err
	}
	query := u.WithContext(ctx)
	if options != nil {
		if options.Debug {
			query = query.Debug()
		}
		for _, pre := range options.Preloads {
			query = query.Preload(pre)
		}
	}
	resp := query.Model(&result).Create(&result)
	if resp.Error != nil {
		return dtos.InternalErrS[T](resp.Error.Error()), resp.Error
	}
	return dtos.SuccessS[T](*result, resp.RowsAffected), nil
}

func DbCount[T any](u *gorm.DB, ctx context.Context, filter any, options *Opt) (int64, error) {
	var mdl T
	var count int64
	cntQuery := u.WithContext(ctx)
	if options != nil {
		if options.Debug {
			cntQuery = cntQuery.Debug()
		}
	}
	cnt := cntQuery.Model(mdl).Where(filter).Count(&count)
	if cnt.Error != nil {
		return 0, cnt.Error
	}
	return count, nil
}

func DbUpdateByFilter[T any](u *gorm.DB, ctx context.Context, filter, updateDto any, options *Opt) (dtos.GResp[T], error) {
	if isEmptyStruct(filter) {
		return dtos.BadRequestMsgS[T](respC.EmptyFilter.Msg()), errors.New(respC.EmptyFilter.Msg())
	}
	result := new(T)
	query := u.WithContext(ctx)
	if options != nil {
		if options.Debug {
			query = query.Debug()
		}
		for _, pre := range options.Preloads {
			query = query.Preload(pre)
		}
	}
	resp := query.Clauses(clause.Returning{}).Model(&result).Where(filter).Updates(updateDto)
	if resp.Error != nil {
		return dtos.InternalErrS[T](resp.Error.Error()), resp.Error
	}
	if resp.RowsAffected == 0 {
		return dtos.NotFoundErrS[T](respC.NoRecordsUpdated.Msg()), errors.New(respC.NoRecordsUpdated.Msg())
	}
	return dtos.SuccessCS[T](*result, respC.UpdateSuccess, resp.RowsAffected), nil
}

func DbUpdateOneById[T any](u *gorm.DB, ctx context.Context, id string, updateDto any, options *Opt) (dtos.GResp[T], error) {
	//usr:=users
	return DbUpdateByFilter[T](u, ctx, models.Base{ID: id}, updateDto, options)
}

func DbDeleteByFilter[T any](u *gorm.DB, ctx context.Context, filter any, options *Opt) (dtos.GResp[T], error) {
	if isEmptyStruct(filter) {
		return dtos.BadRequestMsgS[T](respC.EmptyFilter.Msg()), errors.New(respC.EmptyFilter.Msg())
	}
	result := new(T)
	query := u.WithContext(ctx)
	if options != nil {
		if options.Debug {
			query = query.Debug()
		}
		for _, pre := range options.Preloads {
			query = query.Preload(pre)
		}
	}

	resp := query.Clauses(clause.Returning{}).Model(result).Where(util.UnwrapAny(filter)).Delete(&result)

	if resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			return dtos.NotFoundErrS[T](resp.Error.Error()), resp.Error
		}
		cmn.LogTrace("the error is", resp)
		return dtos.NotFoundErrS[T](resp.Error.Error()), resp.Error
	}
	if resp.RowsAffected == 0 {
		return dtos.NotFoundErrS[T]("NO records were deleted"), errors.New("NO records were deleted")
	}
	return dtos.SuccessCS[T](*result, respC.DeleteSuccess, resp.RowsAffected), nil

}
func DbDeleteOneById[T any](u *gorm.DB, ctx context.Context, id string, options *Opt) (dtos.GResp[T], error) {
	return DbDeleteByFilter[T](u, ctx, models.Base{ID: id}, options)
}

//-===============================  Batch Operations ===================================
//======================================================================================

func DbCreateMany[T any](u *gorm.DB, ctx context.Context, value any, options *Opt) (dtos.GResp[[]T], error) {
	var result []T
	if err := mapstructure.Decode(value, &result); err != nil {
		return dtos.InternalErrS[[]T](err.Error()), err
	}
	query := u.WithContext(ctx)
	if options != nil {
		if options.Debug {
			query = query.Debug()
		}
		for _, pre := range options.Preloads {
			query = query.Preload(pre)
		}
	}
	resp := query.Clauses(clause.OnConflict{DoNothing: true}, clause.Returning{}).Model(&result).Create(&result)
	if resp.Error != nil {
		return dtos.InternalErrS[[]T](resp.Error.Error()), resp.Error
	}
	return dtos.SuccessS[[]T](result, resp.RowsAffected), nil
}

func DbUpdateMany[T any](u *gorm.DB, ctx context.Context, filter, updateDto any, options *Opt) (dtos.GResp[int64], error) {
	if isEmptyStruct(filter) {
		return dtos.BadRequestMsgS[int64](respC.EmptyFilter.Msg()), errors.New(respC.EmptyFilter.Msg())
	}
	result := new(T)
	query := u.WithContext(ctx)
	if options != nil {
		if options.Debug {
			query = query.Debug()
		}
		for _, pre := range options.Preloads {
			query = query.Preload(pre)
		}
	}
	resp := query.Clauses(clause.Returning{}).Model(&result).Where(filter).Updates(updateDto)
	if resp.Error != nil {
		return dtos.InternalErrS[int64](resp.Error.Error()), resp.Error
	}
	if resp.RowsAffected == 0 {
		return dtos.NotFoundErrS[int64](respC.NoRecordsUpdated.Msg()), errors.New(respC.NoRecordsUpdated.Msg())
	}
	return dtos.SuccessCS[int64](resp.RowsAffected, respC.UpdateSuccess, resp.RowsAffected), nil
}

func DbDeleteMany[T any](u *gorm.DB, ctx context.Context, filter any, options *Opt) (dtos.GResp[int64], error) {
	if isEmptyStruct(filter) {
		return dtos.BadRequestMsgS[int64](respC.EmptyFilter.Msg()), errors.New(respC.EmptyFilter.Msg())
	}
	result := new(T)
	query := u.WithContext(ctx)
	if options != nil {
		if options.Debug {
			query = query.Debug()
		}
		for _, pre := range options.Preloads {
			query = query.Preload(pre)
		}
	}
	resp := query.Model(result).Where(util.UnwrapAny(filter)).Delete(&result)

	if resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			return dtos.NotFoundErrS[int64](resp.Error.Error()), resp.Error
		}
		//cmn.LogTrace("the error is", resp)
		return dtos.InternalErrS[int64](resp.Error.Error()), resp.Error
	}
	if resp.RowsAffected == 0 {
		return dtos.NotFoundErrS[int64]("NO records were deleted"), errors.New("NO records were deleted")
	}
	return dtos.SuccessCS[int64](resp.RowsAffected, respC.DeleteSuccess, resp.RowsAffected), nil

}

//--------- New Methods

func DbUpsertManyWithValues[T any](u *gorm.DB, ctx context.Context, value any, options Opt) (dtos.GResp[[]T], error) {
	var result []T
	if err := mapstructure.Decode(value, &result); err != nil {
		return dtos.InternalErrS[[]T](err.Error()), err
	}

	query := u.WithContext(ctx)
	if options.Debug {
		query = query.Debug()
	}
	for _, pre := range options.Preloads {
		query = query.Preload(pre)
	}
	resp := query.Clauses(clause.OnConflict{
		Columns:   options.Columns,
		DoUpdates: clause.AssignmentColumns(options.UpdateColumns),
	}, clause.Returning{}).Model(&result).Create(&result)
	if resp.Error != nil {
		return dtos.InternalErrS[[]T](resp.Error.Error()), resp.Error
	}
	return dtos.SuccessS[[]T](result, resp.RowsAffected), nil
}
