package ringbuffer

import (
	"fmt"
	"strconv"
	"sync"
)

// BufferType provides constraints on the types that may be used for a NewRingBuffer
type BufferType interface {
	string | int | int32 | int64 | float32 | float64 | bool
}

type RingBuffer[T BufferType] struct {
	buffer []T
	mut    sync.Mutex
	size   int
	write  int
	count  int
}

// NewRingBuffer creates a new ring buffer with a fixed size and specified type
// constrained by the BufferType interface
func NewRingBuffer[T BufferType](size int) *RingBuffer[T] {
	return &RingBuffer[T]{
		buffer: make([]T, size),
		size:   size,
	}
}

// Add inserts a new element into the thread-safe buffer, overwriting the oldest element
// if the buffer is full
func (rb *RingBuffer[T]) Add(value T) {
	rb.mut.Lock()
	defer rb.mut.Unlock()

	rb.buffer[rb.write] = value
	// rb.write acts as a logical pointer that moves forward each time Add(...) is called.
	// When write reaches the buffer size (rb.size), it wraps around to the beginning of
	//	the buffer by using the modulo operator
	rb.write = (rb.write + 1) % rb.size

	if rb.count < rb.size {
		// Only increment the count of elements in the buffer when the buffer isn't full
		rb.count++
	}
}

// Get returns the contents of the buffer in "First-In First-Out" (FIFO) order
func (rb *RingBuffer[T]) Get() []T {
	rb.mut.Lock()
	defer rb.mut.Unlock()

	result := make([]T, rb.count)
	for i := 0; i < rb.count; i++ {
		index := (rb.write + rb.size - rb.count + i) % rb.size
		result = append(result, rb.buffer[index])
	}
	return result
}

// Len returns the number of elements in the thread-safe buffer
func (rb *RingBuffer[T]) Len() int {
	rb.mut.Lock()
	defer rb.mut.Unlock()
	return rb.count
}

// Reset recreates a new ring buffer of the same exact size
func (rb *RingBuffer[T]) Reset() {
	rb.mut.Lock()
	defer rb.mut.Unlock()

	rb.buffer = make([]T, rb.size)
	rb.write = 0 // reset the write logical pointer to the start of the buffer
	rb.count = 0 // there's nothing in the buffer, of course
}

// Size returns the size of the ring buffer itself, as opposed to the number of elements
// within the buffer
func (rb *RingBuffer[T]) Size() int {
	rb.mut.Lock()
	defer rb.mut.Unlock()
	return rb.size
}

// String converts the size, write pointer, count of elements, and contents of the ring
// buffer into a string, then returns that string
func (rb *RingBuffer[T]) String() string {
	rb.mut.Lock()
	defer rb.mut.Unlock()

	bufferStr := "size= " + strconv.Itoa(rb.size) +
		", write= " + strconv.Itoa(rb.write) +
		", count= " + strconv.Itoa(rb.count) +
		", buffer= {"
	lastElement := len(rb.buffer) - 1

	for i := range rb.buffer {
		// FYI: https://ectobit.com/blog/check-type-of-generic-parameter/
		// Checking the type of generic parameter is kinda weird in Go ...
		switch any(rb.buffer).(type) {
		case string:
			if i == lastElement {
				bufferStr += fmt.Sprintf("%v", rb.buffer[i])
			} else {
				bufferStr += fmt.Sprintf("%v", rb.buffer[i]) + ","
			}
		case float32, float64:
			if i == lastElement {
				bufferStr += fmt.Sprintf("%.2f", rb.buffer[i])
			} else {
				bufferStr += fmt.Sprintf("%.2f", rb.buffer[i]) + ","
			}
		case int, int32, int64:
			if i == lastElement {
				bufferStr += fmt.Sprintf("%d", rb.buffer[i])
			} else {
				bufferStr += fmt.Sprintf("%d", rb.buffer[i]) + ","
			}
		case bool:
			if i == lastElement {
				bufferStr += fmt.Sprintf("%t", rb.buffer[i])
			} else {
				bufferStr += fmt.Sprintf("%t", rb.buffer[i]) + ","
			}
		}
	}
	bufferStr += "}"
	return bufferStr
}
