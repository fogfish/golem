# Type Safe Channels (`pipe`)

The module defines type safe channels with few operations over it.

## Quick Example

```go
// Create unbounded channel
recv, send := pipe.New[string](0)
defer close(send) 

// Applies function over channel
ch := pipe.Map(recv, func(s string) int {return len(s)})
```
