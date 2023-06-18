package maps

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalJSON(t *testing.T) {
	t.Run("int key", func(t *testing.T) {
		m := NewOrderedMap[int, any](10)
		m.Set(1, "bar")
		m.Set(7, "baz")
		m.Set(2, 28)
		m.Set(3, 100)
		m.Set(4, "baz")
		m.Set(5, "28")
		m.Set(6, "100")
		m.Set(8, "baz")
		m.Set(8, "baz")
		m.Set(9, "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque auctor augue accumsan mi maximus, quis viverra massa pretium. Phasellus imperdiet sapien a interdum sollicitudin. Duis at commodo lectus, a lacinia sem.")

		b, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, `{"1":"bar","7":"baz","2":28,"3":100,"4":"baz","5":"28","6":"100","8":"baz","9":"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque auctor augue accumsan mi maximus, quis viverra massa pretium. Phasellus imperdiet sapien a interdum sollicitudin. Duis at commodo lectus, a lacinia sem."}`, string(b))
	})

	t.Run("string key", func(t *testing.T) {
		m := NewOrderedMap[string, any](2)
		m.Set("test", "bar")
		m.Set("abc", true)

		b, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, `{"test":"bar","abc":true}`, string(b))
	})

	t.Run("typed string key", func(t *testing.T) {
		type myString string
		m := NewOrderedMap[myString, any](2)
		m.Set("test", "bar")
		m.Set("abc", true)

		b, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, `{"test":"bar","abc":true}`, string(b))
	})

	t.Run("typed int key", func(t *testing.T) {
		type myInt uint32
		m := NewOrderedMap[myInt, any](5)
		m.Set(1, "bar")
		m.Set(7, "baz")
		m.Set(2, 28)
		m.Set(3, 100)
		m.Set(4, "baz")

		b, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, `{"1":"bar","7":"baz","2":28,"3":100,"4":"baz"}`, string(b))
	})

	t.Run("typed int* key", func(t *testing.T) {
		type myInt uint32
		m := NewOrderedMap[myInt, any](5)
		m.Set(1, int8(1))
		m.Set(2, uint8(2))
		m.Set(3, int16(3))
		m.Set(4, int16(4))
		m.Set(5, int32(5))
		m.Set(6, uint32(6))
		m.Set(7, int64(7))
		m.Set(8, int64(8))
		b, err := json.Marshal(m)
		assert.NoError(t, err)
		assert.Equal(t, `{"1":1,"2":2,"3":3,"4":4,"5":5,"6":6,"7":7,"8":8}`, string(b))
	})

	t.Run("empty map", func(t *testing.T) {
		om := NewOrderedMap[string, any](0)

		b, err := json.Marshal(om)
		assert.NoError(t, err)
		assert.Equal(t, `{}`, string(b))
	})
}

func TestUnmarshallJSON(t *testing.T) {
	t.Run("int key", func(t *testing.T) {
		data := `{"1":"bar","7":"baz","2":28,"3":100,"4":"baz","5":"28","6":"100","8":"baz"}`

		om := NewOrderedMap[int, any](10)
		require.NoError(t, json.Unmarshal([]byte(data), &om))

		assertOrderedPairsEqual(t, om,
			[]int{1, 7, 2, 3, 4, 5, 6, 8},
			[]any{"bar", "baz", float64(28), float64(100), "baz", "28", "100", "baz"})
	})

	t.Run("string key", func(t *testing.T) {
		data := `{"test":"bar","abc":true}`

		om := NewOrderedMap[string, any](10)
		require.NoError(t, json.Unmarshal([]byte(data), &om))

		assertOrderedPairsEqual(t, om,
			[]string{"test", "abc"},
			[]any{"bar", true})
	})

	t.Run("typed string key", func(t *testing.T) {
		data := `{"test":"bar","abc":true}`

		type myString string
		om := NewOrderedMap[myString, any](10)
		require.NoError(t, json.Unmarshal([]byte(data), &om))

		assertOrderedPairsEqual(t, om,
			[]myString{"test", "abc"},
			[]any{"bar", true})
	})

	t.Run("typed int key", func(t *testing.T) {
		data := `{"1":"bar","7":"baz","2":28,"3":100,"4":"baz","5":"28","6":"100","8":"baz"}`

		type myInt uint32
		om := NewOrderedMap[myInt, any](10)
		require.NoError(t, json.Unmarshal([]byte(data), &om))

		assertOrderedPairsEqual(t, om,
			[]myInt{1, 7, 2, 3, 4, 5, 6, 8},
			[]any{"bar", "baz", float64(28), float64(100), "baz", "28", "100", "baz"})
	})

	t.Run("when fed with an input that's not an object", func(t *testing.T) {
		for _, data := range []string{"true", `["foo"]`, "42", `"foo"`} {
			om := NewOrderedMap[int, any](10)
			require.Error(t, json.Unmarshal([]byte(data), &om))
		}
	})

	t.Run("empty map", func(t *testing.T) {
		data := `{}`

		om := NewOrderedMap[int, any](0)
		require.NoError(t, json.Unmarshal([]byte(data), &om))

		assertLenEqual(t, om, 0)
	})
}

func assertOrderedPairsEqual[K comparable, V any](
	t *testing.T, orderedMap *OrderedMap[K, V], expectedKeys []K, expectedValues []V,
) {
	t.Helper()

	assertOrderedPairsEqualFromNewest(t, orderedMap, expectedKeys, expectedValues)
	assertOrderedPairsEqualFromOldest(t, orderedMap, expectedKeys, expectedValues)
}

func assertOrderedPairsEqualFromNewest[K comparable, V any](
	t *testing.T, orderedMap *OrderedMap[K, V], expectedKeys []K, expectedValues []V,
) {
	t.Helper()

	if assert.Equal(t, len(expectedKeys), len(expectedValues)) && assert.Equal(t, len(expectedKeys), orderedMap.Len()) {
		i := orderedMap.Len() - 1
		for entry := orderedMap.Newest(); entry != nil; entry = entry.Prev() {
			assert.Equal(t, expectedKeys[i], entry.Key, "from newest index=%d on key", i)
			assert.Equal(t, expectedValues[i], entry.Value, "from newest index=%d on value", i)
			i--
		}
	}
}

func assertOrderedPairsEqualFromOldest[K comparable, V any](
	t *testing.T, orderedMap *OrderedMap[K, V], expectedKeys []K, expectedValues []V,
) {
	t.Helper()

	if assert.Equal(t, len(expectedKeys), len(expectedValues)) && assert.Equal(t, len(expectedKeys), orderedMap.Len()) {
		i := 0
		for entry := orderedMap.Oldest(); entry != nil; entry = entry.Next() {
			assert.Equal(t, expectedKeys[i], entry.Key, "from oldest index=%d on key", i)
			assert.Equal(t, expectedValues[i], entry.Value, "from oldest index=%d on value", i)
			i++
		}
	}
}

func assertLenEqual[K comparable, V any](t *testing.T, orderedMap *OrderedMap[K, V], expectedLen int) {
	t.Helper()

	assert.Equal(t, expectedLen, orderedMap.Len())

	// also check the list length, for good measure
	assert.Equal(t, expectedLen, orderedMap.list.Len())
}
