package arango

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/arangodb/go-driver"
)

type DigitalPaymentFinanceAccount struct {
	DocumentModel
	CompanyId             string  `json:"company_id"`
	FinanceAccountId      string  `json:"finance_account_id"`
	Name                  string  `json:"name"`
	OnHoldAmount          float64 `json:"on_hold_amount"`
	CreditBalance         float64 `json:"credit_balance"`
	DebitBalance          float64 `json:"debit_balance"`
	OnHoldAmountDisbursed float64 `json:"on_hold_amount_disbursement"`
}

func TestTransactionConcurrent(t *testing.T) {
	fmt.Println("running test concurrent")
	c := context.Background()

	db, err := InitArangoDB(
		"http://localhost:8529",
		"database",
		"service",
		"password",
		c,
	)

	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	tcollection := driver.TransactionCollections{
		Read:  []string{},
		Write: []string{},
		Exclusive: []string{
			"digital_payment_finance_accounts",
		},
	}

	wg := sync.WaitGroup{}

	now := time.Now()

	repo := NewArangoBaseRepository(db, "digital_payment_finance_accounts")

	for i := 1; i < 51; i++ {
		if i%10 == 0 {
			fmt.Println("waiting for transaction kelar")
			wg.Wait()
		}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
			c := context.Background()
			updateWithTransactionV2(c, repo, tcollection)
			// updateWithTransactionV1(c, repo, tcollection)
			// updateWithoutTransaction(c, db, tcollection)
		}(i)
	}

	wg.Wait()

	fmt.Println(time.Since(now))

	fmt.Println("done")
}

func updateWithoutTransaction(c context.Context, db ArangoDB, tcollection driver.TransactionCollections) {
	repo := NewArangoBaseRepository(db, "digital_payment_finance_accounts")

	var fa DigitalPaymentFinanceAccount
	err := repo.Where("company_id", "slate-transaction-testing-without-transaction").Get(&fa)
	if err != nil {
		fmt.Println("err when get ", err.Error())
	}

	fa.DebitBalance += 10000
	repo.Update(context.Background(), &fa)
}

func updateWithTransactionV2(c context.Context, repo ArangoBaseRepositoryInterface, tcollection driver.TransactionCollections) {

	transaction, err := repo.BeginTransaction(c, tcollection.Read, tcollection.Write, tcollection.Exclusive, nil)
	if err != nil {
		fmt.Println("err when begin transaction ", err.Error())
	}

	var fa DigitalPaymentFinanceAccount
	err = repo.Where("company_id", "slate-transaction-testing-with-transaction").GetWithContext(transaction.Context, &fa)
	if err != nil {
		fmt.Println("err when get ", err.Error())
	}

	fa.DebitBalance += 10000
	repo.Update(transaction.Context, &fa)

	repo.CommitTransaction(c, transaction)

	fmt.Println("transaction done")
}

func updateWithTransactionV1(c context.Context, repo ArangoBaseRepositoryInterface, tcollection driver.TransactionCollections) {
	tID, err := repo.DB().BeginTransaction(c, tcollection, &driver.BeginTransactionOptions{LockTimeout: 50 * time.Second})
	if err != nil {
		fmt.Println("begin transaction err", err)
	}

	tctx := driver.WithTransactionID(c, tID)

	var fa DigitalPaymentFinanceAccount
	err = repo.Where("company_id", "slate-transaction-testing-with-transaction").GetWithContext(tctx, &fa)
	if err != nil {
		fmt.Println("err when get ", err.Error())
	}

	fa.DebitBalance += 10000
	repo.Update(tctx, &fa)

	err = repo.DB().CommitTransaction(c, tID, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("transaction done")
}
