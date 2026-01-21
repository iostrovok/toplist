# toplist

`toplist` is a high-performance, concurrent-safe skip list implementation in Go. It is designed to handle large volumes of sorted data with support for both synchronous and asynchronous operations.

## Key Features

- **Lock-Free Base Layer**: Uses atomic operations (CAS) for thread-safe insertions and deletions without traditional mutexes in the base data layer.
- **Asynchronous Processing**: Includes a built-in worker queue for handling high-load write operations concurrently.
- **Dynamic Rebuilding**: Automatically rebalances the skip list structure to maintain $O(\log N)$ search performance as the list grows.
- **Memory Efficient**: Uses pointer marking to track logical deletions, allowing for efficient physical cleanup.
- **Debug Friendly**: Built-in debug modes and map tracking to trace operations and list consistency.

## Usage

### Initialization

```go
import "github.com/iostrovok/toplist"

// Create a new list with a background worker pool
tl := toplist.New()

```

### Basic Operations
#### Save (Insert or Update)
Inserts a new value or updates an existing one at the given index.
``` go
err := tl.Save(100, "some data")
```

#### Find
Retrieves an element by its index.


``` go
node, found := tl.Find(100)
if found {
    fmt.Printf("Value: %v\n", node.Value)
}
```

#### Delete
Removes an element from the list.

``` go
err := tl.Delete(100)
```

### Advanced Features
#### Batch/Async Processing
Use Run to enqueue operations without waiting for immediate completion, providing a callback for the result.

``` go
tl.Run(toplist.SaveAction, 200, "async data", func(action toplist.Action, index int64, err error) {
    if err == nil {
        fmt.Printf("Operation %s on index %d successful\n", action, index)
    }
})
```

#### List Maintenance
Manually trigger a rebalance or physical cleanup of marked nodes.

``` go
// Rebuild the skip list levels for optimal search speed
tl.Build()

// Physically remove nodes marked for deletion
removedCount := tl.Clean()
```

## Internal Architecture
The list maintains multiple levels of indices. While the base level is modified using atomic CAS operations, 
the upper index levels can be reconstructed periodically via the Build method (or automatically by the internal queue) 
to ensure search efficiency stays logarithmic even under heavy modification.
