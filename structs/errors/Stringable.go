package errors

type StringableError interface {
	Stringify() string
}
