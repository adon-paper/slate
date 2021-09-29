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
	//	WithMany(
	//	SubQuery("has_payment_request").
	//		Traversal("paper_chain_payment_requests._id", "ANY", true).
	//		Returns(
	//			"amount_due:has_payment_request.document_info.amount_due",
	//			"due_date:has_payment_request.document_info.due_date",
	//			"grand_total:has_payment_request.document_info.totals.grandTotalUnformatted",
	//			"invoice_id:has_payment_request.document_info.uuid",
	//			"invoice_number:has_payment_request.document_number",
	//			"status:has_payment_request.document_info.status",
	//		),
	//	"invoices",
	//).WithOne(
	//	SubQuery("companies").
	//		WhereColumn("uuid", "==", "paper_chain_payment_requests.supplier_id"),
	//	"supplier",
	//).
	Join(SubQuery("company").WhereColumn("company._id","==","paper_chain_payment_requests.company_id")).where("count(100)","<=","100").
	ToQuery()

	fmt.Println(q)

	var maps map[string]interface{}

	q, maps = repo.Where("name", "==", "chad").Where("name", "==", "rekt").ToQuery()

	fmt.Println(q)
	fmt.Println(maps)
}
