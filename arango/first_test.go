package arango

import (
	"testing"
)

func TestFirst(t *testing.T) {
	// db, _ := InitArangoDBTest()
	// repo := NewArangoBaseRepository(db, "paper_chain_payment_requests")

	// f := FinancingFacility{}
	// f.Limit = 123
	// f.BNPL = &BNPL{
	// 	Setting: Setting{
	// 		LateFee: 321,
	// 	},
	// }
	// repo.First(context.Background(), &f)
}

type FinancingFacility struct {
	DocumentModel
	Limit         float64        `json:"limit"`
	LimitDue      float64        `json:"limit_due"`
	CompanyUUID   string         `json:"company_uuid"`
	CompanyEmail  string         `json:"company_email"`
	CompanyName   string         `json:"company_name"`
	CompanyPhone  string         `json:"company_phone"`
	BNPL          *BNPL          `json:"bnpl,omitempty"`
	APF           *APF           `json:"apf,omitempty"`
	PaperTrade    *PaperTrade    `json:"paper_trade,omitempty"`
	GrowthAccount *GrowthAccount `json:"growth_account,omitempty"`
}

type BillingTime struct {
	BillingName string `json:"billing_name"`
	StartDate   int    `json:"start_date"`
	EndDate     int    `json:"end_date"`
	DueDate     int    `json:"due_date"`
}

type Setting struct {
	LateFee      float64 `json:"late_fee"`
	InterestRate float64 `json:"interest_rate"`
}

type BNPL struct {
	BaseFinancing
	Setting     Setting       `json:"setting"`
	BillingTime []BillingTime `json:"billing_time"`
}

type APF struct {
	BaseFinancing
	Setting struct {
		LateFee           float64 `json:"late_fee"`
		LatePeriod        int64   `json:"late_period"`
		LoanDueDate       int64   `json:"loan_due_date"`
		InterestRate      float64 `json:"interest_rate"`
		GracePeriodStatus string  `json:"grace_period_status"`
	} `json:"setting"`
}

type DiscountFeeSetting struct {
	Status          int64   `json:"status"`
	MaxDueDate      int64   `json:"max_due_date"`
	MinDueDate      int64   `json:"min_due_date"`
	MaxDiscountRate float64 `json:"max_discount_rate"`
	MinDiscountRate float64 `json:"min_discount_rate"`
	Type            string  `json:"type"` // "flat", "prorate"
}

type PaperTrade struct {
	BaseFinancing
	Setting struct {
		VATDisbursementSetting struct {
			Status int64   `json:"status"`
			PPN    float64 `json:"ppn"`
			PPH23  float64 `json:"pph_23"`
		} `json:"vat_disbursement_setting"`
		VATInvoiceSetting int64 `json:"vat_invoice_setting"`
	} `json:"setting"`
	DiscountFeeSetting []DiscountFeeSetting `json:"discount_fee_setting"`
}

type GrowthAccount struct {
	BaseFinancing
	Setting struct {
		AutoRepayment    bool    `json:"auto_repayment"`
		RepaymentPercent float64 `json:"repayment_percent"`
	} `json:"setting"`
}

type BaseFinancing struct {
	Category string `json:"category"`
	Status   string `json:"status"` // "active", "inactive"
}
