package maps

import (
	_ "iter"
	"sync"
	"testing"
)

func BenchmarkMaps(b *testing.B) {
	// builtin map
	{
		m := make(map[int]int)
		b.Run("BuiltInMap_Set", func(b *testing.B) {

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m[i] = i
			}

		})
		b.Run("BuiltInMap_Get", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = m[i]
			}
		})
	}

	// sync.Map
	{
		m := &sync.Map{}
		b.Run("BuiltInMap_Set", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Store(i, i)
			}
		})
		b.Run("BuiltInMap_Get", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = m.Load(i)
			}
		})
	}

	//orderedmap
	{
		m := NewOrderedMap[int, int](0)
		b.Run("OrderedMap_Set", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Set(i, i)
			}
		})
		b.Run("OrderedMap_Get", func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = m.Get(i)
			}
		})
	}

	// // github.com/rsc/omap
	// {
	// 	m := &omap.Map[int, int]{}
	// 	b.Run("Omap_Set", func(b *testing.B) {
	// 		b.ResetTimer()
	// 		for i := 0; i < b.N; i++ {
	// 			m.Set(i, i)
	// 		}
	// 	})
	// 	b.Run("Omap_Get", func(b *testing.B) {
	// 		b.ResetTimer()
	// 		for i := 0; i < b.N; i++ {
	// 			_, _ = m.Get(i)
	// 		}
	// 	})
	// }
}
