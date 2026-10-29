package account

type Handler interface {
}

type handler struct {
	service Service
}

func NewUserHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}
