package account

import "github.com/shopspring/decimal"

type AccountResponse struct {
	AccountId     string                `json:"account_id"`
	AccountNumber string                `json:"account_number"`
	Amount        decimal.Decimal       `json:"amount"`
	Currency      string                `json:"currency"`
	UserId        string                `json:"user_id"`
	Type          string                `json:"type"`
	Issuer        string                `json:"issuer"`
	Color         string                `json:"color"`
	IsMainAccount bool                  `json:"is_main_account"`
	Progress      int                   `json:"progress"`
	Flags         []AccountFlagResponse `json:"flags"`
}

type AccountFlagResponse struct {
	FlagType  string `json:"flag_type"`
	FlagValue string `json:"flag_value"`
}
