package transaction

type TransactionResponse struct {
	TransactionId string `json:"transaction_id"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	IsBank        bool   `json:"is_bank"`
}

type TransactionPaginationResponse struct {
	Data       []TransactionResponse `json:"data"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalRows  int64                 `json:"total_rows"`
	TotalPages int                   `json:"total_pages"`
}
