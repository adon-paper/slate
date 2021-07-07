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

	q, _ = repo.Where("buyer_id", "==", "request.BuyerId").
		Where("source", "==", "digpayout").
		Where("_key", "==", "request.PaymentRequestKey").
		WithMany(
			SubQuery("has_payment_request").Traversal("paper_chain_payment_requests._id", "ANY"),
			"invoices",
		).ToQuery()

	fmt.Println(q)
}
