package tuple

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTuple(t *testing.T) {
	t2 := &Tuple2[int, string]{1, "a"}
	assert.Equal(t, 1, t2.V1)
	assert.Equal(t, "a", t2.V2)

	t3 := &Tuple3[int, string, bool]{1, "a", true}
	assert.Equal(t, 1, t3.V1)
	assert.Equal(t, "a", t3.V2)
	assert.Equal(t, true, t3.V3)

	t4 := &Tuple4[int, string, bool, float64]{1, "a", true, 1.0}
	assert.Equal(t, 1, t4.V1)
	assert.Equal(t, "a", t4.V2)
	assert.Equal(t, true, t4.V3)
	assert.Equal(t, 1.0, t4.V4)

	t5 := &Tuple5[int, string, bool, float64, int]{1, "a", true, 1.0, 1}
	assert.Equal(t, 1, t5.V1)
	assert.Equal(t, "a", t5.V2)
	assert.Equal(t, true, t5.V3)
	assert.Equal(t, 1.0, t5.V4)
	assert.Equal(t, 1, t5.V5)

	t6 := &Tuple6[int, string, bool, float64, int, time.Time]{1, "a", true, 1.0, 1, time.Time{}}
	assert.Equal(t, 1, t6.V1)
	assert.Equal(t, "a", t6.V2)
	assert.Equal(t, true, t6.V3)
	assert.Equal(t, 1.0, t6.V4)
	assert.Equal(t, 1, t6.V5)
	assert.Equal(t, time.Time{}, t6.V6)

	t7 := &Tuple7[int, string, bool, float64, int, time.Time, int]{1, "a", true, 1.0, 1, time.Time{}, 1}
	assert.Equal(t, 1, t7.V1)
	assert.Equal(t, "a", t7.V2)
	assert.Equal(t, true, t7.V3)
	assert.Equal(t, 1.0, t7.V4)
	assert.Equal(t, 1, t7.V5)
	assert.Equal(t, time.Time{}, t7.V6)
	assert.Equal(t, 1, t7.V7)

	t8 := &Tuple8[int, string, bool, float64, int, time.Time, int, string]{1, "a", true, 1.0, 1, time.Time{}, 1, "a"}
	assert.Equal(t, 1, t8.V1)
	assert.Equal(t, "a", t8.V2)
	assert.Equal(t, true, t8.V3)
	assert.Equal(t, 1.0, t8.V4)
	assert.Equal(t, 1, t8.V5)
	assert.Equal(t, time.Time{}, t8.V6)
	assert.Equal(t, 1, t8.V7)
	assert.Equal(t, "a", t8.V8)

	t9 := &Tuple9[int, string, bool, float64, int, time.Time, int, string, bool]{1, "a", true, 1.0, 1, time.Time{}, 1, "a", true}
	assert.Equal(t, 1, t9.V1)
	assert.Equal(t, "a", t9.V2)
	assert.Equal(t, true, t9.V3)
	assert.Equal(t, 1.0, t9.V4)
	assert.Equal(t, 1, t9.V5)
	assert.Equal(t, time.Time{}, t9.V6)
	assert.Equal(t, 1, t9.V7)
	assert.Equal(t, "a", t9.V8)
	assert.Equal(t, true, t9.V9)

	t10 := &Tuple10[int, string, bool, float64, int, time.Time, int, string, bool, float64]{
		1, "a", true, 1.0, 1, time.Time{}, 1, "a", true, 1.0}
	assert.Equal(t, 1, t10.V1)
	assert.Equal(t, "a", t10.V2)
	assert.Equal(t, true, t10.V3)
	assert.Equal(t, 1.0, t10.V4)
	assert.Equal(t, 1, t10.V5)
	assert.Equal(t, time.Time{}, t10.V6)
	assert.Equal(t, 1, t10.V7)
	assert.Equal(t, "a", t10.V8)
	assert.Equal(t, true, t10.V9)
	assert.Equal(t, 1.0, t10.V10)

	t11 := &Tuple11[int, string, bool, float64, int, time.Time, int, string, bool, float64, int]{
		1, "a", true, 1.0, 1, time.Time{}, 1, "a", true, 1.0, 1}
	assert.Equal(t, 1, t11.V1)
	assert.Equal(t, "a", t11.V2)
	assert.Equal(t, true, t11.V3)
	assert.Equal(t, 1.0, t11.V4)
	assert.Equal(t, 1, t11.V5)
	assert.Equal(t, time.Time{}, t11.V6)
	assert.Equal(t, 1, t11.V7)
	assert.Equal(t, "a", t11.V8)
	assert.Equal(t, true, t11.V9)
	assert.Equal(t, 1.0, t11.V10)

	t12 := &Tuple12[int, string, bool, float64, int, time.Time, int, string, bool, float64, int, string]{
		1, "a", true, 1.0, 1, time.Time{}, 1, "a", true, 1.0, 1, "a"}
	assert.Equal(t, 1, t12.V1)
	assert.Equal(t, "a", t12.V2)
	assert.Equal(t, true, t12.V3)
	assert.Equal(t, 1.0, t12.V4)
	assert.Equal(t, 1, t12.V5)
	assert.Equal(t, time.Time{}, t12.V6)
	assert.Equal(t, 1, t12.V7)
	assert.Equal(t, "a", t12.V8)
	assert.Equal(t, true, t12.V9)
	assert.Equal(t, 1.0, t12.V10)
	assert.Equal(t, 1, t12.V11)
	assert.Equal(t, "a", t12.V12)
}
