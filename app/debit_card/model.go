package debitcard

type DebitCard struct {
	CardID      string `gorm:"column:card_id"`
	UserID      string `gorm:"column:user_id"`
	Name        string `gorm:"column:name"`
	Issuer      string `gorm:"column:issuer"`
	Number      string `gorm:"column:number"`
	Color       string `gorm:"column:color"`
	BorderColor string `gorm:"column:border_color"`
	Status      string `gorm:"column:status"`
}

func (DebitCard) TableName() string {
	return "debit_cards"
}

type DebitCardDetail struct {
	CardID string `gorm:"column:card_id;primaryKey"`
	UserID string `gorm:"column:user_id"`
	Issuer string `gorm:"column:issuer"`
	Number string `gorm:"column:number"`
}

func (DebitCardDetail) TableName() string {
	return "debit_card_details"
}

type DebitCardDesign struct {
	CardID      string `gorm:"column:card_id;primaryKey"`
	UserID      string `gorm:"column:user_id"`
	Color       string `gorm:"column:color"`
	BorderColor string `gorm:"column:border_color"`
}

func (DebitCardDesign) TableName() string {
	return "debit_card_design"
}

type DebitCardStatus struct {
	CardId string `gorm:"column:card_id;primaryKey"`
	UserId string `gorm:"column:user_id"`
	Status string `gorm:"column:status"`
}

func (DebitCardStatus) TableName() string {
	return "debit_card_status"
}
