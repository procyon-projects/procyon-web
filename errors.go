package web

type RouterError struct {
	message string
}

func NewRouterError(message string) RouterError {
	return RouterError{message}
}

func (err RouterError) Error() string {
	return err.message
}

type NoHandlerFoundError struct {
	message string
}

func NewNoHandlerFoundError(message string) NoHandlerFoundError {
	return NoHandlerFoundError{message}
}

func (err NoHandlerFoundError) Error() string {
	return err.message
}

type NoHandlerParameterResolver struct {
	message string
}

func NewNoHandlerParameterResolver(message string) NoHandlerParameterResolver {
	return NoHandlerParameterResolver{message}
}

func (err NoHandlerParameterResolver) Error() string {
	return err.message
}
