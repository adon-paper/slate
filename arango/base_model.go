package arango

type ArangoQueryBuilder struct {
	Filters   []ArangoFilterQueryBuilder
	Joins     []ArangoJoinQueryBuilder
	First     int
	Rows      int
	SortField string
	SortOrder int
}

type ArangoDatalistResponse struct {
	Result       []ArangoInterface `json:"result"`
	TotalRecords int64             `json:"total_records"`
}

type ArangoFilterQueryBuilder struct {
	Key          string
	ArgumentKey  string
	Operator     string
	Value        interface{}
	CustomFilter string
	AndOr        string
}

type ArangoJoinQueryBuilder struct {
	CollectionFrom string
	FromKey        string
	CollectionTo   string
	ToKey          string
	ResultKey      string
}

type PaginationFilters struct {
	SortField string `json:"sortField"`
	SortOrder int    `json:"sortOrder"`
	First     int    `json:"first"`
	Rows      int    `json:"rows"`
}
