package ringbuffer

import (
	"errors"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		expectedType string
		capacity     int
	}{
		{"new buffer, zero capacity", "zero", 0},
		{"new buffer, string", "string", 5},
		{"new buffer, int", "int", 5},
		{"new buffer, byte", "byte", 5},
		{"new buffer, float", "float", 5},
		{"new buffer, bool", "bool", 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			switch test.expectedType {
			case "zero":
				_, err := New[string](test.capacity)
				if err == nil {
					t.Errorf("unexpected error when creating a zero length buffer: %s", err)
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
		{"zero capacity error test", 0},
		{"buffer too small error test", 5},
		{"resize buffer", 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			switch test.name {
			case "zero capacity":
				rb, err := New[int](5)
				rb, err = rb.NewSize(test.capacity)
				if !errors.Is(err, errBufferSizeIsZero) {
					t.Errorf("zero length buffer should be producing an error but incorrectly returns: %s", err)
					t.Fail()
				}
			case "buffer too small":
				rb, err := New[int](test.capacity)
				rb.WriteValues([]int{1, 2, 3, 4, 5})
				rb, err = rb.NewSize(2)
				if !errors.Is(err, errBufferSizeTooSmall) {
					t.Errorf("buffer size smaller than the number of values should produce an error but incorrectly returns: %s", err)
					t.Fail()
				}
			case "resize buffer":
				rbOld, _ := New[int](test.capacity)
				for i := 0; i < test.capacity-2; i++ {
					err := rbOld.Write(i)
					if err != nil {
						t.Errorf("failed to write to buffer: %s", err)
						t.Fail()
					}
				}
				rbNew, err := rbOld.NewSize(test.capacity - 2)
				if err != nil {
					t.Errorf("unexpected error when resizing buffer from %v to %v: %s", test.capacity, test.capacity-2, err)
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
			}
		})
	}
}

func TestRead(t *testing.T) {
	t.Run("Read()", func(t *testing.T) {
		rb, _ := New[string](3)
		rb.WriteValues([]string{"a", "b", "c"})
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

		if !rb.isEmpty {
			t.Errorf("buffer should be empty, but it is not")
			t.Fail()
		}

		// Write testString to buffer and check that the object state is as it should be
		for _, str := range testString {
			err := rb.Write(str)
			if err != nil {
				t.Errorf("test failed to write to buffer: %s", err)
				t.Fail()
			}
		}
		if !reflect.DeepEqual(testString, rb.Read()) {
			t.Errorf("incorrect result on Read(), expected %s but got %s", testString, rb.Read())
			t.Fail()
		}
		if rb.elementCount != len(testString) {
			t.Errorf("incorrect element count, expected %d but got %d", len(testString), rb.elementCount)
			t.Fail()
		}
		if rb.isFull {
			t.Errorf("buffer isFull when it should not be, expected %v but got %v", false, rb.isFull)
			t.Fail()
		}
		if rb.isEmpty {
			t.Errorf("buffer isEmpty when it should not be, expected %v but got %v", false, rb.isEmpty)
			t.Fail()
		}

		// Now test overwriting old values of the buffer with testString2. Make sure it
		// writes and reads as expected
		for _, str := range testString2 {
			err := rb.Write(str)
			if err != nil {
				t.Errorf("test failed to write to buffer: %s", err)
				t.Fail()
			}
		}
		if !reflect.DeepEqual(expected, rb.Read()) {
			t.Errorf("incorrect result on Read(), expected %s but got %s", expected, rb.Read())
			t.Fail()
		}
		if rb.elementCount != capacity {
			t.Errorf("incorrect element count, expected %d but got %d", capacity, rb.elementCount)
			t.Fail()
		}
		if !rb.isFull {
			t.Errorf("buffer !isFull when it SHOULD be FULL, expected %v but got %v", true, rb.isFull)
			t.Fail()
		}
		if rb.isEmpty {
			t.Errorf("buffer isEmpty when it should not be, expected %v but got %v", false, rb.isEmpty)
			t.Fail()
		}
	})
}

func TestWriteValues(t *testing.T) {
	t.Run("TestWriteValues()", func(t *testing.T) {
		rb, _ := New[string](3)
		empty := []string{}
		tooManyValues := []string{"a", "b", "c", "d"}

		// test an empty string or value
		if err := rb.WriteValues(empty); err == nil {
			t.Errorf("an error was expected when zero values are written; %s", err)
			t.Fail()
		}
		// test writing too many values
		if err := rb.WriteValues(tooManyValues); err == nil {
			t.Errorf("an error was expected when too values are written; %s", err)
			t.Fail()
		}
		err := rb.WriteValues([]string{"1", "2", "3"})
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
			t.Parallel()
			switch test.bufferType {
			case "string":
				rb, _ := New[string](test.capacity)
				err := rb.WriteValues([]string{"a", "b"})
				if err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				t.Log(rb.String())
			case "int":
				rb, _ := New[int](test.capacity)
				err := rb.WriteValues([]int{1, 2})
				if err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				t.Log(rb.String())
			case "uint":
				rb, _ := New[uint](test.capacity)
				err := rb.WriteValues([]uint{1, 2})
				if err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				t.Log(rb.String())
			case "byte":
				rb, _ := New[byte](test.capacity)
				err := rb.WriteValues([]byte{1, 2})
				if err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				t.Log(rb.String())
			case "float":
				rb, _ := New[float32](test.capacity)
				err := rb.WriteValues([]float32{0.5, 1.5})
				if err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				t.Log(rb.String())
			case "bool":
				rb, _ := New[bool](test.capacity)
				err := rb.WriteValues([]bool{true, true})
				if err != nil {
					t.Errorf("failed to write to buffer: %s", err)
					t.Fail()
				}
				t.Log(rb.String())
			}
		})
	}
}

func TestReset(t *testing.T) {
	t.Run("Reset()", func(t *testing.T) {
		rb, _ := New[string](5)
		testString := []string{"a", "b", "c", "d", "e"}
		if err := rb.WriteValues(testString); err != nil {
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
		if !rb.isEmpty {
			t.Errorf("buffer is not empty when it should be, expected %v but got %v", true, rb.isEmpty)
			t.Fail()
		}
		if rb.isFull {
			t.Errorf("buffer isFull when it should NOT be, expected %v but got %v", false, rb.isFull)
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
			if err := rb.WriteValues(test.values); err != nil {
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
			t.Errorf("incorrect size, expected %d but got %d", 3, rb.Size())
			t.Fail()
		}
	})
}

func TestIsFull(t *testing.T) {
	t.Run("IsFull()", func(t *testing.T) {
		rb, _ := New[string](3)
		if rb.IsFull() {
			t.Errorf("incorrect IsFull() state, expected %v but got %v", false, rb.IsFull())
		}
		if err := rb.WriteValues([]string{"a", "b", "c"}); err != nil {
			t.Errorf("failed to write to buffer: %s", err)
			t.Fail()
		}
		if rb.isFull != rb.IsFull() {
			t.Errorf("rb.IsFull() does not match the state of rb.isFull, expected %v but got %v", true, rb.IsFull())
			t.Fail()
		}
		if !rb.IsFull() {
			t.Errorf("incorrect isFull() state, expected %v but got %v", true, rb.IsFull())
			t.Fail()
		}
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("IsEmpty()", func(t *testing.T) {
		rb, _ := New[string](3)
		if !rb.IsEmpty() {
			t.Errorf("incorrect IsEmpty() state, expected %v but got %v", true, rb.IsEmpty())
			t.Fail()
		}
		if rb.isEmpty != rb.IsEmpty() {
			t.Errorf("incorrect IsEmpty() state, expected %v but got %v", true, rb.IsEmpty())
		}
		if err := rb.Write("a"); err != nil {
			t.Errorf("failed to write to buffer: %s", err)
			t.Fail()
		}
		if rb.IsEmpty() {
			t.Errorf("incorrect IsEmpty() state, expected %v but got %v", false, rb.IsEmpty())
		}
	})
}
