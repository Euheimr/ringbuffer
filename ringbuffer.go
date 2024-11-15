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
	// All data, including undefined elements, of which take a default value to their
	// respective type. Bool defaults to False, int defaults to 0, string defaults to
	// "", ... etc
	buffer     []T
	mut        sync.Mutex // Handles thread safety and concurrency
	capacity   int        // Total size of the buffer
	values     int        // Number of elements stored within the buffer
	writeIndex int        // The next index to write into the buffer when Write() is called
	isFull     bool       // Flag to tell when the buffer has been overwritten at some point
	isEmpty    bool       // Flag for telling when the buffer has been Reset() or the buffer was newly created but has no stored values
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

// String converts the capacity, writeIndex pointer, count of elements, and contents of
// the ring buffer into a string, then returns that string
func (rb *RingBuffer[T]) String() string {
	rb.mut.Lock()
	defer rb.mut.Unlock()

	bufferStr := "capacity= " + strconv.Itoa(rb.capacity) +
		", writeIndex= " + strconv.Itoa(rb.writeIndex) +
		", values= " + strconv.Itoa(rb.values) +
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

	result = make([]T, 0, rb.values)

	for i := 0; i < rb.values; i++ {
		index := (rb.writeIndex + rb.capacity - rb.values + i) % rb.capacity
		result = append(result, rb.buffer[index])
	}
	return result
}

// Write inserts one element into the thread-safe buffer, overwriting the oldest element
// if the buffer is full
func (rb *RingBuffer[T]) Write(value T) error {
	rb.mut.Lock()
	defer rb.mut.Unlock()

	rb.buffer[rb.writeIndex] = value
	// rb.writeIndex acts as a logical pointer that moves forward each time Write(...)
	//	is called.
	// When writeIndex reaches the buffer capacity, it wraps around to the beginning of
	//	the buffer by using the modulo operator
	rb.writeIndex = (rb.writeIndex + 1) % rb.capacity

	if rb.values < rb.capacity {
		// Only increment the values of elements in the buffer when the buffer isn't full
		rb.values++
	}
	if rb.values == rb.capacity {
		rb.isFull = true
	} else if rb.values == 0 && rb.writeIndex == rb.values {
		rb.isEmpty = true
	} else {
		rb.isEmpty = false
	}
	return nil
}

func (rb *RingBuffer[T]) WriteValues(values []T) error {
	if len(values) > rb.capacity {
		return errBufferSizeTooSmall
	} else if len(values) == 0 {
		return errDataLengthIsZero
	}

	for _, val := range values {
		if err := rb.Write(val); err != nil {
			return err
		}
	}
	return nil
}

// Reset recreates a new ring buffer of the same exact capacity
func (rb *RingBuffer[T]) Reset() {
	rb.mut.Lock()
	defer rb.mut.Unlock()

	rb.buffer = make([]T, rb.capacity)
	rb.values = 0     // there's nothing (no elements/values) in the buffer, of course
	rb.writeIndex = 0 // reset the logical pointer to the beginning of the buffer
	rb.isEmpty = true // a zero'd array is, by definition, empty
	rb.isFull = false // the buffer is no longer full
}

// Length returns the number of elements / values within the buffer.
//
// For getting the total capacity of the buffer, use Size()
func (rb *RingBuffer[T]) Length() int {
	rb.mut.Lock()
	defer rb.mut.Unlock()
	return rb.values
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
