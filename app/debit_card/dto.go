package debitcard

type CardDetailResponse struct {
	Issuer string `json:"issuer"`
	Number string `json:"number"`
}

type CardDesignResponse struct {
	Color       string `json:"color"`
	BorderColor string `json:"border_color"`
}

type DebitCardResponse struct {
	CardId string             `json:"card_id"`
	Name   string             `json:"name"`
	Status string             `json:"status"`
	Detail CardDetailResponse `json:"detail,omitempty"`
	Design CardDesignResponse `json:"design,omitempty"`
}

type UpdateCardStatusRequest struct {
	UserId string `json:"user_id" validate:"required"`
	Status string `json:"status" validate:"required,oneof=Active InActive"`
	CardId string `json:"card_id" validate:"required"`
}
