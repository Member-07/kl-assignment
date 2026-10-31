package banner

type BannerResponse struct {
	BannerId    string  `json:"banner_id"`
	Title       *string `json:"title"`
	Description string  `json:"description"`
	Image       *string `json:"image"`
}
