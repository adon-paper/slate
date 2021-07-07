package arango

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/noldwidjaja/slate/helper"
)

type (
	traversalDirection string
)

const (
	INBOUND  traversalDirection = "INBOUND"
	OUTBOUND traversalDirection = "OUTBOUND"
	ANY      traversalDirection = "ANY"
)

type arangoQueryTraversal struct {
	enabled   bool
	direction traversalDirection
	sourceId  string
}

type ArangoQuery struct {
	collection string
	traversal  arangoQueryTraversal
	query      string
	filterArgs map[string]interface{}
	joins      []*ArangoQuery
	withs      []*ArangoQuery
	sortField  string
	sortOrder  string
	offset     int
	limit      int
	first      bool
	alias      string
}

func SubQuery(collection string) *ArangoQuery {
	return &ArangoQuery{
		collection: collection,
		alias:      collection,
	}
}

func (r *ArangoQuery) Where(column string, operator string, value interface{}) *ArangoQuery {
	if strings.Contains(column, ".") {
		r.query += " FILTER " + column + " " + operator + " @" + r.collection + "_" + column
	} else {
		r.query += " FILTER " + r.collection + "." + column + " " + operator + " @" + r.collection + "_" + column
	}

	if r.filterArgs == nil {
		r.filterArgs = make(map[string]interface{})
	}

	if strings.ToUpper(operator) == "LIKE" {
		value = "%" + value.(string) + "%"
	}

	r.filterArgs[r.collection+"_"+column] = value

	return r
}

func (r *ArangoQuery) WhereOr(column string, operator string, value interface{}) *ArangoQuery {
	if strings.Contains(column, ".") {
		r.query += " OR " + column + " " + operator + " @" + r.collection + "_" + column
	} else {
		r.query += " OR " + r.collection + "." + column + " " + operator + " @" + r.collection + "_" + column
	}

	if r.filterArgs == nil {
		r.filterArgs = make(map[string]interface{})
	}

	if strings.ToUpper(operator) == "LIKE" {
		value = "%" + value.(string) + "%"
	}

	r.filterArgs[r.collection+"_"+column] = value

	return r
}

func (r *ArangoQuery) WhereColumn(column string, operator string, value string) *ArangoQuery {
	if strings.Contains(column, ".") {
		r.query += " FILTER " + column + " " + operator + " " + value
	} else {
		r.query += " FILTER " + r.collection + "." + column + " " + operator + " " + value
	}
	return r
}

func (r *ArangoQuery) Join(query *ArangoQuery) *ArangoQuery {
	q, f := query.toQueryWithoutReturn()
	r.query += q

	r.joins = append(r.joins, query)
	r.filterArgs = helper.MergeMaps(r.filterArgs, f)
	return r
}

func (r *ArangoQuery) WithOne(repo *ArangoQuery, alias string) *ArangoQuery {
	r.first = true
	r.with(repo, alias)
	return r
}

func (r *ArangoQuery) WithMany(repo *ArangoQuery, alias string) *ArangoQuery {
	r.first = false
	r.with(repo, alias)
	return r
}

func (r *ArangoQuery) with(query *ArangoQuery, alias string) *ArangoQuery {
	query.alias = alias
	q, f := query.ToQuery()
	r.query += ` LET ` + alias + ` = ( 
      ` + q + ` 
      )
   `

	r.withs = append(r.withs, query)
	r.filterArgs = helper.MergeMaps(r.filterArgs, f)
	return r
}

func (r *ArangoQuery) Offset(offset int) *ArangoQuery {
	r.offset = offset
	return r
}

func (r *ArangoQuery) Limit(limit int) *ArangoQuery {
	r.limit = limit

	return r
}

func (r *ArangoQuery) Sort(sortField, sortOrder string) *ArangoQuery {
	if strings.Contains(sortField, ".") {
		r.sortField = sortField
	} else {
		r.sortField = r.collection + "." + sortField
	}

	if sortOrder != "ASC" {
		r.sortOrder = "DESC"
	} else {
		r.sortOrder = "ASC"
	}

	return r
}

func (r *ArangoQuery) Traversal(source string, direction traversalDirection) *ArangoQuery {
	r.traversal.enabled = true
	r.traversal.direction = direction
	r.traversal.sourceId = source

	return r
}

func (r *ArangoQuery) toQueryWithoutReturn() (string, map[string]interface{}) {
	var (
		limitQuery string
		sortQuery  string
		finalQuery string
	)

	if r.limit > 0 {
		limitQuery = fmt.Sprintf("LIMIT %s,%s", strconv.Itoa(r.offset), strconv.Itoa(r.limit))
	}

	if r.sortField != "" {
		sortQuery = fmt.Sprintf("SORT %s %s", r.sortField, r.sortOrder)
	}

	if r.traversal.enabled {
		finalQuery = fmt.Sprintf("FOR %s in %s %s %s %s %s %s ",
			r.alias,
			r.traversal.direction,
			r.traversal.sourceId,
			r.collection,
			r.query,
			limitQuery,
			sortQuery,
		)
	} else {
		finalQuery = fmt.Sprintf("FOR %s in %s %s %s %s ",
			r.alias,
			r.collection,
			r.query,
			limitQuery,
			sortQuery,
		)
	}

	args := r.filterArgs

	return finalQuery, args
}

func (r *ArangoQuery) ToQuery() (string, map[string]interface{}) {
	var (
		returnData string
		limitQuery string
		sortQuery  string
		finalQuery string
	)

	returnData = "MERGE("

	if len(r.withs) > 0 {
		returnData += "{"
		for index, with := range r.withs {
			alias := with.alias
			if with.first {
				alias = fmt.Sprintf(" FIRST(%s) ", alias)
			}

			if index == 0 {
				returnData += fmt.Sprintf("%s :%s", with.alias, alias)
			} else {
				returnData += fmt.Sprintf(", %s :%s", with.alias, alias)
			}
		}
		returnData += "}, "
	}

	if len(r.joins) > 0 {
		for _, join := range r.joins {
			returnData += fmt.Sprintf("%s, ", join.collection)
		}
	}

	returnData += fmt.Sprintf("%s)", r.collection)

	if r.limit > 0 {
		limitQuery = fmt.Sprintf("LIMIT %s,%s", strconv.Itoa(r.offset), strconv.Itoa(r.limit))
	}

	if r.sortField != "" {
		sortQuery = fmt.Sprintf("SORT %s %s", r.sortField, r.sortOrder)
	}

	if r.traversal.enabled {
		finalQuery = fmt.Sprintf("FOR %s in %s %s %s %s %s %s RETURN %s",
			r.alias,
			r.traversal.direction,
			r.traversal.sourceId,
			r.collection,
			r.query,
			limitQuery,
			sortQuery,
			returnData,
		)
	} else {
		finalQuery = fmt.Sprintf("FOR %s in %s %s %s %s RETURN %s",
			r.alias,
			r.collection,
			r.query,
			limitQuery,
			sortQuery,
			returnData,
		)
	}

	args := r.filterArgs

	r.clearQuery()

	return finalQuery, args
}

func (r *ArangoQuery) clearQuery() {
	r.query = ""
	r.filterArgs = make(map[string]interface{})
	r.joins = nil
	r.withs = nil
	r.sortField = ""
	r.sortOrder = ""
	r.offset = 0
	r.limit = 0
	r.alias = r.collection
	r.traversal = arangoQueryTraversal{}
}
