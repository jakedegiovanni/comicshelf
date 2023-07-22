package main

var (
	followed = make(map[string]struct{})
)

func follow(series string) {
	followed[series] = struct{}{}
}

func unfollow(series string) {
	delete(followed, series)
}

func following(series string) bool {
	_, ok := followed[series]
	return ok
}
