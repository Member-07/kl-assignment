package model

type UserWithGreeting struct {
	UserID      string `gorm:"column:user_id"`
	Name        string `gorm:"column:name"`
	GreetingMsg string `gorm:"column:message"`
}
