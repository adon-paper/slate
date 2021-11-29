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

	// q, _ = repo.Where("ref_id", "==", "REF1637343607lsJrm").WithMany(
	// 	SubQuery("digital_payment_transactions").
	// 		WhereColumn("payment_request_id", "==", "digital_payment_requests._key"),
	// 	"has_transactions",
	// ).WithOne(
	// 	SubQuery("digital_payment_transactions").
	// 		WhereColumn("payment_request_id", "==", " digital_payment_requests._key").
	// 		Where("digital_payment_transactions.status", "in", `["PENDING","WAITING"]`).
	// 		Sort("created_at", "DESC"),
	// 	"has_open_transaction",
	// ).WithOne(
	// 	SubQuery("digital_payment_transactions").
	// 		WhereColumn("payment_request_id", "==", " digital_payment_requests._key"),
	// 	"has_paid_transaction",
	// ).WithOne(
	// 	SubQuery("digital_payment_documents").
	// 		WhereColumn("digital_payment_documents.payment_request_id", "==", " digital_payment_requests._key"),
	// 	"has_digital_payment_document",
	// ).ToQuery()

	q, _ = repo.Where("_id", "==", "paymentrequestId").
		WithOne(
			SubQuery("is_payment_out_request").
				Traversal("paper_chain_payment_requests._id", "ANY").
				WithMany(
					SubQueryWithAlias("has_payment_out_request", "pi").
						Traversal("is_payment_out_request._id", "ANY", true),
					"purchase_invoices",
				),
			"payment_out_request_details",
		).WithOne(
		SubQuery("companies").
			WhereColumn("uuid", "==", "paper_chain_payment_requests.buyer_id"),
		"buyer",
	).ToQuery()

	fmt.Println(q)

	// var maps map[string]interface{}

	// q, maps = repo.Where("name", "==", "chad").Where("name", "==", "rekt").ToQuery()

	// fmt.Println(q)
	// fmt.Println(maps)
}
