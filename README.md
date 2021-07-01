# slate

### Importing 

go get -u github.com/noldwidjaja/slate


### Usage Example 

#### Main file
```
import "github.com/noldwidjaja/slate/arango"

func main() {    
    ctx := context.Background()

    arangoDB, err := arango.InitArangoDB(os.Getenv(constant.EnvArangoDBHost), os.Getenv(constant.EnvArangoDBDatabase),
        os.Getenv(constant.EnvArangoDBUser), os.Getenv(constant.EnvArangoDBPassword), ctx)

    repo := NewRepository(arangoDB)
}

```

#### Repository File
```
import "github.com/noldwidjaja/slate/arango"

type model struct {
	arango.DocumentModel

    // insert your datas here
    Name string `json:"name"`
}

type RepositoryInterface interface {
	First(c context.Context, request model) (model, error)
	All(c context.Context, request model, sortField string, sortOrder, first, rows int) ([]model, int64, error)
	Create(c context.Context, request *model) error
	Update(c context.Context, request *model) error
    Delete(c context.Context, request *model) error
}

type repository struct {
	arango.ArangoBaseRepository
}

func NewRepository(arangoDB arangodb.ArangoDB) RepositoryInterface {
	repo := repository{
		ArangoBaseRepository: baseArango.ArangoBaseRepository{
			Collection: "digpay_disbursement_transactions",
			ArangoDB:   arangoDB,
		},
	}

	return &repo
}


func (r *repository) First(c context.Context, request model) (model, error) {
	err := r.ArangoBaseRepository.First(c, &request)

	return request, err
}

func (r *repository) All(c context.Context, request model, sortField string, sortOrder, first, rows int) ([]model, int64, error) {
	var response []model

	paginationFilters := baseArango.PaginationFilters{
		SortField: sortField,
		SortOrder: sortOrder,
		First:     first,
		Rows:      rows,
	}

	result, totalRecords, err := r.ArangoBaseRepository.All(c, &request, paginationFilters)
	if err != nil {
		return response, 0, err
	}

	bytedata, _ := json.Marshal(result)
	json.Unmarshal(bytedata, &response)

	return response, totalRecords, err
}

func (r *repository) Create(c context.Context, request *model) error {
	return r.ArangoBaseRepository.Create(c, request)
}

func (r *repository) Update(c context.Context, request *model) error {
	return r.ArangoBaseRepository.Update(c, request)
}

func (r *repository) Delete(c context.Context, request *model) error {
	return r.ArangoBaseRepository.Delete(c, request)
}

```


#### Repository Usage 

```
    var firstModel model 
    // Searches the first model with name test
    firstModel.Name = "test"
    firstModel, err := repo.First(c, firstModel)
```