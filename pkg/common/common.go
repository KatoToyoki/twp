package common

import (
	"errors"
	"regexp"

	"github.com/jykuo-love-shiritori/twp/pkg/constants"
)

type QueryParams struct {
	Offset int32 `query:"offset"` // the start index of the query
	Limit  int32 `query:"limit"`  // the number of the query
}

func NewQueryParams(offset int32, limit int32) QueryParams {
	return QueryParams{Offset: offset, Limit: limit}
}

func (q *QueryParams) Validate() error {
	if q.Offset < 0 || q.Limit < 0 || q.Limit > constants.QUERY_LIMIT {
		return errors.New("invalid query parameter")
	}
	return nil
}

func HasSpecialChars(input string) bool {
	regexPattern := `[.*+?()|{}\\^$]`
	re := regexp.MustCompile(regexPattern)
	return re.MatchString(input)
}
