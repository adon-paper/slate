package arango

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/arangodb/go-driver"
)

func (r *ArangoBaseRepository) Where(column string, operator string, value interface{}) *ArangoBaseRepository {
	r.ArangoQuery = *r.ArangoQuery.Where(column, operator, value)
	return r
}

func (r *ArangoBaseRepository) WhereOr(column string, operator string, value interface{}) *ArangoBaseRepository {
	r.ArangoQuery = *r.ArangoQuery.WhereOr(column, operator, value)
	return r
}

func (r *ArangoBaseRepository) WhereColumn(column string, operator string, value string) *ArangoBaseRepository {
	r.ArangoQuery = *r.ArangoQuery.WhereColumn(column, operator, value)
	return r
}

func (r *ArangoBaseRepository) Join(from, fromKey, To, toKey string) *ArangoBaseRepository {
	r.ArangoQuery = *r.ArangoQuery.Join(from, fromKey, To, toKey)
	return r
}

func (r *ArangoBaseRepository) WithOne(repo *ArangoQuery, alias string) *ArangoBaseRepository {
	r.ArangoQuery = *r.ArangoQuery.WithOne(repo, alias)
	return r
}

func (r *ArangoBaseRepository) WithMany(repo *ArangoQuery, alias string) *ArangoBaseRepository {
	r.ArangoQuery = *r.ArangoQuery.WithMany(repo, alias)
	return r
}

func (r *ArangoBaseRepository) JoinEdge(from, fromKey, edge, alias, direction string) *ArangoBaseRepository {
	r.ArangoQuery = *r.ArangoQuery.JoinEdge(from, fromKey, edge, alias, direction)
	return r
}

func (r *ArangoBaseRepository) Offset(offset int) *ArangoBaseRepository {
	r.ArangoQuery = *r.ArangoQuery.Offset(offset)
	return r
}

func (r *ArangoBaseRepository) Limit(limit int) *ArangoBaseRepository {
	r.ArangoQuery = *r.ArangoQuery.Limit(limit)
	return r
}

func (r *ArangoBaseRepository) Sort(sortField, sortOrder string) *ArangoBaseRepository {
	r.ArangoQuery = *r.ArangoQuery.Sort(sortField, sortOrder)
	return r
}

func (r *ArangoBaseRepository) Raw() (string, map[string]interface{}) {
	return r.ArangoQuery.Raw()
}

func (r *ArangoBaseRepository) clearQuery() {
	r.ArangoQuery.clearQuery()
	r.collection = r.Collection
}

func (r *ArangoBaseRepository) Get(request interface{}) error {

	r.query, r.filterArgs = r.Raw()

	return r.executeQuery(request)
}

func (r *ArangoBaseRepository) Count(request interface{}) error {
	var (
		returnData string
		limitQuery string
		sortQuery  string
	)

	returnData = "COLLECT WITH COUNT INTO total RETURN total"

	r.query = fmt.Sprintf("FOR %s in %s %s %s %s RETURN %s",
		r.collection,
		r.collection,
		r.query,
		limitQuery,
		sortQuery,
		returnData,
	)

	return r.executeQuery(request)
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
