package chanx

// Batch reads from a channel and calls fn with a slice of batchSize.
func Batch[T any](ch <-chan T, batchSize int, fn func([]T)) {
	for batchSize <= 1 { // sanity check,
		for v := range ch {
			fn([]T{v})
		}

		return
	}

	// batchSize > 1

	var batch = make([]T, 0, batchSize)
	for {
		select {
		case v, ok := <-ch:
			if !ok { // closed
				fn(batch)
				return
			}

			batch = append(batch, v)
			if len(batch) == batchSize { // full
				fn(batch)
				batch = make([]T, 0, batchSize) // reset
			}
		default:
			if len(batch) > 0 { // partial
				fn(batch)
				batch = make([]T, 0, batchSize) // reset
			} else { // empty
				// wait for more
				v, ok := <-ch
				if !ok {
					return
				}

				batch = append(batch, v)
			}
		}
	}

}
