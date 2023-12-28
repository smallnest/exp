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

// BatchFilter reads from a channel, filters values, and calls fn with a slice of batchSize.
func BatchFilter[T, V any](ctx context.Context, ch <-chan T, batchSize int, filter func(T) V, fn func([]V)) {
	for batchSize <= 1 { // sanity check,
		for v := range ch {
			fn([]V{filter(v)})
		}

		return
	}

	// batchSize > 1
	var batch = make([]V, 0, batchSize)
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

			batch = append(batch, filter(v))
			if len(batch) == batchSize { // full
				fn(batch)
				batch = make([]V, 0, batchSize) // reset
			}
		default:
			if len(batch) > 0 { // partial
				fn(batch)
				batch = make([]V, 0, batchSize) // reset
			} else { // empty
				// wait for more
				select {
				case <-ctx.Done():
					return
				case v, ok := <-ch:
					if !ok {
						return
					}

					batch = append(batch, filter(v))
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

// FlatBatchFilter reads from a channel of slices, flats values, filters values, and calls fn with a slice of batchSize.
func FlatBatchFilter[T, V any](ctx context.Context, ch <-chan []T, batchSize int, filter func(T) V, fn func([]V)) {
	for batchSize <= 1 { // sanity check,
		for v := range ch {
			var fiv = make([]V, len(v))
			for i, vv := range v {
				fiv[i] = filter(vv)
			}
			fn(fiv)
		}

		return
	}

	// batchSize > 1
	var batch = make([]V, 0, batchSize)
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

			for _, vv := range v {
				batch = append(batch, filter(vv))
			}

			if len(batch) >= batchSize { // full
				fn(batch)
				batch = make([]V, 0, batchSize) // reset
			}
		default:
			if len(batch) > 0 { // partial
				fn(batch)
				batch = make([]V, 0, batchSize) // reset
			} else { // empty
				// wait for more
				select {
				case <-ctx.Done():
					return
				case v, ok := <-ch:
					if !ok {
						return
					}

					for _, vv := range v {
						batch = append(batch, filter(vv))
					}
					if len(batch) >= batchSize { // full
						fn(batch)
						batch = make([]V, 0, batchSize) // reset
					}
				}
			}
		}
	}

}
