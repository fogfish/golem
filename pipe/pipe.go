package pipe

func New[T any](n int) (<-chan T, chan<- T) {
	in := make(chan T, n)
	eg := make(chan T, n)

	mq := newq[T]()

	go func() {
		var (
			x  T
			ok bool
		)

		for {
			select {
			case x, ok = <-in:
				if !ok {
					return
				}
				enq(&x, mq)

			case emit(eg, mq) <- head(mq):
				deq(mq)
			}
		}
	}()

	return eg, in
}
