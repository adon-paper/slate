package arango

import (
	"fmt"
	"testing"
)

func TestArango(t *testing.T) {
	db, _ := InitArangoDBTest()
	repo := NewArangoBaseRepository(db, "companies")

	var q string
	var maps map[string]interface{}

	// q, _ = repo.Where("_id", "==", "company_id").
	// 	WhereOr("company_email", "==", "test@test.com").
	// 	WhereColumn("companies.company_name", "==", "Karya Anak Rumahan").
	// 	Join("companies", "_id", "digital_payment_requests", "company_id").
	// 	With("companies", "_id", "digital_payment_requests", "company_id", "dpr").
	// 	Raw()

	// q, _ = repo.Where("company_email", "LIKE", "arnold.widjaja@paper").
	// 	With("companies", "_id", "digital_payment_requests", "company_id", "dpr").
	// 	Raw()

	// q, _ = repo.Where("company_email", "LIKE", "arnold.widjaja@paper").
	// 	WithEdge("companies", "_id", "has_service", "ff", "any").
	// 	Raw()

	// q, _ = repo.Where("company_email", "LIKE", "arnold.widjaja@paper").
	// 	JoinEdge("companies", "_id", "has_service", "c", "any").
	// 	Raw()

	q, maps = repo.Where("company_email", "LIKE", "arnold.widjaja@paper").
		JoinEdge("companies", "_id", "has_service", "ff", "any").
		WithOne(
			SubQuery("digital_payment_requests").Where("_key", "==", "1").
				WithOne(
					SubQuery("digital_payment_transaction").
						WhereColumn("_key", "==", "dpr.payment_request"),
					"dpt",
				),
				"dpr",
		).
		Raw()

	fmt.Println(q)
	fmt.Println(maps)
}
