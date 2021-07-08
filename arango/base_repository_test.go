package arango

import (
	"fmt"
	"testing"
)

func TestArango(t *testing.T) {
	db, _ := InitArangoDBTest()
	repo := NewArangoBaseRepository(db, "paper_chain_payment_requests")

	var q string
	// var maps map[string]interface{}

	q, _ = repo.WithMany(
		SubQuery("has_payment_request").
			Traversal("paper_chain_payment_requests._id", "ANY").
			Returns("invoice_id:invoices._id", "number:invoices.document_number"),
		"invoices",
	).Returns(
		"buyer_id:paper_chain_payment_requests.buyer_id",
		"pr_id:paper_chain_payment_requests._id",
		"pr_key:paper_chain_payment_requests._key",
		"full_pr:paper_chain_payment_requests",
		"invoices:invoices",
	).
		ToQuery()

	fmt.Println(q)
}
