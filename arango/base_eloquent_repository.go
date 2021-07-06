package arango

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/arangodb/go-driver"
)

func (r *ArangoBaseRepository) Where(column string, operator string, value interface{}) *ArangoBaseRepository {
	if strings.Contains(column, ".") {
		r.query += " FILTER " + column + " " + operator + " @" + r.Collection + "_" + column
	} else {
		r.query += " FILTER " + r.Collection + "." + column + " " + operator + " @" + r.Collection + "_" + column
	}

	if r.filterArgs == nil {
		r.filterArgs = make(map[string]interface{})
	}

	if strings.ToUpper(operator) == "LIKE" {
		value = "%" + value.(string) + "%"
	}

	r.filterArgs[r.Collection+"_"+column] = value

	return r
}

func (r *ArangoBaseRepository) WhereOr(column string, operator string, value interface{}) *ArangoBaseRepository {
	if strings.Contains(column, ".") {
		r.query += " OR " + column + " " + operator + " @" + r.Collection + "_" + column
	} else {
		r.query += " OR " + r.Collection + "." + column + " " + operator + " @" + r.Collection + "_" + column
	}

	if r.filterArgs == nil {
		r.filterArgs = make(map[string]interface{})
	}

	if strings.ToUpper(operator) == "LIKE" {
		value = "%" + value.(string) + "%"
	}

	r.filterArgs[r.Collection+"_"+column] = value

	return r
}

func (r *ArangoBaseRepository) WhereColumn(column string, operator string, value string) *ArangoBaseRepository {
	if strings.Contains(column, ".") {
		r.query += " OR " + column + " " + operator + " " + value
	} else {
		r.query += " OR " + r.Collection + "." + column + " " + operator + " " + value
	}
	return r
}

func (r *ArangoBaseRepository) Join(from, fromKey, To, toKey string) *ArangoBaseRepository {
	r.query += ` FOR ` + To + ` IN ` + To + `
		FILTER ` + To + "." + toKey + "==" + from + "." + fromKey + `
	`

	r.joins = append(r.joins, To)

	return r
}

func (r *ArangoBaseRepository) With(from, fromKey, to, toKey, alias string) *ArangoBaseRepository {
	r.query += ` LET ` + alias + ` = (
		FOR ` + to + ` IN ` + to + `
		FILTER ` + to + "." + toKey + "==" + from + "." + fromKey

	// if len(filters) > 0 {
	// 	for _, filter := range filters {
	// 		// filter()
	// 	}
	// }

	r.query += `
		RETURN ` + to + `
	)`

	r.withs = append(r.withs, alias)

	return r
}

func (r *ArangoBaseRepository) JoinEdge(from, fromKey, edge, alias, direction string) *ArangoBaseRepository {
	r.query += `
		FOR ` + alias + ` IN ` + direction + " " + from + "." + fromKey + " " + edge + `
	`

	r.joins = append(r.joins, alias)

	return r
}

func (r *ArangoBaseRepository) WithEdge(from, fromKey, edge, alias, direction string) *ArangoBaseRepository {
	r.query += ` LET ` + alias + ` = (
		FOR ` + alias + ` IN ` + direction + " " + from + "." + fromKey + " " + edge + `
		RETURN ` + alias + ` 
	)`

	r.withs = append(r.withs, alias)

	return r
}

func (r *ArangoBaseRepository) Offset(offset int) *ArangoBaseRepository {
	r.offset = offset
	return r
}

func (r *ArangoBaseRepository) Limit(limit int) *ArangoBaseRepository {
	r.limit = limit

	return r
}

func (r *ArangoBaseRepository) Sort(sortField, sortOrder string) *ArangoBaseRepository {
	if strings.Contains(sortField, ".") {
		r.sortField = sortField
	} else {
		r.sortField = r.Collection + "." + sortField
	}

	if sortOrder != "ASC" {
		r.sortOrder = "DESC"
	} else {
		r.sortOrder = "ASC"
	}

	return r
}

func (r *ArangoBaseRepository) Get(request interface{}) error {

	r.query, r.filterArgs = r.Raw()

	return r.executeQuery(request)
}

func (r *ArangoBaseRepository) Raw() (string, map[string]interface{}) {
	var (
		returnData string
		limitQuery string
		sortQuery  string
	)

	returnData = "MERGE("
	for _, join := range r.joins {
		returnData += join + ", "
	}

	if len(r.withs) > 0 {
		returnData += "{"
		for index, with := range r.withs {
			if index == 0 {
				returnData += with + " :" + with
			} else {
				returnData += ", " + with + " :" + with
			}
		}
		returnData += "}, "
	}

	returnData += r.Collection + ")"

	if r.limit > 0 {
		limitQuery = "LIMIT " + strconv.Itoa(r.offset) + "," + strconv.Itoa(r.limit)
	}

	if r.sortField != "" {
		sortQuery = "SORT " + r.sortField + " " + r.sortOrder
	}

	rawQuery := fmt.Sprintf("FOR %s in %s %s %s %s RETURN %s",
		r.Collection,
		r.Collection,
		r.query,
		limitQuery,
		sortQuery,
		returnData,
	)

	args := r.filterArgs

	r.clearQuery()

	return rawQuery, args
}

func (r *ArangoBaseRepository) executeQuery(request interface{}) error {
	c := context.Background()

	ctx := driver.WithQueryCount(c)

	fmt.Println(r.query)

	data, err := r.ArangoDB.DB().Query(ctx, r.query, r.filterArgs)
	if err != nil {
		fmt.Println(err)
		return err
	}

	r.clearQuery()

	defer data.Close()

	if data.Count() > 0 {
		v := reflect.Indirect(reflect.ValueOf(request))

		if v.Kind() == reflect.Slice {
			var response []interface{}
			for data.HasMore() {
				var d interface{}
				data.ReadDocument(c, &d)
				response = append(response, d)
			}

			jsonResponse, _ := json.Marshal(response)
			json.Unmarshal(jsonResponse, &request)
			return nil
		}

		data.ReadDocument(c, &request)
		return nil
	}

	return nil
}

func (r *ArangoBaseRepository) clearQuery() {
	r.query = ""
	r.filterArgs = make(map[string]interface{})
	r.joins = []string{}
	r.withs = []string{}
	r.sortField = ""
	r.sortOrder = ""
	r.offset = 0
	r.limit = 0
}
