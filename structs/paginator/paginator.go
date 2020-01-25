package paginator

type Paginator struct {
	Limit  interface{}
	Page   int
	Offset interface{}
}
