package dir3

type Y1 struct {
	*baseCache
}

func NewY1() *Y1 {
	return &Y1{
		baseCache: newBaseCache(),
	}
}
