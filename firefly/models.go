package firefly

import "time"

type FireflyWebhookResponse struct {
	UUID        string  `json:"uuid"`
	UserID      int     `json:"user_id"`
	UserGroupID int     `json:"user_group_id"`
	Trigger     string  `json:"trigger"`
	Response    string  `json:"response"`
	URL         string  `json:"url"`
	Version     string  `json:"version"`
	Content     Content `json:"content"`
}

type Content struct {
	ID           int           `json:"id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	User         int           `json:"user"`
	GroupTitle   interface{}   `json:"group_title"`
	Transactions []Transaction `json:"transactions"`
	Links        []Link        `json:"links"`
}

type Transaction struct {
	User                         int         `json:"user"`
	TransactionJournalID         string      `json:"transaction_journal_id"`
	Type                         string      `json:"type"`
	Date                         time.Time   `json:"date"`
	Order                        int         `json:"order"`
	CurrencyID                   string      `json:"currency_id"`
	CurrencyCode                 string      `json:"currency_code"`
	CurrencySymbol               string      `json:"currency_symbol"`
	CurrencyDecimalPlaces        int         `json:"currency_decimal_places"`
	ForeignCurrencyID            string      `json:"foreign_currency_id"`
	ForeignCurrencyCode          interface{} `json:"foreign_currency_code"`
	ForeignCurrencySymbol        interface{} `json:"foreign_currency_symbol"`
	ForeignCurrencyDecimalPlaces interface{} `json:"foreign_currency_decimal_places"`
	Amount                       string      `json:"amount"`
	ForeignAmount                interface{} `json:"foreign_amount"`
	Description                  string      `json:"description"`
	SourceID                     string      `json:"source_id"`
	SourceName                   string      `json:"source_name"`
	SourceIban                   string      `json:"source_iban"`
	SourceType                   string      `json:"source_type"`
	DestinationID                string      `json:"destination_id"`
	DestinationName              string      `json:"destination_name"`
	DestinationIban              interface{} `json:"destination_iban"`
	DestinationType              string      `json:"destination_type"`
	BudgetID                     string      `json:"budget_id"`
	BudgetName                   interface{} `json:"budget_name"`
	CategoryID                   string      `json:"category_id"`
	CategoryName                 interface{} `json:"category_name"`
	BillID                       string      `json:"bill_id"`
	BillName                     interface{} `json:"bill_name"`
	Reconciled                   bool        `json:"reconciled"`
	Notes                        interface{} `json:"notes"`
	Tags                         []string    `json:"tags"`
	InternalReference            interface{} `json:"internal_reference"`
	ExternalID                   interface{} `json:"external_id"`
	OriginalSource               string      `json:"original_source"`
	RecurrenceID                 interface{} `json:"recurrence_id"`
	BunqPaymentID                interface{} `json:"bunq_payment_id"`
	ImportHashV2                 string      `json:"import_hash_v2"`
	SepaCc                       interface{} `json:"sepa_cc"`
	SepaCtOp                     interface{} `json:"sepa_ct_op"`
	SepaCtID                     interface{} `json:"sepa_ct_id"`
	SepaDb                       interface{} `json:"sepa_db"`
	SepaCountry                  interface{} `json:"sepa_country"`
	SepaEp                       interface{} `json:"sepa_ep"`
	SepaCi                       interface{} `json:"sepa_ci"`
	SepaBatchID                  interface{} `json:"sepa_batch_id"`
	InterestDate                 interface{} `json:"interest_date"`
	BookDate                     interface{} `json:"book_date"`
	ProcessDate                  interface{} `json:"process_date"`
	Due_date                     interface{} `json:"due_date"`
	PaymentDate                  interface{} `json:"payment_date"`
	InvoiceDate                  interface{} `json:"invoice_date"`
	Longitude                    interface{} `json:"longitude"`
	Latitude                     interface{} `json:"latitude"`
	ZoomLevel                    interface{} `json:"zoom_level"`
}

type Link struct {
	Rel string `json:"rel"`
	URI string `json:"uri"`
}

type FireflyTxRequest struct {
	Transactions []Transaction `json:"transactions"` // Il tag JSON deve rimanere "transactions"
}
