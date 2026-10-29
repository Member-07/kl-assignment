package response

const (
	SUCCESS string = "success"
)

type StandardResponse[T any] struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func MappingStandardResponse[T any](message string, data T) StandardResponse[T] {
	return StandardResponse[T]{
		Message: message,
		Data:    data,
	}
}
