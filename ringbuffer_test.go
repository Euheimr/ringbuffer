package ringbuffer

import (
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

//func TestWrite(t *testing.T) {
//	t.Run("TestWrite()", func(t *testing.T) {
//		t.Log("NOT IMPLEMENTED")
//		t.Fail()
//	})
//}

//func TestWriteValues(t *testing.T) {
//	t.Run("TestWriteValues()", func(t *testing.T) {
//		t.Log("NOT IMPLEMENTED")
//		t.Fail()
//	})
//}

func TestString(t *testing.T) {
	tests := []struct {
		name       string
		bufferType string
	}{
		{"String()<-string", "string"},
		{"String()<-int", "int"},
		{"String()<-byte", "byte"},
		{"String()<-float", "float"},
		{"String()<-bool", "bool"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			switch test.bufferType {
			case "string":
				rb, _ := New[string](2)
				rb.WriteValues([]string{"a", "b"})
				t.Log(rb.String())
			case "int":
				rb, _ := New[int](2)
				rb.WriteValues([]int{1, 2})
				t.Log(rb.String())
			case "uint":
				rb, _ := New[uint](2)
				rb.WriteValues([]uint{1, 2})
				t.Log(rb.String())
			case "byte":
				rb, _ := New[byte](2)
				rb.WriteValues([]byte{1, 2})
				t.Log(rb.String())
			case "float":
				rb, _ := New[float32](2)
				rb.WriteValues([]float32{0.5, 1.5})
				t.Log(rb.String())
			case "bool":
				rb, _ := New[bool](2)
				rb.WriteValues([]bool{true, true})
				t.Log(rb.String())
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
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("incorrect result on Read(), expected %s but got %s", expected, result)
			t.Fail()
		}
		rb.WriteValues([]string{"d"})
		expected = []string{"b", "c", "d"}
		result = rb.Read()
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("incorrect result on Read(), expected %s but got %s", expected, result)
			t.Fail()
		}
	})
}

//func TestReset(t *testing.T) {
//	t.Run("Reset()", func(t *testing.T) {
//		t.Log("NOT IMPLEMENTED")
//		t.Fail()
//	})
//}

//func TestLength(t *testing.T) {
//	t.Run("Length()", func(t *testing.T) {
//		t.Log("NOT IMPLEMENTED")
//		t.Fail()
//	})
//}

//func TestSize(t *testing.T) {
//	t.Run("Size()", func(t *testing.T) {
//		t.Log("NOT IMPLEMENTED")
//		t.Fail()
//	})
//}

//func TestIsFull(t *testing.T) {
//	t.Run("IsFull()", func(t *testing.T) {
//		t.Log("NOT IMPLEMENTED")
//		t.Fail()
//	})
//}

//func TestIsEmpty(t *testing.T) {
//	t.Run("IsEmpty()", func(t *testing.T) {
//		t.Log("NOT IMPLEMENTED")
//		t.Fail()
//	})
//}
