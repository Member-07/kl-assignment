package transaction

type Transaction struct {
	TransactionId string  `gorm:"column:transaction_id;primaryKey"`
	UserId        *string `gorm:"column:user_id"`
	Name          *string `gorm:"column:name"`
	Image         *string `gorm:"column:image"`
	IsBank        *bool   `gorm:"column:isBank"`
}

func (Transaction) TableName() string {
	return "transactions"
}
