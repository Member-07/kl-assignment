package account

import "time"

type AccountWithDetail struct {
	AccountId     string  `gorm:"column:account_id"`
	UserId        string  `gorm:"column:user_id"`
	Type          string  `gorm:"column:type"`
	Currency      string  `gorm:"column:currency"`
	AccountNumber string  `gorm:"column:account_number"`
	Issuer        string  `gorm:"column:issuer"`
	Amount        float64 `gorm:"column:amount"`
	Color         string  `gorm:"column:color"`
	IsMainAccount bool    `gorm:"column:is_main_account"`
	Progress      int     `gorm:"column:progress"`
}

type AccountFlag struct {
	FlagId    int       `gorm:"column:flag_id;primaryKey;autoIncrement;type:int"`
	AccountId string    `gorm:"column:account_id;type:varchar(50)"`
	UserId    string    `gorm:"column:user_id;type:varchar(50)"`
	FlagType  string    `gorm:"column:flag_type;type:varchar(50)"`
	FlagValue string    `gorm:"column:flag_value;type:varchar(30)"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp"`
}

func (AccountFlag) TableName() string {
	return "account_flags"
}
