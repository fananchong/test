package dir4

var Handlers map[int]func()

func SetHandler(i int, handler func()) {
	Handlers[i] = handler
}

func GetHandler(i int) func() {
	return Handlers[i]
}
