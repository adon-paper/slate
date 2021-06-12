package arango

import (
	"context"
	"fmt"
	"sync"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type arangoDB struct {
	client driver.Client
	db     driver.Database
}

type ArangoDB interface {
	Client() driver.Client
	DB() driver.Database
	InitCollections(cols []string) error
}

func InitArangoDB(host, database, username, password string, context context.Context) (ArangoDB, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{host},
		ConnLimit: 32,
	})
	if err != nil {
		return nil, fmt.Errorf("error when creating new connection to arango : %s", err.Error())
	}
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(username, password),
	})
	if err != nil {
		return nil, fmt.Errorf("error when creating new arango client : %s", err.Error())
	}

	db, err := connect(database, c, context)
	if err != nil {
		return nil, fmt.Errorf("error when opening database connection : %s", err.Error())
	}
	return &arangoDB{
		client: c,
		db:     db,
	}, nil
}

func connect(database string, client driver.Client, c context.Context) (driver.Database, error) {
	ctx := driver.WithQueryCount(c)
	db, err := client.Database(ctx, database)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (a *arangoDB) Client() driver.Client {
	return a.client
}

func (a *arangoDB) DB() driver.Database {
	return a.db
}

func (a *arangoDB) InitCollections(cols []string) error {
	wg := sync.WaitGroup{}
	for _, col := range cols {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()
			//options := &driver.CreateCollectionOptions{ Type: }
			colExist, err := a.db.CollectionExists(ctx, col)
			if err == nil && !colExist {
				a.db.CreateCollection(ctx, col, nil)
			}
		}()
	}
	wg.Wait()
	return nil
}
