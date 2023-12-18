package chanx

import "context"

// Batch reads from a channel and calls fn with a slice of batchSize.
func Batch[T any](ctx context.Context, ch <-chan T, batchSize int, fn func([]T)) {
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
		case <-ctx.Done():
			if len(batch) > 0 {
				fn(batch)
			}
			return
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
				select {
				case <-ctx.Done():
					if len(batch) > 0 {
						fn(batch)
					}
					return
				case v, ok := <-ch:
					if !ok {
						return
					}

					batch = append(batch, v)
				}

			}
		}
	}

}

// FlatBatch reads from a channel of slices, flats values, and calls fn with a slice of batchSize.
func FlatBatch[T any](ctx context.Context, ch <-chan []T, batchSize int, fn func([]T)) {
	for batchSize <= 1 { // sanity check,
		for v := range ch {
			fn(v)
		}

		return
	}

	// batchSize > 1
	var batch = make([]T, 0, batchSize)
	for {
		select {
		case <-ctx.Done():
			if len(batch) > 0 {
				fn(batch)
			}
			return
		case v, ok := <-ch:
			if !ok { // closed
				fn(batch)
				return
			}

			batch = append(batch, v...)
			if len(batch) >= batchSize { // full
				fn(batch)
				batch = make([]T, 0, batchSize) // reset
			}
		default:
			if len(batch) > 0 { // partial
				fn(batch)
				batch = make([]T, 0, batchSize) // reset
			} else { // empty
				// wait for more
				select {
				case <-ctx.Done():
					if len(batch) > 0 {
						fn(batch)
					}
					return
				case v, ok := <-ch:
					if !ok {
						return
					}

					batch = append(batch, v...)
					if len(batch) >= batchSize { // full
						fn(batch)
						batch = make([]T, 0, batchSize) // reset
					}
				}
			}
		}
	}

}
