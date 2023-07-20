package pipe

type Step[T Accumulator] func(T)

type Steps[T Accumulator] []Step[T]

type Accumulator interface {
	Complete() bool
	Error() error
}

func Run[T Accumulator](accumulator T, steps ...Step[T]) {

	for _, step := range steps {

		step(accumulator)

		if accumulator.Error() != nil {
			return
		}

		if accumulator.Complete() {
			return
		}
	}
}
