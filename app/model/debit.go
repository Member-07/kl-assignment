package model

type DebitCard struct {
	CardID    string `gorm:"column:card_id;primaryKey;type:varchar(50)"`
	UserID    string `gorm:"column:user_id;type:varchar(50)"`
	Name      string `gorm:"column:name;type:varchar(100)"`
	DummyCol7 string `gorm:"column:dummy_col_7;type:varchar(255)"`
}

func (DebitCard) TableName() string {
	return "debit_cards"
}

type DebitCardStatus struct {
	CardID    string `gorm:"column:card_id;primaryKey;type:varchar(50)"`
	UserID    string `gorm:"column:user_id;type:varchar(50)"`
	Status    string `gorm:"column:status;type:varchar(20)"`
	DummyCol8 string `gorm:"column:dummy_col_8;type:varchar(255)"`
}

func (DebitCardStatus) TableName() string {
	return "debit_card_status"
}

type DebitCardDesign struct {
	CardID      string `gorm:"column:card_id;primaryKey;type:varchar(50)"`
	UserID      string `gorm:"column:user_id;type:varchar(50)"`
	Color       string `gorm:"column:color;type:varchar(10)"`
	BorderColor string `gorm:"column:border_color;type:varchar(10)"`
	DummyCol9   string `gorm:"column:dummy_col_9;type:varchar(255)"`
}

func (DebitCardDesign) TableName() string {
	return "debit_card_design"
}

type DebitCardDetail struct {
	CardID     string `gorm:"column:card_id;primaryKey;type:varchar(50)"`
	UserID     string `gorm:"column:user_id;type:varchar(50)"`
	Issuer     string `gorm:"column:issuer;type:varchar(100)"`
	Number     string `gorm:"column:number;type:varchar(25)"`
	DummyCol10 string `gorm:"column:dummy_col_10;type:varchar(255)"`
}

func (DebitCardDetail) TableName() string {
	return "debit_card_details"
}
