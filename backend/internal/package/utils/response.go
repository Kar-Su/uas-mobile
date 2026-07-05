package utils

type Response[T any] struct {
	Success bool            `json:"success"`
	Message string          `json:"message" example:"Operation successful"`
	Path    string          `json:"path,omitempty" example:"/api/resource"`
	Error   any             `json:"error,omitempty"`
	Data    T               `json:"data,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}

type responseOptions struct {
	path string
	meta *PaginationMeta
}

type ResponseOption func(*responseOptions)

func WithPath(path string) ResponseOption {
	return func(o *responseOptions) {
		o.path = path
	}
}

func WithMeta(meta *PaginationMeta) ResponseOption {
	return func(o *responseOptions) {
		o.meta = meta
	}
}

func BuildResponse[T any](message string, data T, opts ...ResponseOption) Response[T] {
	o := &responseOptions{}
	for _, opt := range opts {
		opt(o)
	}
	return Response[T]{
		Success: true,
		Message: message,
		Path:    o.path,
		Data:    data,
		Meta:    o.meta,
	}
}

func BuildResponseFailed(message string, err any, opts ...ResponseOption) Response[any] {
	o := &responseOptions{}
	for _, opt := range opts {
		opt(o)
	}
	return Response[any]{
		Success: false,
		Message: message,
		Path:    o.path,
		Error:   err,
	}
}
