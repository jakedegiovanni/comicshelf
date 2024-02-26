package comicshelf

type Set[T comparable] map[T]struct{}

func (s Set[T]) Has(t T) bool {
	_, ok := s[t]
	return ok
}

func (s Set[T]) Put(t T) {
	s[t] = struct{}{}
}

func (s Set[T]) Delete(t T) {
	delete(s, t)
}
