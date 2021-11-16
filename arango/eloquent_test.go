package arango

import (
	"fmt"
	"testing"
)

func TestEloquent(t *testing.T) {
	db, _ := InitArangoDBTest()
	repo := NewArangoBaseRepository(db, "paper_chain_payment_requests")

	var q string
	// // var maps map[string]interface{}

	q, _ = repo.
		WithMany(SubQueryWithAlias("payment_reconciliation_transactions", "prt1").
			WhereColumn("company_id", "==", "digital_payment_finance_accounts.company_id").
			Where("status", "IN", []string{"POSTED", "settlement"}).
			Where("created_at", ">=", "2021-09-13T16:44:44.266164424+07:00").
			Returns("total:prt1.amount.grand_total"),
			"debit").
		WithMany(SubQueryWithAlias("payment_reconciliation_transactions", "prt2").
			WhereColumn("company_id", "==", "digital_payment_finance_accounts.company_id").
			Where("status", "IN", []string{"DRAFT", "PENDING", "FRAUD"}).
			Where("created_at", ">=", "2021-09-13T16:44:44.266164424+07:00").
			Returns("total:prt2.amount.grand_total"),
			"onhold").
		WithMany(SubQuery("digpay_disbursement_transactions").
			WhereColumn("digpay_disbursement_transactions.company_id", "==", "digital_payment_finance_accounts.company_id").
			Where("digpay_disbursement_transactions.disbursement_request_no", "LIKE", "DISB%").
			Returns("total:digpay_disbursement_transactions.disbursement_amount"),
			"disbursed").
		WhereColumn("digital_payment_finance_accounts.debit_balance", "!=", "sum(debit[*].total)").
		WhereOrColumn("digital_payment_finance_accounts.credit_balance", "!=", "sum(disbursed[*].total)").
		WhereOrColumn("digital_payment_finance_accounts.on_hold_amount", "!=", "sum(onhold[*].total)").
		Where("digital_payment_finance_accounts.on_hold_amount", "!=", "0").
		Where("digital_payment_finance_accounts.on_hold_amount", "!=", nil).
		Returns("debit :sum(debit[*].total)", "onhold :sum(onhold[*].total)", "disbursed :sum(disbursed[*].total)", " digital_payment_finance_accounts").ToQuery()

	fmt.Println(q)

	var maps map[string]interface{}

	q, maps = repo.Where("name", "==", "chad").Where("name", "==", "rekt").ToQuery()

	fmt.Println(q)
	fmt.Println(maps)
}
