package generic

import (
	"context"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	models "github.com/birukbelay/gocmn/src/base"
	"github.com/birukbelay/gocmn/src/dtos"
	cmn "github.com/birukbelay/gocmn/src/logger"
	respC "github.com/birukbelay/gocmn/src/resp_const"
	"github.com/birukbelay/gocmn/src/util"
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
	query = searchTags(query, pagi.Tags)
	query = addAllOfTags(query, pagi.AllTags)
	query = addAnyOfTags(query, pagi.AnyTags)
	//start & end date
	query, err := addStartEndDate(query, pagi.StartDate, pagi.EndDate)
	if err != nil {
		return dtos.PInternalErrS[[]T](err.Error()), err
	}

	//================.  Authorization Queries. ==========
	if options != nil {
		if options.AuthKey != nil && options.AuthVal != nil {
			queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
			query = query.Where(queryStr, options.AuthVal)
		}
		if len(options.WhereQuery) > 0 {
			for _, where := range options.WhereQuery {
				query = query.Where(where.Query, where.Args...)
			}
		}
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
	//================.  Authorization Queries. ==========
	if options != nil {
		if options.AuthKey != nil && options.AuthVal != nil {
			queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
			query = query.Where(queryStr, options.AuthVal)
		}
		if len(options.WhereQuery) > 0 {
			for _, where := range options.WhereQuery {
				query = query.Where(where.Query, where.Args...)
			}
		}
	}

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
		return dtos.BadReqM[T](respC.EmptyFilter.Msg()), errors.New(respC.EmptyFilter.Msg())
	}
	result := new(T)
	query := u.WithContext(ctx)
	query = DebugPreloadSelect(query, options, nil)
	if options != nil {
		if options.AuthKey != nil && options.AuthVal != nil {
			queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
			query = query.Where(queryStr, options.AuthVal)
		}
		if len(options.WhereQuery) > 0 {
			for _, where := range options.WhereQuery {
				query = query.Where(where.Query, where.Args...)
			}
		}
	}

	query = query.Model(result).Where(filter)
	resp := query.First(result)
	if resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			return dtos.NotFoundErrS[T](resp.Error.Error()), resp.Error
		}
		//cmn.LogTrace("the error is", resp)
		return dtos.InternalErrMS[T](resp.Error.Error()), resp.Error
	}
	return dtos.SuccessS[T](*result, resp.RowsAffected), nil
}

func DbGetOneByID[T any](u *gorm.DB, ctx context.Context, id string, options *Opt) (dtos.GResp[T], error) {
	return DbGetOne[T](u, ctx, models.Base{ID: id}, options)
}

func DbCreateOne[T any](u *gorm.DB, ctx context.Context, value any, options *Opt) (dtos.GResp[T], error) {
	result := new(T)
	if err := mapstructure.Decode(value, &result); err != nil {
		return dtos.BadReqM[T](err.Error()), err
	}
	query := u.WithContext(ctx).Clauses(clause.Returning{})
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
		return dtos.InternalErrMS[T](resp.Error.Error()), resp.Error
	}
	return dtos.SuccessS[T](*result, resp.RowsAffected), nil
}

// DbUpsertOneAllFields Update all, on conflict of these columns
func DbUpsertOneAllFields[T any](u *gorm.DB, ctx context.Context, value any, cols []clause.Column, options *Opt) (dtos.GResp[T], error) {
	result := new(T)
	if err := mapstructure.Decode(value, &result); err != nil {
		return dtos.BadReqM[T](err.Error()), err
	}
	query := u.WithContext(ctx)
	if options != nil {
		//================.  Authorization Queries. ==========
		if options.AuthKey != nil && options.AuthVal != nil {
			queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
			query = query.Where(queryStr, options.AuthVal)
		}
		if options.Debug {
			query = query.Debug()
		}
		for _, pre := range options.Preloads {
			query = query.Preload(pre)
		}
	}
	resp := query.Clauses(clause.OnConflict{
		Columns:   cols,
		UpdateAll: true,
	}, clause.Returning{}).Model(&result).Create(&result)
	if resp.Error != nil {
		return dtos.InternalErrMS[T](resp.Error.Error()), resp.Error
	}
	return dtos.SuccessS[T](*result, resp.RowsAffected), nil
}

// DbUpsertOneListedFields Update listed fields(updateCols), on conflict of these columns
func DbUpsertOneListedFields[T any](u *gorm.DB, ctx context.Context, value any, conflictingCols []clause.Column, updateCols []string, options *Opt) (dtos.GResp[T], error) {
	result := new(T)
	if err := mapstructure.Decode(value, &result); err != nil {
		return dtos.BadReqM[T](err.Error()), err
	}
	query := u.WithContext(ctx)
	if options != nil {
		if options.Debug {
			query = query.Debug()
		}
		for _, pre := range options.Preloads {
			query = query.Preload(pre)
		}
		//================.  Authorization Queries. ==========
		if options.AuthKey != nil && options.AuthVal != nil {
			queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
			query = query.Where(queryStr, options.AuthVal)
		}

	}
	resp := query.Clauses(clause.OnConflict{
		Columns:   conflictingCols,
		DoUpdates: clause.AssignmentColumns(updateCols),
	}).Model(&result).Create(&result)
	if resp.Error != nil {
		return dtos.InternalErrMS[T](resp.Error.Error()), resp.Error
	}
	return dtos.SuccessS[T](*result, resp.RowsAffected), nil
}

func DbCount[T any](u *gorm.DB, ctx context.Context, filter any, options *Opt) (dtos.GResp[int64], error) {
	var mdl T
	var count int64
	cntQuery := u.WithContext(ctx)
	if options != nil {
		if options.Debug {
			cntQuery = cntQuery.Debug()
		}
		//================.  Authorization Queries. ==========
		if options.AuthKey != nil && options.AuthVal != nil {
			queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
			cntQuery = cntQuery.Where(queryStr, options.AuthVal)
		}

	}
	cnt := cntQuery.Model(mdl).Where(filter).Count(&count)
	if cnt.Error != nil {
		return dtos.InternalErrMS[int64](cnt.Error.Error()), cnt.Error
	}
	return dtos.SuccessS(count, cnt.RowsAffected), nil
}

func DbFetchWihtIn[T any](u *gorm.DB, ctx context.Context, inKey string, names []string, options *Opt) (dtos.PResp[[]T], error) {
	var dataList []T
	query := u.WithContext(ctx)
	query = DebugPreloadSelect(query, options, nil)
	query = query.Where(fmt.Sprintf("%s IN (?)", inKey), names)
	//================.  Authorization Queries. ==========
	if options != nil {
		if options.AuthKey != nil && options.AuthVal != nil {
			queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
			query = query.Where(queryStr, options.AuthVal)
		}
	}
	//use a differnet session after this
	query = query.Session(&gorm.Session{})
	resp := query.Model(dataList).Find(&dataList)
	if resp.Error != nil {
		return dtos.PInternalErrS[[]T](resp.Error.Error()), resp.Error
	}

	var count int64
	cnt := query.Model(dataList).Count(&count)
	if cnt.Error != nil {
		return dtos.PInternalErrS[[]T](cnt.Error.Error()), cnt.Error
	}

	return dtos.ReturnOffsetPaginatedS[[]T](dataList, count, true, resp.RowsAffected), nil
}

func DbUpdateByFilter[T any](u *gorm.DB, ctx context.Context, filter, updateDto any, options *Opt) (dtos.GResp[T], error) {
	if isEmptyStruct(filter) {
		return dtos.BadReqM[T](respC.EmptyFilter.Msg()), errors.New(respC.EmptyFilter.Msg())
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
		//================.  Authorization Queries. ==========
		if options.AuthKey != nil && options.AuthVal != nil {
			queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
			query = query.Where(queryStr, options.AuthVal)
		}

	}
	resp := query.Clauses(clause.Returning{}).Model(&result).Where(filter).Updates(updateDto)
	if resp.Error != nil {
		return dtos.InternalErrMS[T](resp.Error.Error()), resp.Error
	}
	if resp.RowsAffected == 0 {
		return dtos.NotModifiedErr[T](respC.NoRecordsUpdated.Msg()), errors.New(respC.NoRecordsUpdated.Msg())
	}
	return dtos.SuccessCS[T](*result, respC.UpdateSuccess, resp.RowsAffected), nil
}

func DbUpdateOneById[T any](u *gorm.DB, ctx context.Context, id string, updateDto any, options *Opt) (dtos.GResp[T], error) {
	//usr:=users
	return DbUpdateByFilter[T](u, ctx, models.Base{ID: id}, updateDto, options)
}

func DbDeleteByFilter[T any](u *gorm.DB, ctx context.Context, filter any, options *Opt) (dtos.GResp[T], error) {
	if isEmptyStruct(filter) {
		return dtos.BadReqM[T](respC.EmptyFilter.Msg()), errors.New(respC.EmptyFilter.Msg())
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
		//================.  Authorization Queries. ==========
		if options.AuthKey != nil && options.AuthVal != nil {
			queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
			query = query.Where(queryStr, options.AuthVal)
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
		return dtos.NotModifiedErr[T]("NO records were deleted"), errors.New("NO records were deleted")
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
		return dtos.InternalErrMS[[]T](err.Error()), err
	}
	query := u.WithContext(ctx)
	if options != nil {
		if options.Debug {
			query = query.Debug()
		}
		for _, pre := range options.Preloads {
			query = query.Preload(pre)
		}
		//================.  Authorization Queries. ==========
		if options.AuthKey != nil && options.AuthVal != nil {
			queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
			query = query.Where(queryStr, options.AuthVal)
		}

	}
	resp := query.Clauses(clause.OnConflict{DoNothing: true}, clause.Returning{}).Model(&result).Create(&result)
	if resp.Error != nil {
		return dtos.InternalErrMS[[]T](resp.Error.Error()), resp.Error
	}
	return dtos.SuccessS[[]T](result, resp.RowsAffected), nil
}

func DbUpdateMany[T any](u *gorm.DB, ctx context.Context, filter, updateDto any, options *Opt) (dtos.GResp[int64], error) {
	if isEmptyStruct(filter) {
		return dtos.BadReqM[int64](respC.EmptyFilter.Msg()), errors.New(respC.EmptyFilter.Msg())
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
		//================.  Authorization Queries. ==========
		if options.AuthKey != nil && options.AuthVal != nil {
			queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
			query = query.Where(queryStr, options.AuthVal)
		}

	}
	resp := query.Clauses(clause.Returning{}).Model(&result).Where(filter).Updates(updateDto)
	if resp.Error != nil {
		return dtos.InternalErrMS[int64](resp.Error.Error()), resp.Error
	}
	if resp.RowsAffected == 0 {
		return dtos.NotModifiedErr[int64](respC.NoRecordsUpdated.Msg()), errors.New(respC.NoRecordsUpdated.Msg())
	}
	return dtos.SuccessCS[int64](resp.RowsAffected, respC.UpdateSuccess, resp.RowsAffected), nil
}

func DbDeleteMany[T any](u *gorm.DB, ctx context.Context, filter any, options *Opt) (dtos.GResp[int64], error) {
	if isEmptyStruct(filter) {
		return dtos.BadReqM[int64](respC.EmptyFilter.Msg()), errors.New(respC.EmptyFilter.Msg())
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
		//================.  Authorization Queries. ==========
		if options.AuthKey != nil && options.AuthVal != nil {
			queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
			query = query.Where(queryStr, options.AuthVal)
		}

	}
	resp := query.Model(result).Where(util.UnwrapAny(filter)).Delete(&result)

	if resp.Error != nil {
		if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
			return dtos.NotFoundErrS[int64](resp.Error.Error()), resp.Error
		}
		//cmn.LogTrace("the error is", resp)
		return dtos.InternalErrMS[int64](resp.Error.Error()), resp.Error
	}
	if resp.RowsAffected == 0 {
		return dtos.NotModifiedErr[int64]("NO records were deleted"), errors.New("NO records were deleted")
	}
	return dtos.SuccessCS[int64](resp.RowsAffected, respC.DeleteSuccess, resp.RowsAffected), nil

}

//--------- New Methods

func DbUpsertManyWithValues[T any](u *gorm.DB, ctx context.Context, value any, options Opt) (dtos.GResp[[]T], error) {
	var result []T
	if err := mapstructure.Decode(value, &result); err != nil {
		return dtos.InternalErrMS[[]T](err.Error()), err
	}

	query := u.WithContext(ctx)
	if options.Debug {
		query = query.Debug()
	}
	//================.  Authorization Queries. ==========
	if options.AuthKey != nil && options.AuthVal != nil {
		queryStr := fmt.Sprintf("%s = ?", *options.AuthKey)
		query = query.Where(queryStr, options.AuthVal)
	}

	for _, pre := range options.Preloads {
		query = query.Preload(pre)
	}
	resp := query.Clauses(clause.OnConflict{
		Columns:   options.Columns,
		DoUpdates: clause.AssignmentColumns(options.UpdateColumns),
	}, clause.Returning{}).Model(&result).Create(&result)
	if resp.Error != nil {
		return dtos.InternalErrMS[[]T](resp.Error.Error()), resp.Error
	}
	return dtos.SuccessS[[]T](result, resp.RowsAffected), nil
}
