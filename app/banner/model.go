package banner

type Banner struct {
	BannerId    string  `gorm:"column:banner_id;primaryKey"`
	UserID      *string `gorm:"column:user_id"`
	Title       *string `gorm:"column:title"`
	Description string  `gorm:"column:description"`
	Image       *string `gorm:"column:image"`
	DummyCol11  *string `gorm:"column:dummy_col_11"`
}

func (Banner) TableName() string {
	return "banners"
}
