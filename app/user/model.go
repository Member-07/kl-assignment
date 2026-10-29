package user

type UserWithGreeting struct {
	UserId   string `gorm:"column:user_id"`
	Name     string `gorm:"column:name"`
	Greeting string `gorm:"column:greeting"`
}
