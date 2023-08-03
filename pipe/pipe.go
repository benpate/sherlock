package pipe

type Step[T any] func(T) bool

type Steps[T any] []Step[T]

func Run[T any](accumulator T, steps ...Step[T]) bool {
	for _, step := range steps {
		if step(accumulator) {
			return true
		}
	}

	return false
}
