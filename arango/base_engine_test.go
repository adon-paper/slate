package arango

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/arangodb/go-driver"
)

type transaction struct {
	DocumentModel
	Description string `json:"description"`
}

func TestInitArangoDB(t *testing.T) {
	success := false

	fmt.Println("running")

	c := context.Background()

	db, err := InitArangoDB(
		"http://localhost:8529",
		"paper-import",
		"service",
		"service123",
		c,
	)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	tcollection := driver.TransactionCollections{
		Read: []string{
			"documents",
			"edge",
		},
		Write: []string{},
		Exclusive: []string{
			"transaction",
		},
	}

	tcollection.Read = append(tcollection.Read, "transaction")

	tID, err := db.DB().BeginTransaction(c, tcollection, &driver.BeginTransactionOptions{LockTimeout: 30 * time.Second})
	if err != nil {
		fmt.Println("begin transaction err", err)
	}

	fmt.Println(tID)

	repo := NewArangoBaseRepository(db, "transaction")

	var transaction transaction
	if success {
		transaction.Description = "success"
	} else {
		transaction.Description = "fail"
	}

	tctx := driver.WithTransactionID(c, tID)

	err = repo.Create(tctx, &transaction)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("transaction object id", transaction.Id)

	if success {
		err = db.DB().CommitTransaction(c, tID, nil)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err = db.DB().AbortTransaction(c, tID, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}
