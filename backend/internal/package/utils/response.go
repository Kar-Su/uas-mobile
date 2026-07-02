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

// ResponseOption is a functional option for BuildResponse and BuildResponseFailed.
type ResponseOption func(*responseOptions)

// WithPath sets the path field on the response.
func WithPath(path string) ResponseOption {
	return func(o *responseOptions) {
		o.path = path
	}
}

// WithMeta attaches pagination metadata to the response.
// Use this to produce a paginated response where meta sits alongside data.
func WithMeta(meta *PaginationMeta) ResponseOption {
	return func(o *responseOptions) {
		o.meta = meta
	}
}

// BuildResponse builds a success response.
// Pass WithMeta(NewPaginationMeta(...)) to produce a paginated response.
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

// BuildResponseFailed builds a failure response.
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
