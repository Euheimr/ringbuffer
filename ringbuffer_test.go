package ringbuffer

import (
	"errors"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		expectedType string
		capacity     int
	}{
		{"negative capacity", -1},
		{"zero capacity", 0},
		{"string", 5},
		{"int", 5},
		{"byte", 5},
		{"float", 5},
		{"bool", 5},
	}

	for _, test := range tests {
		t.Run(test.expectedType, func(t *testing.T) {
			switch test.expectedType {
			case "negative capacity":
				rb, err := New[string](test.capacity)
				if rb == nil || errors.Is(err, errCapacityNegativeOrZero) {
					t.Log("expected error when creating a zero length buffer, all is ok!")
				} else {
					t.Errorf("zero size buffer should be producing an error but incorrectly returns: %s", err)
					t.Fail()
				}
			case "zero size":
				rb, err := New[string](test.capacity)
				if rb == nil || errors.Is(err, errCapacityNegativeOrZero) {
					t.Log("expected error when creating a zero length buffer, all is ok!")
				} else {
					t.Errorf("zero size buffer should be producing an error but incorrectly returns: %s", err)
					t.Fail()
				}
			case "string":
				rb, _ := New[string](test.capacity)
				buffType := reflect.TypeOf(rb.buffer[0]).String()
				if buffType != "string" {
					t.Errorf("incorrect buffer type, expected %s but got %s",
						test.expectedType, buffType)
					t.Fail()
				}
			case "int":
				rb, _ := New[int](test.capacity)
				buffType := reflect.TypeOf(rb.buffer[0]).String()
				if buffType != "int" {
					t.Errorf("incorrect buffer type, expected %s but got %s",
						test.expectedType, buffType)
					t.Fail()
				}
			case "uint":
				rb, _ := New[uint](test.capacity)
				buffType := reflect.TypeOf(rb.buffer[0]).String()
				if buffType != "uint" {
					t.Errorf("incorrect buffer type, expected %s but got %s",
						test.expectedType, buffType)
					t.Fail()
				}
			case "byte":
				rb, _ := New[byte](test.capacity)
				buffType := reflect.TypeOf(rb.buffer[0]).String()
				if buffType != "uint8" {
					t.Errorf("incorrect buffer type, expected %s but got %s",
						test.expectedType, buffType)
					t.Fail()
				}
			case "float":
				rb, _ := New[float32](test.capacity)
				buffType := reflect.TypeOf(rb.buffer[0]).String()
				if buffType != "float32" {
					t.Errorf("incorrect buffer type, expected %s but got %s",
						test.expectedType, buffType)
					t.Fail()
				}
			case "bool":
				rb, _ := New[bool](test.capacity)
				buffType := reflect.TypeOf(rb.buffer[0]).String()
				if buffType != "bool" {
					t.Errorf("incorrect buffer type, expected %s but got %s",
						test.expectedType, buffType)
					t.Fail()
				}
			}
		})
	}
}

func TestNewSize(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
	}{
		{"zero capacity", 0},
		{"buffer too small", 5},
		{"resize buffer", 5},
		{"resize buffer with no values", 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.name {
			case "zero capacity":
				rb, err := New[int](2)
				if err = rb.WriteMany([]int{0, 1}); err != nil {
					t.Logf("failed to write to buffer: %s", err)
					t.Fail()
				}

				rb, err = rb.NewSize(test.capacity)
				if rb == nil || errors.Is(err, errCapacityNegativeOrZero) {
					t.Log("zero size buffer error expected, all is ok!")
				} else {
					t.Errorf("zero length buffer should be producing an error but incorrectly returns: %s", err)
					t.Fail()
				}

			case "buffer too small":
				rb, _ := New[int](test.capacity)
				if err := rb.WriteMany([]int{1, 2, 3, 4, 5}); err != nil {
					t.Logf("failed to write to buffer: %s", err)
					t.Fail()
				}

				rb, err := rb.NewSize(2)
				if rb == nil || errors.Is(err, errCapacityResizeTooSmall) {
					t.Log("buffer size too small error expected, all is ok!")
				} else {
					t.Errorf("buffer size smaller than the number of values should produce an error but incorrectly returns: %s", err)
					t.Fail()
				}
			case "resize buffer":
				rbOld, _ := New[int](test.capacity)
				for i := 0; i < test.capacity-2; i++ {
					rbOld.Write(i)
				}
				rbNew, err := rbOld.NewSize(test.capacity - 2)
				if err != nil {
					t.Errorf("unexpected error when resizing buffer from %v to %v: %s", test.capacity, test.capacity-2, err)
					t.Fail()
				}
				if !reflect.DeepEqual(rbOld.Read(), rbNew.Read()) {
					t.Errorf("old buffer is different from the new buffer")
					t.Fail()
				}
				if rbOld.capacity == rbNew.capacity {
					t.Errorf("old ring buffer is the same capacity as the new buffer")
					t.Fail()
				}
				if rbOld.elementCount != rbNew.elementCount {
					t.Errorf("old buffer element count is different from the new buffer")
					t.Fail()
				}
				if rbOld.writeIndex != rbNew.writeIndex {
					t.Errorf("old buffer writeIndex is different from the new buffer")
					t.Fail()
				}
			case "resize buffer with no values":
				rbOld, _ := New[int](test.capacity)
				rbNew, _ := rbOld.NewSize(test.capacity - 2)
				if !reflect.DeepEqual(rbOld.Read(), rbNew.Read()) {
					t.Errorf("old empty buffer is different from the new empty buffer")
					t.Fail()
				}
				if !rbNew.IsEmpty() {
					t.Errorf("resized buffer should be empty but IsEmpty() == %v", rbNew.IsEmpty())
					t.Fail()
				}
				if rbNew.IsFull() {
					t.Errorf("resized buffer should not be full but IsFull() == %v", rbNew.IsFull())
					t.Fail()
				}
			}
		})
	}
}

func TestRead(t *testing.T) {
	t.Run("Read()", func(t *testing.T) {
		rb, _ := New[string](3)
		rb.WriteMany([]string{"a", "b", "c"})
		expected := []string{"a", "b", "c"}
		result := rb.Read()
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("incorrect result on Read(), expected %s but got %s", expected, result)
			t.Fail()
		}
		rb.Write("d")
		expected = []string{"b", "c", "d"}
		result = rb.Read()
		if !reflect.DeepEqual(expected, result) {
			t.Errorf("incorrect result on Read(), expected %s but got %s", expected, result)
			t.Fail()
		}
	})
}

func TestWrite(t *testing.T) {
	capacity := 3
	testString := []string{"test1", "test2"}
	testString2 := []string{"test3", "test4"}
	expected := []string{"test2", "test3", "test4"}

	t.Run("TestWrite()", func(t *testing.T) {
		rb, _ := New[string](capacity)

		if !rb.IsEmpty() {
			t.Errorf("buffer should be empty but it is not, expected %v but got %v", true, rb.IsEmpty())
			t.Fail()
		}

		// Write testString to buffer and check that the object state is as it should be
		for _, str := range testString {
			rb.Write(str)
		}
		if !reflect.DeepEqual(testString, rb.Read()) {
			t.Errorf("incorrect result on Read(), expected %s but got %s", testString, rb.Read())
			t.Fail()
		}
		if rb.elementCount != len(testString) {
			t.Errorf("incorrect element count, expected %d but got %d", len(testString), rb.elementCount)
			t.Fail()
		}
		if rb.IsFull() {
			t.Errorf("buffer isFull when it should not be, expected %v but got %v", false, rb.IsFull())
			t.Fail()
		}
		if rb.IsEmpty() {
			t.Errorf("buffer isEmpty when it should not be, expected %v but got %v", false, rb.IsEmpty())
			t.Fail()
		}

		// Now test overwriting old values of the buffer with testString2. Make sure it
		// writes and reads as expected
		for _, str := range testString2 {
			rb.Write(str)
		}
		if !reflect.DeepEqual(expected, rb.Read()) {
			t.Errorf("incorrect result on Read(), expected %s but got %s", expected, rb.Read())
			t.Fail()
		}
		if rb.elementCount != capacity {
			t.Errorf("incorrect element count, expected %d but got %d", capacity, rb.elementCount)
			t.Fail()
		}
		if !rb.IsFull() {
			t.Errorf("buffer !IsFull() when it SHOULD be FULL, expected %v but got %v", true, rb.IsFull())
			t.Fail()
		}
		if rb.IsEmpty() {
			t.Errorf("buffer IsEmpty() when it should NOT be, expected %v but got %v", false, rb.IsEmpty())
			t.Fail()
		}
	})
}

func TestWriteMany(t *testing.T) {
	t.Run("TestWriteMany()", func(t *testing.T) {
		rb, _ := New[string](3)
		empty := []string{}
		tooManyValues := []string{"a", "b", "c", "d"}

		// test an empty string or value
		if err := rb.WriteMany(empty); err == nil {
			t.Errorf("an error was expected when zero values are written; %s", err)
			t.Fail()
		}
		// test writing too many values
		if err := rb.WriteMany(tooManyValues); err == nil {
			t.Errorf("an error was expected when too values are written; %s", err)
			t.Fail()
		}
		err := rb.WriteMany([]string{"1", "2", "3"})
		if err != nil {
			t.Errorf("an error was not expected when writing values: %s", err)
			t.Fail()
		}
	})
}

func TestString(t *testing.T) {
	tests := []struct {
		name       string
		bufferType string
		capacity   int
	}{
		{"String()<-string", "string", 2},
		{"String()<-int", "int", 2},
		{"String()<-byte", "byte", 2},
		{"String()<-float", "float", 2},
		{"String()<-bool", "bool", 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.bufferType {
			case "string":
				rb, _ := New[string](test.capacity)
				err := rb.WriteMany([]string{"a", "b"})
				if err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				t.Logf("buffer[%s]: %s", test.bufferType, rb.String())
			case "int":
				rb, _ := New[int](test.capacity)
				err := rb.WriteMany([]int{1, 2})
				if err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				t.Logf("buffer[%s]: %s", test.bufferType, rb.String())
			case "uint":
				rb, _ := New[uint](test.capacity)
				err := rb.WriteMany([]uint{1, 2})
				if err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				t.Logf("buffer[%s]: %s", test.bufferType, rb.String())
			case "byte":
				rb, _ := New[byte](test.capacity)
				err := rb.WriteMany([]byte{1, 2})
				if err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				t.Logf("buffer[%s]: %s", test.bufferType, rb.String())
			case "float":
				rb, _ := New[float32](test.capacity)
				err := rb.WriteMany([]float32{0.5, 1.5})
				if err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				t.Logf("buffer[%s]: %s", test.bufferType, rb.String())
			case "bool":
				rb, _ := New[bool](test.capacity)
				err := rb.WriteMany([]bool{true, true})
				if err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				t.Logf("buffer[%s]: %s", test.bufferType, rb.String())
			}
		})
	}
}

func TestReset(t *testing.T) {
	t.Run("Reset()", func(t *testing.T) {
		rb, _ := New[string](5)
		testString := []string{"a", "b", "c", "d", "e"}
		if err := rb.WriteMany(testString); err != nil {
			t.Errorf("failed to write to buffer: %s", err)
			t.Fail()
		}
		rb.Reset()
		if reflect.DeepEqual(testString, rb.Read()) {
			t.Errorf("incorrect result on reset, expected %s but got %s", []string{}, rb.Read())
			t.Fail()
		}
		if rb.elementCount != 0 {
			t.Errorf("incorrect element count, expected %d but got %d", 0, rb.elementCount)
			t.Fail()
		}
		if rb.writeIndex != 0 {
			t.Errorf("incorrect writeIndex, expected %d but got %d", 0, rb.writeIndex)
			t.Fail()
		}
		if !rb.IsEmpty() {
			t.Errorf("buffer is not empty when it should be, expected %v but got %v", true, rb.IsEmpty())
			t.Fail()
		}
		if rb.IsFull() {
			t.Errorf("buffer isFull when it should NOT be, expected %v but got %v", false, rb.IsFull())
			t.Fail()
		}
	})
}

func TestLength(t *testing.T) {
	tests := []struct {
		name         string
		values       []string
		elementCount int
	}{
		{"Length() test 1", []string{"a"}, 1},
		{"Length() test 2", []string{"b", "c"}, 3},
	}
	rb, _ := New[string](3)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := rb.WriteMany(test.values); err != nil {
				t.Errorf("failed to write to buffer: %s", err)
				t.Fail()
			}
			if rb.elementCount != rb.Length() {
				t.Errorf("incorrect element count, expected %d but got %d", test.elementCount, rb.Length())
				t.Fail()
			}
		})
	}
}

func TestSize(t *testing.T) {
	t.Run("Size()", func(t *testing.T) {
		rb, _ := New[string](3)
		if rb.Size() != 3 {
			t.Errorf("incorrect capacity, expected %d but got %d", 3, rb.Size())
			t.Fail()
		}
		rb, _ = rb.NewSize(5)
		if rb.Size() != 5 {
			t.Errorf("incorrect capacity, expected %d but got %d", 5, rb.Size())
			t.Fail()
		}
	})
}

func TestIsFull(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"normal"},
		{"after Reset()"},
	}
	rb, _ := New[string](3)
	testStr := []string{"a", "b", "c"}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.name {
			case "normal":
				if rb.IsFull() {
					t.Errorf("incorrect IsFull() state, expected %v but got %v", false, rb.IsFull())
					t.Fail()
				}
				if err := rb.WriteMany(testStr); err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				if !rb.IsFull() {
					t.Errorf("incorrect isFull() state, expected %v but got %v", true, rb.IsFull())
					t.Fail()
				}
			case "after Reset()":
				rb.Reset()
				if err := rb.WriteMany(testStr); err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				if !rb.IsFull() {
					t.Errorf("incorrect isFull() state, expected %v but got %v", true, rb.IsFull())
					t.Fail()
				}
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	t.Run("IsEmpty()", func(t *testing.T) {
		rb, _ := New[string](3)
		if !rb.IsEmpty() {
			t.Errorf("incorrect IsEmpty() state, expected %v but got %v", true, rb.IsEmpty())
			t.Fail()
		}
		rb.Write("a")
		if rb.IsEmpty() {
			t.Errorf("incorrect IsEmpty() state, expected %v but got %v", false, rb.IsEmpty())
			t.Fail()
		}

		rb.Reset()
		if !rb.IsEmpty() {
			t.Errorf("incorrect IsEmpty() state, expected %v but got %v", true, rb.IsEmpty())
			t.Fail()
		}
	})
}
