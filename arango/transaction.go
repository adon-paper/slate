package arango

import (
	"context"

	"github.com/arangodb/go-driver"
)

type Transaction struct {
	ID          driver.TransactionID
	Context     context.Context
	Collections driver.TransactionCollections
}
