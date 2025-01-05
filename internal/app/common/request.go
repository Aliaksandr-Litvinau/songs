package common

import "context"

// RequestReader interface for reading request parameters
type RequestReader interface {
	// PathParam returns the parameter value from the URL path
	PathParam(name string) (string, error)

	// QueryParam returns the value of the query parameter
	QueryParam(name string) string

	// DefaultQueryParam returns the value of the query parameter or the default value
	DefaultQueryParam(name, defaultValue string) string

	// DecodeBody decodes the request body into a structure
	DecodeBody(interface{}) error

	// Context returns the request context
	Context() context.Context
}
