package generic

import (
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/birukbelay/gocmn/src/dtos"
)

func addStartEndDate(query *gorm.DB, startDate, endDate string) (*gorm.DB, error) {
	//====================   Start Date & END Date =======================
	if startDate != "" {
		startDate, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format: %v", err)
		}
		query = query.Where("created_at >= ?", startDate)
	}
	// Add EndDate filter if it exists and is non-empty
	if endDate != "" {
		endDate, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %v", err)
		}
		// Set endDate to end of day (23:59:59) to include the full day
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		query = query.Where("created_at <= ?", endDate)
	}
	return query, nil

}
func addSearchFilters(query *gorm.DB, pagi dtos.PaginationInput) (q *gorm.DB) {
	//=================================  Prefix search Matching ===============
	//for indexed name matching, eg: `name LIKE abe%`
	if pagi.Like != "" && pagi.PrefixColLike != "" {
		query = query.Where(fmt.Sprintf("%s LIKE ?", pagi.PrefixColLike), pagi.Like+"%")
	}
	//=============  For text search, WARNING! YOU will need to index the cols for text serarch ======
	if pagi.Query != "" && len(pagi.TxtSearchCols) > 0 {
		searchVector := fmt.Sprintf("to_tsvector('english', %s)", strings.Join(pagi.TxtSearchCols, " || ' ' || "))
		query = query.Where(searchVector+" @@ to_tsquery(?)", ToTsQuery(pagi.Query))
	}
	return query
}
func ToTsQuery(input string) string {
	// Remove dangerous characters
	cleaned := strings.Map(func(r rune) rune {
		switch r {
		case ':', '!', '&', '|', '(', ')', '<', '>', '\'', '\\':
			return -1 // remove
		}
		return r
	}, input)

	// Split into words
	words := strings.Fields(cleaned)
	if len(words) == 0 {
		return ""
	}

	// Add prefix operator
	for i, w := range words {
		words[i] = w + ":*"
	}

	// Join with AND (you can use OR: "|")
	return strings.Join(words, " & ")
}
func searchTags(query *gorm.DB, tags []string) (q *gorm.DB) {
	if len(tags) > 0 {
		query = query.Where("tag IN (?)", tags)
	}
	return query
}
func addAnyOfTags(query *gorm.DB, tags []string) (q *gorm.DB) {
	if len(tags) > 0 {
		query = query.Where("tags && ?", pq.Array(tags))
	}
	return query
}
func addAllOfTags(query *gorm.DB, tags []string) (q *gorm.DB) {
	if len(tags) > 0 {
		query = query.Where("tags @> ?", pq.Array(tags))
	}
	return query
}

// =====================  For Preload, Debug and Select =======================
func DebugPreloadSelect(query *gorm.DB, options *Opt, sel []string) (q *gorm.DB) {
	if options != nil {
		if options.Debug {
			query = query.Debug()
		}
		for _, pre := range options.Preloads {
			query = query.Preload(pre)
		}
	}
	if len(sel) > 0 {
		query = query.Select(sel)
	}
	return query
}
func AddInQueries(query *gorm.DB, options *Opt) (q *gorm.DB) {
	if options != nil {
		if len(options.InQueries) > 0 {
			for n, v := range options.InQueries {
				if len(v) > 0 {
					query = query.Where(fmt.Sprintf("%s IN (?)", n), v)
				}
			}
		}
	}

	return query
}
func AddNotInQueries(query *gorm.DB, options *Opt) (q *gorm.DB) {
	if options != nil {
		if len(options.NotInQueries) > 0 {
			for n, v := range options.NotInQueries {
				if len(v) > 0 {
					query = query.Where(fmt.Sprintf("%s NOT IN (?)", n), v)
				}
			}
		}
	}

	return query
}
