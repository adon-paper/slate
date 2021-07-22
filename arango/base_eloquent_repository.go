package arango

func (r *ArangoBaseRepository) Where(param ...interface{}) *ArangoQuery {
	q := NewQuery(r.Collection).Where(param...)
	q.ArangoDB = r.ArangoDB

	return q
}

func (r *ArangoBaseRepository) WhereOr(column string, operator string, value interface{}) *ArangoQuery {
	q := NewQuery(r.Collection).WhereOr(column, operator, value)
	q.ArangoDB = r.ArangoDB

	return q
}

func (r *ArangoBaseRepository) WhereColumn(column string, operator string, value string) *ArangoQuery {
	q := NewQuery(r.Collection).WhereColumn(column, operator, value)
	q.ArangoDB = r.ArangoDB

	return q
}

func (r *ArangoBaseRepository) WithOne(repo *ArangoQuery, alias string) *ArangoQuery {
	q := NewQuery(r.Collection).WithOne(repo, alias)
	q.ArangoDB = r.ArangoDB

	return q
}

func (r *ArangoBaseRepository) WithMany(repo *ArangoQuery, alias string) *ArangoQuery {
	q := NewQuery(r.Collection).WithMany(repo, alias)
	q.ArangoDB = r.ArangoDB

	return q
}

func (r *ArangoBaseRepository) Join(query *ArangoQuery) *ArangoQuery {
	q := NewQuery(r.Collection).Join(query)
	q.ArangoDB = r.ArangoDB

	return q
}

func (r *ArangoBaseRepository) Offset(offset int) *ArangoQuery {
	q := NewQuery(r.Collection).Offset(offset)
	q.ArangoDB = r.ArangoDB

	return q
}

func (r *ArangoBaseRepository) Limit(limit int) *ArangoQuery {
	q := NewQuery(r.Collection).Limit(limit)
	q.ArangoDB = r.ArangoDB

	return q
}

func (r *ArangoBaseRepository) Sort(sortField, sortOrder string) *ArangoQuery {
	q := NewQuery(r.Collection).Sort(sortField, sortOrder)
	q.ArangoDB = r.ArangoDB

	return q
}

func (r *ArangoBaseRepository) Traversal(sourceId string, direction traversalDirection) *ArangoQuery {
	q := NewQuery(r.Collection).Traversal(sourceId, direction)
	q.ArangoDB = r.ArangoDB

	return q
}

func (r *ArangoBaseRepository) Returns(returns ...string) *ArangoQuery {
	q := NewQuery(r.Collection).Returns(returns...)
	q.ArangoDB = r.ArangoDB

	return q
}

func (r *ArangoBaseRepository) Get(request interface{}) error {
	q := NewQuery(r.Collection)
	q.ArangoDB = r.ArangoDB

	return q.Get(request)
}

func (r *ArangoBaseRepository) Count(request interface{}) error {
	q := NewQuery(r.Collection)
	q.ArangoDB = r.ArangoDB

	return q.Count(request)
}
