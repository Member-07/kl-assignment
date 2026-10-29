package model

type Banner struct {
	BannerId    string `gorm:"column:banner_id;primaryKey;type:varchar(50)"`
	UserId      string `gorm:"column:user_id;type:varchar(50)"`
	Title       string `gorm:"column:title;type:varchar(255)"`
	Description string `gorm:"column:description;type:text"`
	Image       string `gorm:"column:image;type:varchar(255)"`
	DummyCol11  string `gorm:"column:dummy_col_11;type:varchar(255)"`
}

func (Banner) TableName() string {
	return "banners"
}
