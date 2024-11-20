package ringbuffer

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

// BufferType provides constraints on the types that may be used for a New RingBuffer
type BufferType interface {
	int | int16 | int32 | int64 |
	byte | uint | uint16 | uint32 | uint64 |
	float32 | float64 | bool | string
}

// RingBuffer is effectively a fixed-size container as a data structure. Fields defined
// in this struct are named in the context as a fixed-size container.
type RingBuffer[T BufferType] struct {
	// buffer contains all data, including undefined elements, of which take a default
	// value to their respective type. Bool defaults to False, int defaults to 0, string
	// defaults to "", ... etc
	buffer       []T
	mut          sync.Mutex // Handles thread safety and concurrency
	capacity     int        // Total size of the buffer
	elementCount int        // Number of values stored within the buffer
	writeIndex   int        // The next index to write into the buffer when Write() is called
	isFull       bool       // true when the buffer capacity matches elementCount
	isEmpty      bool       // true when the buffer has been Reset() or elementCount == 0
}

// Error handling statements
var (
	errBufferSizeIsZero = errors.New("failed to create a new ring buffer! " +
		"Capacity / size cannot be zero")
	errBufferSizeTooSmall = errors.New("failed to write to buffer! Ring buffer total " +
		"capacity is too small for all values to be written")
	errDataLengthIsZero = errors.New("failed to write to buffer! The amount of " +
		"data to write is zero")
)

// New is effectively a constructor that creates a new ring buffer with a fixed,
// zero-indexed capacity and specified type constrained by the BufferType interface
func New[T BufferType](capacity int) (*RingBuffer[T], error) {
	if capacity <= 0 {
		return nil, errBufferSizeIsZero
	}
	return &RingBuffer[T]{
		buffer:   make([]T, capacity),
		capacity: capacity,
		isFull:   false,
		isEmpty:  true,
	}, nil
}

// NewSize recreates a new ring buffer with a different capacity / size, but with the same
// data as the old ring buffer. The NEW capacity cannot be smaller than the number of
// values or elements contained in the OLD buffer.
func (rb *RingBuffer[T]) NewSize(capacity int) (*RingBuffer[T], error) {
	rb.mut.Lock()
	defer rb.mut.Unlock()

	var (
		full  bool
		empty bool
	)

	if capacity <= 0 {
		return nil, errBufferSizeIsZero
	}

	if rb.elementCount > capacity {
		return nil, errBufferSizeTooSmall
	}

	if rb.elementCount == capacity {
		full = true
	} else if rb.elementCount == 0 {
		empty = true
	}

	newBuffer := make([]T, capacity)
	newBuffer = rb.buffer

	return &RingBuffer[T]{
		buffer:       newBuffer,
		capacity:     capacity,
		elementCount: rb.elementCount,
		writeIndex:   rb.writeIndex,
		isFull:       full,
		isEmpty:      empty,
	}, nil
}

// String converts the capacity, writeIndex pointer, count of elements, and contents of
// the ring buffer into a string, then returns that string
func (rb *RingBuffer[T]) String() string {
	rb.mut.Lock()
	defer rb.mut.Unlock()

	bufferStr := "capacity= " + strconv.Itoa(rb.capacity) +
		", writeIndex= " + strconv.Itoa(rb.writeIndex) +
		", elementCount= " + strconv.Itoa(rb.elementCount) +
		", isFull= " + strconv.FormatBool(rb.isFull) +
		", isEmpty= " + strconv.FormatBool(rb.isEmpty) +
		", buffer= ["
	lastElement := len(rb.buffer) - 1

	for i := range rb.buffer {
		if i != lastElement {
			bufferStr += fmt.Sprintf("%v", rb.buffer[i]) + ","
		} else {
			bufferStr += fmt.Sprintf("%v", rb.buffer[i]) + "]"
		}
	}
	return bufferStr
}

// Read returns the contents of the buffer in "First-In First-Out" (FIFO) order
func (rb *RingBuffer[T]) Read() (result []T) {
	rb.mut.Lock()
	defer rb.mut.Unlock()

	result = make([]T, 0, rb.elementCount)

	for i := 0; i < rb.elementCount; i++ {
		index := (rb.writeIndex + rb.capacity - rb.elementCount + i) % rb.capacity
		result = append(result, rb.buffer[index])
	}
	return result
}

// Write inserts one element into the thread-safe buffer, overwriting the oldest element
// if the buffer is full
func (rb *RingBuffer[T]) Write(value T) {
	rb.mut.Lock()
	defer rb.mut.Unlock()

	rb.buffer[rb.writeIndex] = value

	// rb.writeIndex acts as a logical pointer that moves forward each time Write(...)
	//	is called.
	// When writeIndex reaches the buffer capacity, it wraps around to the beginning of
	//	the buffer by using the modulo operator
	rb.writeIndex = (rb.writeIndex + 1) % rb.capacity

	// Only increment the elementCount of elements in the buffer when the buffer
	// isn't full
	if rb.elementCount < rb.capacity {
		rb.elementCount++
	}
	if rb.elementCount == rb.capacity {
		rb.isFull = true
	}
	if rb.elementCount > 0 {
		rb.isEmpty = false
	}
}

func (rb *RingBuffer[T]) WriteValues(values []T) error {
	if len(values) > rb.capacity {
		return errBufferSizeTooSmall
	} else if len(values) == 0 {
		return errDataLengthIsZero
	}

	for _, val := range values {
		rb.Write(val)
	}
	return nil
}

// Reset recreates a new ring buffer of the same exact capacity
func (rb *RingBuffer[T]) Reset() {
	rb.mut.Lock()
	defer rb.mut.Unlock()

	rb.buffer = make([]T, rb.capacity)
	rb.elementCount = 0 // there's nothing (no elements/values) in the buffer, of course
	rb.writeIndex = 0   // reset the logical pointer to the beginning of the buffer
	rb.isEmpty = true   // a zero'd array is, by definition, empty
	rb.isFull = false   // the buffer is no longer full
}

// Length returns the number of elements / values within the buffer.
//
// For getting the total capacity of the buffer, use Size()
func (rb *RingBuffer[T]) Length() int {
	rb.mut.Lock()
	defer rb.mut.Unlock()
	return rb.elementCount
}

// Size returns the zero-indexed capacity of the ring buffer itself, as opposed to the
// number of elements within the buffer.
//
// For getting the number of elements in a buffer, use Length()
func (rb *RingBuffer[T]) Size() int {
	rb.mut.Lock()
	defer rb.mut.Unlock()
	return rb.capacity
}

func (rb *RingBuffer[T]) IsFull() bool {
	rb.mut.Lock()
	defer rb.mut.Unlock()
	return rb.isFull
}

func (rb *RingBuffer[T]) IsEmpty() bool {
	rb.mut.Lock()
	defer rb.mut.Unlock()
	return rb.isEmpty
}
