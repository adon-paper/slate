package arango

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"slate/constant"
	"slate/helper"
	"strconv"
	"strings"
	"time"

	"github.com/arangodb/go-driver"
)

type ArangoBaseRepositoryInterface interface {
	first(c context.Context, request ArangoInterface) error
	rawFirst(c context.Context, queryBuilder ArangoQueryBuilder, request ArangoInterface) error
	all(c context.Context, request interface{}, baseFilter PaginationFilters) ([]interface{}, int64, error)
	buildFilter(s interface{}, filters []ArangoFilterQueryBuilder, joinCollection string, prefixes ...string) []ArangoFilterQueryBuilder

	create(c context.Context, request ArangoInterface) error
	update(c context.Context, request ArangoInterface) error
	delete(c context.Context, request ArangoInterface) error
}

type arangoBaseRepository struct {
	logger   *log.Logger
	arangoDB ArangoDB

	collection string
}

func NewArangoBaseRepository(logger *log.Logger, arangoDB ArangoDB) ArangoBaseRepositoryInterface {
	return &arangoBaseRepository{
		logger:   logger,
		arangoDB: arangoDB,
	}
}

func (r *arangoBaseRepository) parseFilterToQuery(queryBuilder ArangoQueryBuilder) (string, map[string]interface{}) {
	var filterQuery string
	filterArgs := make(map[string]interface{})

	for index, filter := range queryBuilder.Filters {
		keyword := "FILTER"
		if index != 0 && (filter.AndOr == "AND" || filter.AndOr == "OR") {
			keyword = filter.AndOr
		}

		if filter.CustomFilter == "" {
			if filter.ArgumentKey == "" {
				filter.ArgumentKey = strings.ReplaceAll(filter.Key, ".", "_")
			}

			if filter.Operator == "" {
				filter.Operator = "=="
			} else if filter.Operator == "LIKE" {
				filter.Value = "%" + filter.Value.(string) + "%"
			}

			filterQuery += " " + keyword + " " + filter.Key + ` ` + filter.Operator + ` @` + filter.ArgumentKey
			filterArgs[filter.ArgumentKey] = filter.Value
		} else {
			filterQuery += " " + keyword + " " + filter.CustomFilter
			filterArgs[filter.ArgumentKey] = filter.Value
		}
	}

	return filterQuery, filterArgs
}

func (r *arangoBaseRepository) parseJoinToQuery(queryBuilder ArangoQueryBuilder) (string, string) {
	var joinQuery string
	resultQuery := "data"

	if len(queryBuilder.Joins) > 0 {
		joinQuery += " FILTER data != null "

		for index, join := range queryBuilder.Joins {
			if join.CollectionFrom != r.collection {
				join.FromKey = "data_" + join.CollectionFrom + "." + join.FromKey
			} else {
				join.FromKey = "data" + "." + join.FromKey
			}

			joinQuery += " FOR data_" + join.CollectionTo + " in " + join.CollectionTo +
				" FILTER data_" + join.CollectionTo + "." + join.ToKey + " == " + join.FromKey

			if join.ResultKey == "" {
				join.ResultKey = join.CollectionTo
			}

			if index == 0 {

				resultQuery = `{
				` + r.collection + `: data,
				` + join.ResultKey + ": data_" + join.CollectionTo
			} else {

				resultQuery += `,
			` + join.ResultKey + ": data_" + join.CollectionTo
			}
		}

		resultQuery += `
			}
		`
	}

	return joinQuery, resultQuery
}

func (r *arangoBaseRepository) buildQuery(queryBuilder ArangoQueryBuilder) (string, string, map[string]interface{}) {

	filterQuery, filterArgs := r.parseFilterToQuery(queryBuilder)
	joinQuery, resultQuery := r.parseJoinToQuery(queryBuilder)

	var sortOrder string
	if queryBuilder.SortOrder > 0 {
		sortOrder = "ASC"
	} else {
		sortOrder = "DESC"
	}

	if queryBuilder.SortField != "" {
		queryBuilder.SortField = "created_at"
	}

	totalRecordsQuery := `
		FOR data IN ` + r.collection +
		joinQuery + " " + filterQuery + `
		COLLECT WITH COUNT INTO length
		RETURN length
	`

	var query string
	if queryBuilder.Rows != 0 {

		query = `
			FOR data IN ` + r.collection +
			joinQuery + " " + filterQuery +
			" LIMIT " + strconv.Itoa(queryBuilder.First) + ", " + strconv.Itoa(queryBuilder.Rows) + `
			SORT data.` + queryBuilder.SortField + ` ` + sortOrder + `
			RETURN ` + resultQuery

	} else {
		query = `
		FOR data IN ` + r.collection +
			joinQuery + " " + filterQuery + ` 
		SORT data.` + queryBuilder.SortField + ` ` + sortOrder + `
		RETURN ` + resultQuery
	}

	return totalRecordsQuery, query, filterArgs
}

func (r *arangoBaseRepository) buildFilter(s interface{}, filters []ArangoFilterQueryBuilder, joinCollection string, prefixes ...string) []ArangoFilterQueryBuilder {
	v := reflect.Indirect(reflect.ValueOf(s))

	var collectionPrefix string
	if joinCollection == "" || joinCollection == r.collection {
		collectionPrefix = "data."
	} else {
		collectionPrefix = "data_" + joinCollection + "."
	}

	var prefix string
	if len(prefixes) > 0 && prefixes[0] != "" {
		prefix = prefixes[0] + "."
	}

	if v.Kind() == reflect.Slice {
		if len(s.([]interface{})) > 0 {
			fmt.Println(s)
		}

	} else {
		fmt.Println(v.NumField())
		for i := 0; i < v.NumField(); i++ {
			tags := strings.Split(v.Type().Field(i).Tag.Get("json"), ",")
			value := v.Field(i).Interface()
			if (strings.Contains(v.Type().Field(i).Type.String(), "") || strings.Contains(v.Type().Field(i).Type.String(), "dto.")) && !strings.Contains(v.Type().Field(i).Type.String(), "Interface") {
				var tag string
				if collection := joinCollection + v.Type().Field(i).Tag.Get("collection"); collection != "" || tags[0] == r.collection {
					tag = ""
				} else {
					tag = tags[0]
				}

				filters = r.buildFilter(value, filters, joinCollection+v.Type().Field(i).Tag.Get("collection"), tag)
			} else {
				if !helper.Empty(value) {
					var filter ArangoFilterQueryBuilder
					if filterKey := v.Type().Field(i).Tag.Get("filter"); filterKey != "" {
						if v.Type().Field(i).Type.String()[:2] == "[]" {
							filter.Key = collectionPrefix + prefix + filterKey
							filter.Value = v.Field(i).Interface()
							filter.Operator = "IN"
							filters = append(filters, filter)
						}
					} else {
						filter.Key = collectionPrefix + prefix + tags[0]
						filter.Value = value
						if tags[0] == "created_at" || tags[0] == "updated_at" {
							filter.Operator = "LIKE"
							filter.Value = value.(time.Time).Format("2006-01-02")
						}
						filters = append(filters, filter)
					}
				}
			}
		}
	}

	return filters
}

func (r *arangoBaseRepository) first(c context.Context, request ArangoInterface) error {
	queryBuilder := ArangoQueryBuilder{
		Filters: r.buildFilter(request, []ArangoFilterQueryBuilder{}, ""),
	}

	return r.rawFirst(c, queryBuilder, request)
}

func (r *arangoBaseRepository) rawFirst(c context.Context, queryBuilder ArangoQueryBuilder, request ArangoInterface) error {

	ctx := driver.WithQueryCount(c)

	_, query, condition := r.buildQuery(queryBuilder)

	fmt.Println(query)
	fmt.Println(condition)

	data, err := r.arangoDB.DB().Query(ctx, query, condition)
	if err != nil {
		return err
	}

	if data.Count() > 0 {
		data.ReadDocument(c, &request)
		return nil
	}

	return constant.ErrorNotFound
}

func (r *arangoBaseRepository) all(c context.Context, request interface{}, baseFilter PaginationFilters) ([]interface{}, int64, error) {
	queryBuilder := ArangoQueryBuilder{
		Filters:   r.buildFilter(request, []ArangoFilterQueryBuilder{}, ""),
		First:     baseFilter.First,
		Rows:      baseFilter.Rows,
		SortField: baseFilter.SortField,
		SortOrder: baseFilter.SortOrder,
	}

	return r.rawAll(c, queryBuilder)
}

func (r *arangoBaseRepository) rawAll(c context.Context, queryBuilder ArangoQueryBuilder) ([]interface{}, int64, error) {
	var response []interface{}
	ctx := driver.WithQueryCount(c)

	totalRecordsQuery, query, condition := r.buildQuery(queryBuilder)

	data, err := r.arangoDB.DB().Query(ctx, query, condition)
	if err != nil {
		return response, 0, err
	}

	if data.Count() > 0 {
		for data.HasMore() {
			var request interface{}
			data.ReadDocument(c, &request)
			response = append(response, request)
		}
	}

	countData, err := r.arangoDB.DB().Query(ctx, totalRecordsQuery, condition)
	if err != nil {
		return response, 0, err
	}

	var totalRecords int64
	if countData.Count() > 0 {

		_, err = countData.ReadDocument(ctx, &totalRecords)
		if err != nil {
			return response, 0, err
		}
	}

	return response, totalRecords, nil
}

func (r *arangoBaseRepository) create(c context.Context, request ArangoInterface) error {
	collection, err := r.arangoDB.DB().Collection(c, r.collection)
	if err != nil {
		return err
	}

	request.InitializeTimestamp()

	insert, err := collection.CreateDocument(c, request)
	if err != nil {
		return err
	}

	request.Set(
		insert.ID.String(),
		insert.Key,
		insert.Rev,
	)
	return nil
}

func (r *arangoBaseRepository) update(c context.Context, request ArangoInterface) error {
	collection, err := r.arangoDB.DB().Collection(c, r.collection)
	if err != nil {
		return err
	}

	request.UpdateTimestamp()

	_, err = collection.UpdateDocument(c, request.GetKey(), request)
	if err != nil {
		return err
	}
	return nil
}

func (r *arangoBaseRepository) delete(c context.Context, request ArangoInterface) error {
	collection, err := r.arangoDB.DB().Collection(c, r.collection)
	if err != nil {
		return err
	}

	_, err = collection.RemoveDocument(c, request.GetKey())
	if err != nil {
		return err
	}
	return nil
}
