package slugerrors

type ErrorType string

const (
	ErrorTypeBadRequest ErrorType = "bad_request"
	ErrorTypeNotFound   ErrorType = "not_found"
	ErrorTypeInternal   ErrorType = "internal"
)

type SlugError interface {
	error
	Slug() string
	ErrorType() ErrorType
}

type slugError struct {
	slug    string
	errType ErrorType
	message string
}

func (e *slugError) Error() string {
	return e.message
}

func (e *slugError) Slug() string {
	return e.slug
}

func (e *slugError) ErrorType() ErrorType {
	return e.errType
}

func NewError(slug string, errType ErrorType, message string) SlugError {
	return &slugError{
		slug:    slug,
		errType: errType,
		message: message,
	}
}
