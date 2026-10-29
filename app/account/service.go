package account

type Service interface {
}

type service struct {
	repo Repository
}

func NewUserService(repo Repository) Service {
	return &service{repo: repo}
}
