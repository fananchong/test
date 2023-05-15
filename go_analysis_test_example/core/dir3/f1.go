package dir3

type baseCache struct {
}

func newBaseCache() *baseCache {
	return &baseCache{}
}

func (ca *baseCache) Walk(walkFunc func()) {
	walkFunc()
}
