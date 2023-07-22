package main

var (
	followed = make(map[string]struct{})
)

func action(series string) {
	if following(series) {
		delete(followed, series)
		return
	}
	followed[series] = struct{}{}
}

func following(series string) bool {
	_, ok := followed[series]
	return ok
}
