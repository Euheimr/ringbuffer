# Ring Buffer (circular buffer)

A ring Buffer is a fixed-size container as a data structure. A lot of ring buffer 
implementations do not allow overwrites when the buffer is full, but I wanted that functionality, so I made it.

## Example usage

```go
package main

import (
	"fmt"
	"github.com/euheimr/ringbuffer"
)

func main() {
   capacity := 3
   rb, err := ringbuffer.New[string](capacity)
   if err != nil {
      fmt.Println(err.Error())
   }

   // Write a single value
   rb.Write("test1")
   fmt.Println(rb.Read())   // rb.Read() == []string{"test1"}

   // Write multiple values
   if err = rb.WriteValues([]string{"test2", "test3", "test4"}); err != nil {
      fmt.Println(err.Error())
   }

   // Read the buffer in order of FIFO (first-in-first-out)
   fmt.Println(rb.Read())  // rb.Read() == []string{"test2", "test3", "test4"}
   
   // You can recreate the buffer with a new size, but must be equal or greater than the
   // number of existing elements. 
   // If you want to make it smaller, call rb.Reset() then rb.NewSize(...)
   rb, err = rb.NewSize(5)
   if err != nil {
	   fmt.Println(err.Error())
   }
   if err = rb.WriteValues([]string{"test5","test6"}); err != nil {
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