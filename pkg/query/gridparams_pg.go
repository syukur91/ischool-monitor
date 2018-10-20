package query

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Usage:
// query, params := FullQuery(gridParams, preQuery, preParams)
// statement := "select * from table " + query
// exec(statement, params)

// FilterClause is
func FilterClause(l GridFilter) (string, []interface{}) {
	query := ""
	var params []interface{}

	if l.HasSubFilter {
		i := 0
		query += "( "
		for _, v := range l.Filters {
			subQuery, subParams := FilterClause(v)
			query += subQuery
			params = append(params, subParams)
			if i != len(l.Filters)-1 {
				query += strings.ToUpper(l.Logic) + " "
			}
			i++
		}
		query += ") "
		return query, params
	}

	query += l.Field + " "
	query += pgOperatorMap[l.Operator].Operator + " "
	if !pgOperatorMap[l.Operator].Unary {
		// add placeholder to where cluase
		query += "? "
		// add actual value to params
		value := ""
		if pgOperatorMap[l.Operator].WildcardBefore {
			value += "%"
		}
		if !pgOperatorMap[l.Operator].Unary {
			value += l.Value.(string)
		}
		if pgOperatorMap[l.Operator].WildcardAfter {
			value += "%"
		}
		params = append(params, value)
	}

	return query, params
}

// FullQuery is
func FullQuery(l *GridParams, preQuery string, preParams []interface{}) (string, []interface{}) {

	query := " "

	// Build WHERE clause
	if len(preQuery) > 0 {
		query += "WHERE " + preQuery + " "
		if l.HasFilter {
			query += "AND "
		}
	} else {
		if l.HasFilter {
			query += "WHERE "
		}
	}

	// Build Filter clause
	i := 0
	for _, v := range l.Filter.Filters {
		subWhere, subParams := FilterClause(v)
		query += subWhere
		preParams = append(preParams, subParams...)
		if i != len(l.Filter.Filters)-1 {
			query += strings.ToUpper(l.Filter.Logic) + " "
		}
		i++
	}

	// Build Sort clause
	sort := ""
	for i, v := range l.Sort {
		if i > 0 {
			sort += ", "
		}

		sort += v.Field + " " + strings.ToUpper(v.Dir)

		if i == len(l.Sort)-1 {
			sort += " "
		}
	}

	if l.HasSort {
		query += "ORDER BY " + sort
	}

	// Build Pagination clause
	query += " OFFSET " + strconv.Itoa(l.Skip) + " "
	query += " LIMIT " + strconv.Itoa(l.PageSize)

	return sqlx.Rebind(sqlx.BindType("postgres"), strings.Replace(query, "  ", " ", -1)), preParams
}

// FilterQuery is
func FilterQuery(l *GridParams, preQuery string, preParams []interface{}) (string, []interface{}) {
	query := " "

	// Build WHERE clause
	if len(preQuery) > 0 {
		query += "WHERE " + preQuery + " "
		if l.HasFilter {
			query += "AND "
		}
	} else {
		if l.HasFilter {
			query += "WHERE "
		}
	}

	// Build Filter clause
	i := 0
	for _, v := range l.Filter.Filters {
		subWhere, subParams := FilterClause(v)
		query += subWhere
		preParams = append(preParams, subParams...)
		if i != len(l.Filter.Filters)-1 {
			query += strings.ToUpper(l.Filter.Logic) + " "
		}
		i++
	}

	return sqlx.Rebind(sqlx.BindType("postgres"), strings.Replace(query, "  ", " ", -1)), preParams
}

// SortQuery is
func SortQuery(query string, l GridParams) string {
	sort := ""
	for i, v := range l.Sort {
		if i > 0 {
			sort += ", "
		}

		sort += v.Field + " " + strings.ToUpper(v.Dir)

		if i == len(l.Sort)-1 {
			sort += " "
		}
	}

	if l.HasSort {
		query += "ORDER BY " + sort
	}

	return query
}

// SortPagingQuery is
func SortPagingQuery(query string, l GridParams) string {
	sort := ""
	for i, v := range l.Sort {
		if i > 0 {
			sort += ", "
		}

		sort += v.Field + " " + strings.ToUpper(v.Dir)

		if i == len(l.Sort)-1 {
			sort += " "
		}
	}

	if l.HasSort {
		query += "ORDER BY " + sort
	}

	query += " OFFSET " + strconv.Itoa(l.Skip) + " "
	query += " LIMIT " + strconv.Itoa(l.PageSize)

	return query
}

func makeOperation(op string, field string) string {
	placeholder := ""
	if operator, ok := pgOperatorMap[op]; ok {
		op = operator.Operator
		if !operator.Unary {
			placeholder = " ?"
		}
	}
	return field + fmt.Sprintf(" %s%s", op, placeholder)
}

func makeOperand(op string, value string, array bool) string {
	if operator, ok := pgOperatorMap[op]; ok {
		if operator.Unary {
			return ""
		}
		if array {
			return fmt.Sprintf("{%s}", value)
		}
	}

	return value
}

func isArray(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Slice:
		return true
	case reflect.Array:
		return true
	default:
		return false
	}
}

var pgOperatorMap = map[string]ComparisonOperator{
	"eq": ComparisonOperator{
		Operator: "=",
		Unary:    false,
	},
	"neq": ComparisonOperator{
		Operator: "<>",
		Unary:    false,
	},
	"contains": ComparisonOperator{
		Operator:       "LIKE",
		WildcardBefore: true,
		WildcardAfter:  true,
		Unary:          false,
	},
	"doesnotcontain": ComparisonOperator{
		Operator:       "NOT LIKE",
		WildcardBefore: true,
		WildcardAfter:  true,
		Unary:          false,
	},
	"startswith": ComparisonOperator{
		Operator:      "LIKE",
		WildcardAfter: true,
		Unary:         false,
	},
	"endswith": ComparisonOperator{
		Operator:       "LIKE",
		Unary:          false,
		WildcardBefore: true,
	},
	"isnull": ComparisonOperator{
		Operator: "IS NULL",
		Unary:    true,
	},
	"isnotnull": ComparisonOperator{
		Operator: "IS NOT NULL",
		Unary:    true,
	},
	"isempty": ComparisonOperator{
		Operator: "= ''",
		Unary:    true,
	},
	"isnotempty": ComparisonOperator{
		Operator: "<> ''",
		Unary:    true,
	},
}
