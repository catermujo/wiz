package wiz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNil(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	var x int
	is.False(IsNil(x))

	var k struct{}
	is.False(IsNil(k))

	var s *string
	is.True(IsNil(s))

	var i *int
	is.True(IsNil(i))

	var b *bool
	is.True(IsNil(b))

	var ifaceWithNilValue any = (*string)(nil) //nolint:staticcheck
	is.True(IsNil(ifaceWithNilValue))
	is.False(ifaceWithNilValue == nil) //nolint:staticcheck
}

func TestToPtr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	result1 := ToPtr([]int{1, 2})

	is.Equal(*result1, []int{1, 2})
}

func TestNil(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	nilFloat64 := Nil[float64]()
	var expNilFloat64 *float64

	nilString := Nil[string]()
	var expNilString *string

	is.Equal(expNilFloat64, nilFloat64)
	is.Nil(nilFloat64)
	is.NotEqual(nil, nilFloat64)

	is.Equal(expNilString, nilString)
	is.Nil(nilString)
	is.NotEqual(nil, nilString)

	is.NotEqual(nilString, nilFloat64)
}

func TestEmptyableToPtr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	is.Nil(EmptyableToPtr(0))
	is.Nil(EmptyableToPtr(""))
	is.Nil(EmptyableToPtr[[]int](nil))
	is.Nil(EmptyableToPtr[map[int]int](nil))
	is.Nil(EmptyableToPtr[error](nil))

	is.Equal(*EmptyableToPtr(42), 42)
	is.Equal(*EmptyableToPtr("nonempty"), "nonempty")
	is.Equal(*EmptyableToPtr([]int{}), []int{})
	is.Equal(*EmptyableToPtr([]int{1, 2}), []int{1, 2})
	is.Equal(*EmptyableToPtr(map[int]int{}), map[int]int{})
	is.Equal(*EmptyableToPtr(assert.AnError), assert.AnError)
}

func TestFromPtr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	str1 := "foo"
	ptr := &str1

	is.Equal("foo", FromPtr(ptr))
	is.Equal("", FromPtr[string](nil))
	is.Equal(0, FromPtr[int](nil))
	is.Nil(FromPtr[*string](nil))
	is.EqualValues(ptr, FromPtr(&ptr))
}

func TestFromPtrOr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	const fallbackStr = "fallback"
	str := "foo"
	ptrStr := &str

	const fallbackInt = -1
	i := 9
	ptrInt := &i

	is.Equal(str, FromPtrOr(ptrStr, fallbackStr))
	is.Equal(fallbackStr, FromPtrOr(nil, fallbackStr))
	is.Equal(i, FromPtrOr(ptrInt, fallbackInt))
	is.Equal(fallbackInt, FromPtrOr(nil, fallbackInt))
}

func TestToSlicePtr(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	str1 := "foo"
	str2 := "bar"
	result1 := ToSlicePtr([]string{str1, str2})

	is.Equal(result1, []*string{&str1, &str2})
}

func TestFromSlicePtr(t *testing.T) {
	is := assert.New(t)

	str1 := "foo"
	str2 := "bar"
	result1 := FromSlicePtr([]*string{&str1, &str2, nil})

	is.Equal(result1, []string{str1, str2, ""})
}

func TestFromSlicePtrOr(t *testing.T) {
	is := assert.New(t)

	str1 := "foo"
	str2 := "bar"
	result1 := FromSlicePtrOr([]*string{&str1, &str2, nil}, "fallback")

	is.Equal(result1, []string{str1, str2, "fallback"})
}

func TestToAnySlice(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	in1 := []int{0, 1, 2, 3}
	in2 := []int{}
	out1 := ToAnySlice(in1)
	out2 := ToAnySlice(in2)

	is.Equal([]any{0, 1, 2, 3}, out1)
	is.Equal([]any{}, out2)
}

func TestFromAnySlice(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	is.NotPanics(func() {
		out1, ok1 := FromAnySlice[string]([]any{"foobar", 42})
		out2, ok2 := FromAnySlice[string]([]any{"foobar", "42"})

		is.Equal([]string{}, out1)
		is.False(ok1)
		is.Equal([]string{"foobar", "42"}, out2)
		is.True(ok2)
	})
}

func TestEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	//nolint:unused
	type test struct{}

	is.Empty(Empty[string]())
	is.Empty(Empty[int64]())
	is.Empty(Empty[test]())
	is.Empty(Empty[chan string]())
}

func TestIsEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	//nolint:unused
	type test struct {
		foobar string
	}

	is.True(IsEmpty(""))
	is.False(IsEmpty("foo"))
	is.True(IsEmpty[int64](0))
	is.False(IsEmpty[int64](42))
	is.True(IsEmpty(test{foobar: ""}))
	is.False(IsEmpty(test{foobar: "foo"}))
}

func TestIsNotEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	//nolint:unused
	type test struct {
		foobar string
	}

	is.False(IsNotEmpty(""))
	is.True(IsNotEmpty("foo"))
	is.False(IsNotEmpty[int64](0))
	is.True(IsNotEmpty[int64](42))
	is.False(IsNotEmpty(test{foobar: ""}))
	is.True(IsNotEmpty(test{foobar: "foo"}))
}

func TestCoalesce(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	newStr := func(v string) *string { return &v }
	var nilStr *string
	str1 := newStr("str1")
	str2 := newStr("str2")

	type structType struct {
		field1 int
		field2 float64
	}
	var zeroStruct structType
	struct1 := structType{1, 1.0}
	struct2 := structType{2, 2.0}

	result1, ok1 := Coalesce[int]()
	result2, ok2 := Coalesce(3)
	result3, ok3 := Coalesce(nil, nilStr)
	result4, ok4 := Coalesce(nilStr, str1)
	result5, ok5 := Coalesce(nilStr, str1, str2)
	result6, ok6 := Coalesce(str1, str2, nilStr)
	result7, ok7 := Coalesce(0, 1, 2, 3)
	result8, ok8 := Coalesce(zeroStruct)
	result9, ok9 := Coalesce(zeroStruct, struct1)
	result10, ok10 := Coalesce(zeroStruct, struct1, struct2)

	is.Equal(0, result1)
	is.False(ok1)

	is.Equal(3, result2)
	is.True(ok2)

	is.Nil(result3)
	is.False(ok3)

	is.Equal(str1, result4)
	is.True(ok4)

	is.Equal(str1, result5)
	is.True(ok5)

	is.Equal(str1, result6)
	is.True(ok6)

	is.Equal(result7, 1)
	is.True(ok7)

	is.Equal(result8, zeroStruct)
	is.False(ok8)

	is.Equal(result9, struct1)
	is.True(ok9)

	is.Equal(result10, struct1)
	is.True(ok10)
}

func TestCoalesceOrEmpty(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	newStr := func(v string) *string { return &v }
	var nilStr *string
	str1 := newStr("str1")
	str2 := newStr("str2")

	type structType struct {
		field1 int
		field2 float64
	}
	var zeroStruct structType
	struct1 := structType{1, 1.0}
	struct2 := structType{2, 2.0}

	result1 := CoalesceOrEmpty[int]()
	result2 := CoalesceOrEmpty(3)
	result3 := CoalesceOrEmpty(nil, nilStr)
	result4 := CoalesceOrEmpty(nilStr, str1)
	result5 := CoalesceOrEmpty(nilStr, str1, str2)
	result6 := CoalesceOrEmpty(str1, str2, nilStr)
	result7 := CoalesceOrEmpty(0, 1, 2, 3)
	result8 := CoalesceOrEmpty(zeroStruct)
	result9 := CoalesceOrEmpty(zeroStruct, struct1)
	result10 := CoalesceOrEmpty(zeroStruct, struct1, struct2)

	is.Equal(0, result1)
	is.Equal(3, result2)
	is.Nil(result3)
	is.Equal(str1, result4)
	is.Equal(str1, result5)
	is.Equal(str1, result6)
	is.Equal(result7, 1)
	is.Equal(result8, zeroStruct)
	is.Equal(result9, struct1)
	is.Equal(result10, struct1)
}
