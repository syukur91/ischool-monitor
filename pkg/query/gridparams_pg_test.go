package query

import (
	"testing"
)

func TestGridParam_FullQuery_Simple(t *testing.T) {

	gridParams := &GridParams{Take: 10, Page: 1, Skip: 0, PageSize: 10}

	gridParams.HasFilter = true
	gridParams.Filter.Logic = "and"

	gf := &GridFilter{}
	gf.Field = "name"
	gf.Operator = "contains"
	gf.Value = "API"

	gf2 := &GridFilter{}
	gf2.Field = "type"
	gf2.Operator = "eq"
	gf2.Value = "API"

	gridParams.Filter.Filters = append(gridParams.Filter.Filters, *gf)
	gridParams.Filter.Filters = append(gridParams.Filter.Filters, *gf2)

	gridParams.HasSort = true
	gridParams.Sort = []GridSort{
		GridSort{
			Field: "name",
			Dir:   "asc",
		},
	}

	preParams := []interface{}{"tenant"}
	preQuery := "tenant_id = ?"

	query, params := FullQuery(gridParams, preQuery, preParams)

	expectedQuery := " WHERE tenant_id = $1 AND name LIKE $2 AND type = $3 ORDER BY name ASC OFFSET 0 LIMIT 10;"
	if query != expectedQuery {
		t.Errorf("expect query %s, but got %s", expectedQuery, query)
		return
	}

	if len(params) != 3 {
		t.Errorf("expect parameters length %d, but got %d", 3, len(params))
		return
	}
}
