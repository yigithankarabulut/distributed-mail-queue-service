package customerror

var (
	ErrNotFound = New("not found", false)
	ErrInternal = New("internal server error", true)
	ErrBadReq   = New("bad request", false)
	// ...
)

type CustomError interface {
	Wrap(err error) CustomError
	Unwrap() error
	AddData(any) CustomError
	DestoryData() CustomError
	Error() string
}

type Error struct {
	Err      error
	Message  string
	Data     any `json:"-"`
	Loggable bool
}

func (e *Error) AddData(data any) CustomError {
	e.Data = data
	return e
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) DestoryData() CustomError {
	e.Data = nil
	return e
}

func (e *Error) Wrap(err error) CustomError {
	e.Err = err
	return e
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error() + ", " + e.Message
	}
	return e.Message
}

func New(m string, l bool) CustomError {
	return &Error{
		Message:  m,
		Loggable: l,
	}
}
