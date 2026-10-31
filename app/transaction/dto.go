package transaction

import "kl-assignment/pkg/response"

type TransactionResponse struct {
	TransactionId string `json:"transaction_id"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	IsBank        bool   `json:"is_bank"`
}

type TransactionPaginationResponse struct {
	Data []TransactionResponse `json:"data"`
	response.Pagination
}
