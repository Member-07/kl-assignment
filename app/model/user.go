package model

type User struct {
	UserID    string `gorm:"column:user_id;primaryKey;type:varchar(50)"`
	Name      string `gorm:"column:name;type:varchar(100)"`
	DummyCol1 string `gorm:"column:dummy_col_1;type:varchar(255)"`
}

func (User) TableName() string {
	return "users"
}

type UserGreeting struct {
	UserID    string `gorm:"column:user_id;primaryKey;type:varchar(50)"`
	Greeting  string `gorm:"column:greeting;type:text"`
	DummyCol2 string `gorm:"column:dummy_col_2;type:varchar(255)"`
}

func (UserGreeting) TableName() string {
	return "user_greetings"
}
