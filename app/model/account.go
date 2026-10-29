package model

import (
	"time"
)

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

type Account struct {
	AccountId     string `gorm:"column:account_id;primaryKey;type:varchar(50)"`
	UserId        string `gorm:"column:user_id;type:varchar(50)"`
	Type          string `gorm:"column:type;type:varchar(50)"`
	Currency      string `gorm:"column:currency;type:varchar(10)"`
	AccountNumber string `gorm:"column:account_number;type:varchar(20)"`
	Issuer        string `gorm:"column:issuer;type:varchar(100)"`
	DummyCol3     string `gorm:"column:dummy_col_3;type:varchar(255)"`
}

func (Account) TableName() string {
	return "accounts"
}

type AccountDetail struct {
	AccountId     string `gorm:"column:account_id;primaryKey;type:varchar(50)"`
	UserId        string `gorm:"column:user_id;type:varchar(50)"`
	Color         string `gorm:"column:color;type:varchar(10)"`
	IsMainAccount bool   `gorm:"column:is_main_account;type:tinyint(1)"`
	Progress      int    `gorm:"column:progress;type:int"`
	DummyCol5     string `gorm:"column:dummy_col_5;type:varchar(255)"`
}

func (AccountDetail) TableName() string {
	return "account_details"
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

type AccountBalance struct {
	AccountId string  `gorm:"column:account_id;primaryKey;type:varchar(50)"`
	UserId    string  `gorm:"column:user_id;type:varchar(50)"`
	Amount    float64 `gorm:"column:amount;type:decimal(15,2)"`
	DummyCol4 string  `gorm:"column:dummy_col_4;type:varchar(255)"`
}

func (AccountBalance) TableName() string {
	return "account_balances"
}
