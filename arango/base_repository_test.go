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
	// 	ToQuery()

	// q, _ = repo.Where("company_email", "LIKE", "arnold.widjaja@paper").
	// 	With("companies", "_id", "digital_payment_requests", "company_id", "dpr").
	// 	ToQuery()

	// q, _ = repo.Where("company_email", "LIKE", "arnold.widjaja@paper").
	// 	WithEdge("companies", "_id", "has_service", "ff", "any").
	// 	ToQuery()

	// q, _ = repo.Where("company_email", "LIKE", "arnold.widjaja@paper").
	// 	JoinEdge("companies", "_id", "has_service", "c", "any").
	// 	ToQuery()

	// q, maps = repo.Where("company_email", "LIKE", "arnold.widjaja@paper").
	// 	JoinEdge("companies", "_id", "has_service", "ff", "any").
	// 	WithOne(
	// 		SubQuery("digital_payment_requests").Where("_key", "==", "1").
	// 			WithOne(
	// 				SubQuery("digital_payment_transaction").
	// 					WhereColumn("_key", "==", "dpr.payment_request"),
	// 				"dpt",
	// 			),
	// 			"dpr",
	// 	).
	// 	ToQuery()

	// q, _ = repo.Where("company_email", "LIKE", "arnold.widjaja@paper").
	// 	WithOne(
	// 		SubQuery("has_service").Traversal("companies._id", "ANY").Where("_key", "==", "1").
	// 			WithOne(
	// 				SubQuery("digital_payment_transaction").
	// 					WhereColumn("_key", "==", "dpr.payment_request"),
	// 				"dpt",
	// 			),
	// 		"dpr",
	// 	).
	// 	ToQuery()

	// fmt.Println(q)

	// FOR companies in companies
	// FILTER companies.company_email LIKE @companies_company_email
	// LET dpr = (
	// 	FOR has_service in ANY companies._id has_service
	// 	FILTER has_service._key == @has_service__key
	// 	RETURN MERGE(has_service)
	// )
	// RETURN MERGE({dpr :dpr}, companies)

	// q, _ = repo.
	// 	Where("company_email", "LIKE", "arnold.widjaja@paper").
	// 	WithOne(
	// 		SubQuery("has_service").
	// 			Where("_key", "==", "1").
	// 			Traversal("companies._id", "ANY"),
	// 		"dpr",
	// 	).
	// 	ToQuery()

	// fmt.Println(q)

	// FOR companies in companies
	// LET ff = (
	// 	FOR has_service in ANY companies._id has_service  FILTER has_service._key == @has_service__key
	// 	LET institution = (
	// 		FOR has_financed in has_financed
	// 		FILTER ff._id == has_service._id
	// 		RETURN MERGE(has_financed)
	// 	)
	// 	RETURN MERGE({institution :institution}, has_service)
	// )
	// RETURN MERGE({ff :ff}, companies)
	q, maps = repo.WithMany(
		SubQuery("has_service").
			Traversal("companies._id", "ANY").
			Where("_key", "==", "1").
			WithMany(
				SubQuery("has_financed").WhereColumn("ff._id", "==", "has_service._id"),
				"institution",
			),
		"ff",
	).ToQuery()
	fmt.Println(q)
	fmt.Println(maps)

	// FOR companies in ANY companies._id companies  FILTER companies.financing_type == @companies_financing_type   RETURN MERGE(companies)
	// q, _ = repo.Traversal("companies._id", "ANY").Where("financing_type", "==", "paper_trade_retailer").ToQuery()
	// fmt.Println(q)

	// FOR companies in companies  FILTER companies.financing_type == @companies_financing_type   RETURN MERGE(companies)
	// q, _ = repo.Where("financing_type", "==", "paper_trade_retailer").ToQuery()
	// fmt.Println(q)

}
