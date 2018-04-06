package errors

type Error interface {
	Error() string
	Message() string
	Code() int
	Status() int
	ErrorStack() []error
	SetMessage(string) Error
	SetStatus(int) Error
	SetCode(int) Error
}

func New(message string) Error {
	return &errorImplementation{
		message: message,
	}
}

func From(err error) Error {
	result := &errorImplementation{
		message: err.Error(),
		errorStack: make([]error, 0),
	}
	if e, ok := err.(Error); ok {
		result.code = e.Code()
		result.status = e.Status()
		result.errorStack = append(result.errorStack, e.ErrorStack()...)
	}
	result.errorStack = append(result.errorStack, err)
	return result
}

type errorImplementation struct {
	message string
	code int
	status int
	errorStack []error
}

func (e *errorImplementation) Error() string {
	return e.message
}

func (e *errorImplementation) Message() string {
	return e.message
}

func (e *errorImplementation) Code() int {
	return e.code
}

func (e *errorImplementation) Status() int {
	return e.status
}

func (e *errorImplementation) ErrorStack() []error {
	return e.errorStack
}

func (e *errorImplementation) SetMessage(m string) Error {
	e.message = m
	return e
}

func (e *errorImplementation) SetStatus(s int) Error {
	e.status = s
	return e
}

func (e *errorImplementation) SetCode(c int) Error {
	e.code = c
	return e
}