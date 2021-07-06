package arango

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/arangodb/go-driver"
)

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
