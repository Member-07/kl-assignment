package model

type Transaction struct {
	TransactionID string `gorm:"column:transaction_id;primaryKey;type:varchar(50)"`
	UserID        string `gorm:"column:user_id;type:varchar(50)"`
	Name          string `gorm:"column:name;type:varchar(100)"`
	Image         string `gorm:"column:image;type:varchar(255)"`
	IsBank        bool   `gorm:"column:isBank;type:tinyint(1)"`
	DummyCol6     string `gorm:"column:dummy_col_6;type:varchar(255)"`
}

func (Transaction) TableName() string {
	return "transactions"
}
