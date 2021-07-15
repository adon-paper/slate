package arango

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/noldwidjaja/slate/constant"
	"github.com/noldwidjaja/slate/helper"

	"github.com/arangodb/go-driver"
)

type ArangoBaseRepositoryInterface interface {
	// Get Raw DB Functions
	DB() driver.Database

	// Base Functions
	BuildFilter(s interface{}, filters []ArangoFilterQueryBuilder, joinCollection string, prefixes ...string) []ArangoFilterQueryBuilder
	RawFirst(c context.Context, queryBuilder ArangoQueryBuilder, request ArangoInterface) error
	RawAll(c context.Context, queryBuilder ArangoQueryBuilder) ([]interface{}, int64, error)
	All(c context.Context, request interface{}, baseFilter PaginationFilters) ([]interface{}, int64, error)
	First(c context.Context, request ArangoInterface) error
	Create(c context.Context, request ArangoInterface) error
	Update(c context.Context, request ArangoInterface) error
	Delete(c context.Context, request ArangoInterface) error

	// Eloquent Style
	Where(column string, operator string, value interface{}) *ArangoBaseRepository
	WhereOr(column string, operator string, value interface{}) *ArangoBaseRepository
	WhereColumn(column string, operator string, value string) *ArangoBaseRepository
	// Join(from, fromKey, To, toKey string) *ArangoBaseRepository
	// JoinEdge(from, fromKey, edge, alias, direction string) *ArangoBaseRepository
	Join(query *ArangoQuery) *ArangoBaseRepository
	WithMany(repo *ArangoQuery, alias string) *ArangoBaseRepository
	WithOne(repo *ArangoQuery, alias string) *ArangoBaseRepository
	Offset(offset int) *ArangoBaseRepository
	Limit(limit int) *ArangoBaseRepository
	Sort(sortField, sortOrder string) *ArangoBaseRepository
	Traversal(sourceId string, direction traversalDirection) *ArangoBaseRepository
	Get(request interface{}) error
	ToQuery() (string, map[string]interface{})
	Count(request interface{}) error
	executeQuery(request interface{}) error
	clearQuery()
}

type ArangoBaseRepository struct {
	ArangoDB   ArangoDB
	Collection string
	ArangoQuery
}

func NewArangoBaseRepository(arangoDB ArangoDB, collection string) ArangoBaseRepositoryInterface {
	return &ArangoBaseRepository{
		ArangoDB:   arangoDB,
		Collection: collection,
		ArangoQuery: ArangoQuery{
			collection: collection,
			alias:      collection,
		},
	}
}

func (r *ArangoBaseRepository) parseFilterToQuery(queryBuilder ArangoQueryBuilder) (string, map[string]interface{}) {
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

			filterQuery += " " + keyword + " " + filter.Key + ` ` + filter.Operator + ` @` + queryBuilder.Alias + filter.ArgumentKey
			filterArgs[filter.ArgumentKey+queryBuilder.Alias] = filter.Value
		} else {
			filterQuery += " " + keyword + " " + filter.CustomFilter
			filterArgs[filter.ArgumentKey+queryBuilder.Alias] = filter.Value
		}
	}

	return filterQuery, filterArgs
}

func (r *ArangoBaseRepository) parseJoinToQuery(queryBuilder ArangoQueryBuilder) (string, string) {
	var joinQuery, resultQuery string

	if len(queryBuilder.Joins) > 0 {
		joinQuery += " FILTER data != null "

		for index, join := range queryBuilder.Joins {
			if join.CollectionFrom != r.Collection {
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

				resultQuery = `
				` + r.Collection + `: data,
				` + join.ResultKey + ": data_" + join.CollectionTo
			} else {

				resultQuery += `,
			` + join.ResultKey + ": data_" + join.CollectionTo
			}
		}
	}

	return joinQuery, resultQuery
}

func (r *ArangoBaseRepository) parseWithToQuery(queryBuilder ArangoQueryBuilder) (string, string, map[string]interface{}) {
	var withQuery, resultQuery, query string
	var filterArgs map[string]interface{}
	if len(queryBuilder.With) > 0 {

		for index, with := range queryBuilder.With {
			withQuery += " LET " + with.Alias + " =( "
			if with.Alias == "" {
				with.Alias = with.Collection + string(rune(index))
			}
			_, query, filterArgs = r.buildQuery(with)
			withQuery += " " + query + " ) "
			resultQuery += with.Alias + ":" + with.Alias
			if index != len(queryBuilder.With)-1 {
				resultQuery += ","
			}
		}
	}

	return withQuery, resultQuery, filterArgs
}

func (r *ArangoBaseRepository) buildQuery(queryBuilder ArangoQueryBuilder) (string, string, map[string]interface{}) {

	filterQuery, filterArgs := r.parseFilterToQuery(queryBuilder)
	joinQuery, joinQueryResultQuery := r.parseJoinToQuery(queryBuilder)
	withQuery, withQueryResultQuery, withFilterArgs := r.parseWithToQuery(queryBuilder)

	alias := "data"
	if queryBuilder.Alias != "" {
		alias = queryBuilder.Alias
	}

	resultQuery := alias
	if joinQueryResultQuery != "" || withQueryResultQuery != "" {
		if joinQueryResultQuery != "" && withQueryResultQuery != "" {
			joinQueryResultQuery += ","
		}
		resultQuery = " { " + alias + "," + joinQueryResultQuery + withQueryResultQuery + " } "
	}

	for index, withFilterArg := range withFilterArgs {
		filterArgs[index] = withFilterArg
	}

	var sortOrder, sort string
	if queryBuilder.SortOrder > 0 {
		sortOrder = "ASC"
	} else {
		sortOrder = "DESC"
	}

	if queryBuilder.SortField != "" {
		sort = `SORT data.` + queryBuilder.SortField + ` ` + sortOrder
	}

	collection := queryBuilder.Collection
	if collection == "" {
		collection = r.Collection
	}

	totalRecordsQuery := `
		FOR ` + alias + ` IN ` + collection +
		joinQuery + " " + filterQuery + `
		COLLECT WITH COUNT INTO length
		RETURN length
	`

	var query string
	if queryBuilder.Rows != 0 {

		query = `
			FOR ` + alias + ` IN ` + collection +
			joinQuery + " " + filterQuery + withQuery +
			" LIMIT " + strconv.Itoa(queryBuilder.First) + ", " + strconv.Itoa(queryBuilder.Rows) + sort + `
			RETURN ` + resultQuery

	} else {
		query = withQuery + `
		FOR ` + alias + ` IN ` + collection +
			joinQuery + " " + filterQuery + sort + `
		RETURN ` + resultQuery
	}

	return totalRecordsQuery, query, filterArgs
}

func (r *ArangoBaseRepository) BuildFilter(s interface{}, filters []ArangoFilterQueryBuilder, joinCollection string, prefixes ...string) []ArangoFilterQueryBuilder {
	v := reflect.Indirect(reflect.ValueOf(s))

	var collectionPrefix string
	if joinCollection == "" || joinCollection == r.Collection {
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
			// Todo : Filter by contents of array[0] in arango's array[*]
		}

	} else {
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanInterface() {
				tags := strings.Split(v.Type().Field(i).Tag.Get("json"), ",")
				value := v.Field(i).Interface()
				if v.Field(i).Kind() == reflect.Struct {
					var tag string
					if collection := joinCollection + v.Type().Field(i).Tag.Get("collection"); collection != "" || tags[0] == r.Collection {
						tag = ""
					} else {
						tag = tags[0]
					}

					filters = r.BuildFilter(value, filters, joinCollection+v.Type().Field(i).Tag.Get("collection"), tag)
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
	}

	return filters
}

func (r *ArangoBaseRepository) First(c context.Context, request ArangoInterface) error {
	queryBuilder := ArangoQueryBuilder{
		Filters: r.BuildFilter(request, []ArangoFilterQueryBuilder{}, ""),
	}

	return r.RawFirst(c, queryBuilder, request)
}

func (r *ArangoBaseRepository) RawFirst(c context.Context, queryBuilder ArangoQueryBuilder, request ArangoInterface) error {

	ctx := driver.WithQueryCount(c)

	_, query, condition := r.buildQuery(queryBuilder)

	data, err := r.ArangoDB.DB().Query(ctx, query, condition)
	if err != nil {
		return err
	}

	defer data.Close()

	if data.Count() > 0 {
		data.ReadDocument(c, &request)
		return nil
	}

	return constant.ErrorNotFound
}

func (r *ArangoBaseRepository) All(c context.Context, request interface{}, baseFilter PaginationFilters) ([]interface{}, int64, error) {
	queryBuilder := ArangoQueryBuilder{
		Filters:   r.BuildFilter(request, []ArangoFilterQueryBuilder{}, ""),
		First:     baseFilter.First,
		Rows:      baseFilter.Rows,
		SortField: baseFilter.SortField,
		SortOrder: baseFilter.SortOrder,
	}

	return r.RawAll(c, queryBuilder)
}

func (r *ArangoBaseRepository) RawAll(c context.Context, queryBuilder ArangoQueryBuilder) ([]interface{}, int64, error) {
	var response []interface{}
	ctx := driver.WithQueryCount(c)

	totalRecordsQuery, query, condition := r.buildQuery(queryBuilder)

	data, err := r.ArangoDB.DB().Query(ctx, query, condition)
	if err != nil {
		return response, 0, err
	}

	defer data.Close()

	if data.Count() > 0 {
		for data.HasMore() {
			var request interface{}
			data.ReadDocument(c, &request)
			response = append(response, request)
		}
	}

	countData, err := r.ArangoDB.DB().Query(ctx, totalRecordsQuery, condition)
	if err != nil {
		return response, 0, err
	}

	defer countData.Close()

	var totalRecords int64
	if countData.Count() > 0 {

		_, err = countData.ReadDocument(ctx, &totalRecords)
		if err != nil {
			return response, 0, err
		}
	}

	return response, totalRecords, nil
}

func (r *ArangoBaseRepository) RawQuery(c context.Context, queryBuilder ArangoQueryBuilder) (string, string, map[string]interface{}) {
	return r.buildQuery(queryBuilder)
}

func (r *ArangoBaseRepository) Create(c context.Context, request ArangoInterface) error {
	collection, err := r.ArangoDB.DB().Collection(c, r.Collection)
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

func (r *ArangoBaseRepository) Update(c context.Context, request ArangoInterface) error {
	collection, err := r.ArangoDB.DB().Collection(c, r.Collection)
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

func (r *ArangoBaseRepository) Delete(c context.Context, request ArangoInterface) error {
	collection, err := r.ArangoDB.DB().Collection(c, r.Collection)
	if err != nil {
		return err
	}

	_, err = collection.RemoveDocument(c, request.GetKey())
	if err != nil {
		return err
	}
	return nil
}

func (r *ArangoBaseRepository) DB() driver.Database {
	return r.ArangoDB.DB()
}
