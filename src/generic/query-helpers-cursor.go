package generic

import (
	"gorm.io/gorm"

	"github.com/birukbelay/gocmn/src/dtos"
	"github.com/birukbelay/gocmn/src/util"
)

func CreateCursor[T any](limit int, forward bool, dataList []T, orderedBy, pagiCursor string) (cursorObj dtos.CursorStruct, dataLst []T, eror error) {
	//cursor := ""
	cursor := dtos.CursorStruct{
		HasMore: false,
		HasPrev: false,
	}
	/*
			if i have a fwd cursor: means i previously have data,
					- so i will generate Back Cursor
					- and a fwd cursor if the Data has more
		   if cursor is bwd: means i have next page data
					- s i will generate a fwd Cursor
					- and a bwd cursor if the data has more
	*/
	if forward {

		if len(dataList) > limit {
			cursor.HasMore = true
			//pop the last element
			dataList = dataList[:len(dataList)-1]
			//generate a cursor
			lastElement := dataList[len(dataList)-1]
			nextCursor, err := util.GenerateCursor(orderedBy, lastElement, true)
			if err != nil {
				return dtos.CursorStruct{}, nil, err
			}
			cursor.NextCursor = nextCursor

		}
		if len(dataList) > 0 && pagiCursor != "" {
			cursor.HasPrev = true
			prevCursor, err := util.GenerateCursor(orderedBy, dataList[0], false)
			if err != nil {
				return dtos.CursorStruct{}, nil, err
			}
			cursor.PrevCursor = prevCursor
		}

	} else { //it is a backward cursor

		//PrevCursor: if the data is bigger than limit, it has prevPage
		if len(dataList) > limit {
			cursor.HasPrev = true
			//pop the last element
			dataList = dataList[:len(dataList)-1]
			//generate a cursor
			//lastElement := dataList[len(dataList)-1]
			prevCursor, err := util.GenerateCursor(orderedBy, dataList[len(dataList)-1], false)
			if err != nil {
				return dtos.CursorStruct{}, nil, err
			}
			cursor.PrevCursor = prevCursor

		}
		//The next cursor is the 0th element
		if len(dataList) > 0 && pagiCursor != "" {
			cursor.HasMore = true
			//because the array is reversed, the prev cursor will be the last item
			nextCursor, err := util.GenerateCursor(orderedBy, dataList[0], true)
			if err != nil {
				return dtos.CursorStruct{}, nil, err
			}
			cursor.NextCursor = nextCursor
		}
		dataList = util.Reverse[T](dataList)
	}
	return cursor, dataList, nil
}

func CreateCursorQuery(query *gorm.DB, curs, orderString, sort_dir string) (q *gorm.DB, ordrStr string, forwd bool, eror error) {
	forward := true
	if curs != "" {
		orderedByName, cursorValue, cursorID, fwd, err := util.DecodeCursor(curs)
		if err != nil {
			return nil, "", forward, err
		}
		forward = fwd
		var queryString string
		if forward {
			//orderString looks like `created_at desc, id`
			//queryString` is stg like: "(created_at > ?) OR (created_at = ? AND id > ?)"
			orderString, queryString = util.GenerateFwdQueryString(*orderedByName, sort_dir)
		} else {
			orderString, queryString = util.GeneratePrevQueryString(*orderedByName, sort_dir)
		}

		if cursorValue != nil && cursorID != nil {
			//`queryString` is stg like: "(%s > ?) OR (%s = ? AND id > ?)"
			//the `cursorValue` is the Value of the field used for sorting,
			// CursorID is the ids value,
			query = query.Where(queryString, *cursorValue, *cursorValue, *cursorID)
		}

	}
	return query, orderString, forward, nil
}
