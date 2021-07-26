package arango

func (r *ArangoBaseRepository) Where(param ...interface{}) *ArangoQuery {
	return NewQuery(r.Collection, r.ArangoDB).Where(param...)
}

func (r *ArangoBaseRepository) WhereOr(column string, operator string, value interface{}) *ArangoQuery {
	return NewQuery(r.Collection, r.ArangoDB).WhereOr(column, operator, value)
}

func (r *ArangoBaseRepository) WhereColumn(column string, operator string, value string) *ArangoQuery {
	return NewQuery(r.Collection, r.ArangoDB).WhereColumn(column, operator, value)
}

func (r *ArangoBaseRepository) WithOne(repo *ArangoQuery, alias string) *ArangoQuery {
	return NewQuery(r.Collection, r.ArangoDB).WithOne(repo, alias)
}

func (r *ArangoBaseRepository) WithMany(repo *ArangoQuery, alias string) *ArangoQuery {
	return NewQuery(r.Collection, r.ArangoDB).WithMany(repo, alias)
}

func (r *ArangoBaseRepository) Join(query *ArangoQuery) *ArangoQuery {
	return NewQuery(r.Collection, r.ArangoDB).Join(query)

}

func (r *ArangoBaseRepository) Offset(offset int) *ArangoQuery {
	return NewQuery(r.Collection, r.ArangoDB).Offset(offset)
}

func (r *ArangoBaseRepository) Limit(limit int) *ArangoQuery {
	return NewQuery(r.Collection, r.ArangoDB).Limit(limit)
}

func (r *ArangoBaseRepository) Sort(sortField, sortOrder string) *ArangoQuery {
	return NewQuery(r.Collection, r.ArangoDB).Sort(sortField, sortOrder)
}

func (r *ArangoBaseRepository) Traversal(sourceId string, direction traversalDirection) *ArangoQuery {
	return NewQuery(r.Collection, r.ArangoDB).Traversal(sourceId, direction)
}

func (r *ArangoBaseRepository) Returns(returns ...string) *ArangoQuery {
	return NewQuery(r.Collection, r.ArangoDB).Returns(returns...)
}

func (r *ArangoBaseRepository) Get(request interface{}) error {
	return NewQuery(r.Collection, r.ArangoDB).Get(request)
}

func (r *ArangoBaseRepository) Count(request interface{}) error {
	return NewQuery(r.Collection, r.ArangoDB).Count(request)
}

func (r *ArangoBaseRepository) NewQuery() *ArangoQuery {
	return NewQuery(r.Collection, r.ArangoDB)
}
