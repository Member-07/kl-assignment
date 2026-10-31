package response

const (
	SUCCESS string = "success"
)

type StandardResponse[T any] struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Data    T      `json:"data,omitempty"`
	Pagination
}

type Pagination struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	TotalRows  int64 `json:"total_rows,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

func MappingStandardResponse[T any](message string, data T) StandardResponse[T] {
	return StandardResponse[T]{
		Message: message,
		Data:    data,
	}
}

func MappingStandardResponseWithPagination[T any](message string, data T, pagination Pagination) StandardResponse[T] {
	res := MappingStandardResponse(message, data)
	res.Pagination = pagination
	return res
}
