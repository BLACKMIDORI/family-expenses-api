package services

type HttpError interface {
	StatusCode() int
	Error() string
}

type InternalServerError struct {
	Message string
}
type BadRequest struct {
	Message string
}
type NotFound struct {
	Message string
}
type ServiceUnavailable struct {
	Message string
}

func (_ InternalServerError) StatusCode() int {
	return 500
}

func (badRequest InternalServerError) Error() string {
	return badRequest.Message
}

func (_ BadRequest) StatusCode() int {
	return 400
}

func (badRequest BadRequest) Error() string {
	return badRequest.Message
}
func (_ NotFound) StatusCode() int {
	return 404
}

func (notFound NotFound) Error() string {
	return notFound.Message
}

func (_ ServiceUnavailable) StatusCode() int {
	return 503
}

func (badRequest ServiceUnavailable) Error() string {
	return badRequest.Message
}
