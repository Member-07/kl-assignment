package user

type UserResponse struct {
	UserId      string `json:"user_id"`
	Name        string `json:"name"`
	GreetingMsg string `json:"greeting_msg"`
}
