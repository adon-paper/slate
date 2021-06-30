package arango

// type ArangoBaseMockRepository struct {
// 	ArangoBaseRepository
// }

// func (r *ArangoBaseMockRepository) first(c context.Context, request ArangoInterface) error {
// 	return r.ArangoBaseRepository.first(c, request)
// }
// func (r *ArangoBaseMockRepository) rawFirst(c context.Context, queryBuilder ArangoQueryBuilder, request ArangoInterface) error {
// 	return r.ArangoBaseRepository.rawFirst(c, queryBuilder, request)
// }
// func (r *ArangoBaseMockRepository) all(c context.Context, request interface{}, baseFilter PaginationFilters) ([]interface{}, int64, error) {
// 	return r.ArangoBaseRepository.all(c, request, baseFilter)
// }
// func (r *ArangoBaseMockRepository) buildFilter(s interface{}, filters []ArangoFilterQueryBuilder, joinCollection string, prefixes ...string) []ArangoFilterQueryBuilder {
// 	var prefix string
// 	if len(prefixes) > 0 && prefixes[0] != "" {
// 		prefix = prefixes[0] + "."
// 	}

// 	return r.ArangoBaseRepository.buildFilter(s, filters, joinCollection, prefix)
// }
