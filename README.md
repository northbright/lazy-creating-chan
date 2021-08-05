# lazy-creating-chan

lazy-creating-chan is an example to show how to use a lazily created channel to pass data between goroutines.

This example code is inspired by the [Done()](https://github.com/golang/go/blob/release-branch.go1.17/src/context/context.go#L358) and [cancel()](https://github.com/golang/go/blob/release-branch.go1.17/src/context/context.go#L397) in the [official context package](https://pkg.go.dev/context).
 
It uses [atomic.Value](https://pkg.go.dev/sync/atomic)'s [Load](https://pkg.go.dev/sync/atomic#Value.Load) and [Store](https://pkg.go.dev/sync/atomic#Value.Store) to load / save a channel.

The channel is created at the first ProgressCh() is called and the work gourtine will send the progress data to the cannel only if the channel is created.
