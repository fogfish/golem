# Type Safe Channels (`pipe`)

Go's concurrency features simplify the creation of streaming data pipelines that effectively utilize I/O and multiple CPUs. [Go Concurrency Patterns: Pipelines and cancellation](https://go.dev/blog/pipelines) explains the topic in-depth. Despite the simplicity, the boilerplate code is required to build channels and establish "chrome" for its management. This module offers consistent operation over channels to form a processing pipelines in a clean and type safe manner. 

The module support sequential `pipe` and parallel `fork` type safe channels  for building pipelines. The module consider channels as "sequential data structure" trying to derive semantic from [stream interface](http://srfi.schemers.org/srfi-41/srfi-41.html)


## Quick Example

Example below is most simplest illustration of composition processing pipeline.

```go
import (
  "github.com/fogfish/golem/pipe"
)

func main() {
  ctx, close := context.WithCancel(context.Background())

  // Generate sequence of integers
  ints := pipe.StdErr(pipe.Unfold(ctx, cap, 0,
    pipe.Pure(func(x int) int { return x + 1 }),
  ))

  // Limit sequence of integers
  ints10 := pipe.TakeWhile(ctx, ints,
    pipe.Pure(func(x int) bool { return x <= 10 }),
  )

  // Calculate squares
  sqrt := pipe.StdErr(pipe.Map(ctx, ints10,
    pipe.Pure(func(x int) int { return x * x }),
  ))

  // Numbers to string
  vals := pipe.StdErr(pipe.Map(ctx, sqrt,
    pipe.Pure(strconv.Itoa),
  ))

  // Output strings
  <-pipe.ForEach(ctx, vals,
    pipe.Pure(
      func(x string) string {
        fmt.Printf("==> %s\n", x)
        return x
      },
    ),
  )

  close()
}
```



## Stream interface

### Supported features
- [x] `emit` takes a function that emits data at a specified frequency to the channel.
- [x] `filter` returns a newly-allocated channel that contains only those elements X of the input channel for which predicate is true.
- [x] `foreach` applies function for each message in the channel.
- [x] `map` applies function over channel messages, emits result to new channel.
- [x] `fold` applies a monoid operation to the values in a channel. The final value is emitted though return channel when the end of the input channel is reached.
- [x] `join` concatenate channels, returns newly-allocated channel composed of elements copied from input channels. 
- [x] `partition` partitions channel in two channels according to a predicate.
- [x] `take` returns a newly-allocated channel containing the first n elements of the input channel.
- [x] `takeWhile` returns a newly-allocated channel that contains those elements from channel while predicate returns true.
- [x] `unfold` the fundamental recursive constructor, it applies a function to each previous seed element in turn to determine the next element.

  
### Not supported feature
- [ ] `drop` returns the suffix of the input channel that starts at the next element after the first n elements.
- [ ] `dropWhile` drops elements from channel while predicate returns true and returns remaining channel suffix.
- [ ] `split` partitions channel into two channels. The split behaves as if it is defined as consequent take, drop.
- [ ] `splitWhile` partitions channel into two channels according to predicate. The splitWhile behaves as if it is defined as consequent takeWhile, dropWhile.
- [ ] `flatten` reduces dimension of channel of channels.
- [ ] `scan` accumulates the partial folds of an input channel into a newly-allocated channel.
- [ ] `zip` takes one or more input channels and returns a newly-allocated channel in which each element is a product of the corresponding elements of the input channels. The output channel is as long as the shortest input stream.
- [ ] `zipWith` takes one or more input channels and returns a newly-allocated channel, each element produced by composition function that map list of input heads to new head. The output stream is as long as the longest input stream.