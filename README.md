# Ring Buffer (circular buffer)
[![Go](https://github.com/Euheimr/ringbuffer/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/Euheimr/ringbuffer/actions/workflows/go.yml) [![Coverage Status](https://coveralls.io/repos/github/Euheimr/ringbuffer/badge.svg?branch=master)](https://coveralls.io/github/Euheimr/ringbuffer?branch=master) [![License](https://img.shields.io/:license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A ring Buffer is a fixed-size container as a data structure. A lot of ring buffer 
implementations do not allow overwrites when the buffer is full, but I wanted that functionality, so I made it.

## Example usage

First, get the package:

   `go get github.com/euheimr/ringbuffer`

Example usage:

```go
package main

import (
	"fmt"
	"github.com/euheimr/ringbuffer"
)

func main() {
   capacity := 3
   // Create a new ring buffer of type `string` and fixed `capacity` or size. The capacity
   //  defines how many elements are in the buffer. 
   // Please note that the buffer itself is zero-indexed. In other words, with a capacity
   //  of 3, the last index of the buffer is 2, and the first element is 0.
   rb, err := ringbuffer.New[string](capacity)
   if err != nil {
      fmt.Println(err.Error())
   }

   // Write a single value
   rb.Write("test1")
   fmt.Println(rb.Read())   // rb.Read() == []string{"test1"}

   // Write multiple values
   if err = rb.WriteMany([]string{"test2", "test3", "test4"}); err != nil {
      fmt.Println(err.Error())
   }

   // Read the buffer in order of FIFO (first-in-first-out)
   fmt.Println(rb.Read())  // rb.Read() == []string{"test2", "test3", "test4"}
   
   // You can recreate the buffer with a new size! However, it MUST be EQUAL or GREATER 
   // than the number of EXISTING elements. 
   // If you want to make it smaller than the existing buffer, you must call rb.Reset() 
   // to clear the buffer data, then rb.NewSize(...) with the smaller size capacity
   rb, err = rb.NewSize(5)
   if err != nil {
	   fmt.Println(err.Error())
   }
   if err = rb.WriteMany([]string{"test5","test6"}); err != nil {
	   fmt.Println(err.Error())
   }

   // newRb.Read() == []string{"test2", "test3", "test4", "test5", "test6"}
   fmt.Println(rb.Read())
   
   // Reset the buffer
   rb.Reset()
   fmt.Println(rb.Read())   // rb.Read() == []string{}
}
```


## Testing & Coverage

Code coverage is at: **100%**

1. Script: `sh ./coverage.sh`
2. Or manually:

   a. `go test -coverprofile coverage.out ./...`
   
   b. `go tool cover -html coverage.out -o coverage.html`
3. Then open up `coverage.html` in a web browser

You can also verbosely print out all tests in a terminal by running:

   `go test -cover -v`