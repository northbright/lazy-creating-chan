# lazy-creating-chan

lazy-creating-chan is an example to show how to use a lazily created channel to pass data between goroutines.

There're 2 goroutines in the example: one main goroutine and one work goroutine.
It uses a channel to send / receive the progress data:
* The work goroutine sends progress data to the channel if it's not nil
* The work goroutine closes the channel if the it's not nil when task is done
* The channel will be lazily created at the first ProgressCh() is called in main goroutine only

This example code is inspired by the [Done()](https://github.com/golang/go/blob/release-branch.go1.17/src/context/context.go#L358) and [cancel()](https://github.com/golang/go/blob/release-branch.go1.17/src/context/context.go#L397) in the [official context package](https://pkg.go.dev/context).
* Use functions([Done()](https://github.com/golang/go/blob/release-branch.go1.17/src/context/context.go#L358) and [cancel()](https://github.com/golang/go/blob/release-branch.go1.17/src/context/context.go#L397)) to return the channels which make it possible to create the channel dynamically when need
* Use [atomic.Value](https://pkg.go.dev/sync/atomic)'s [Load](https://pkg.go.dev/sync/atomic#Value.Load) and [Store](https://pkg.go.dev/sync/atomic#Value.Store) to load / store a channel
* Use a [pre-closed channel](https://github.com/golang/go/blob/release-branch.go1.17/src/context/context.go#L333) to make it possible to create the channel even the task goroutine exited

  ```
  // closedchan is a reusable closed channel.
  var closedchan = make(chan struct{})

  func init() {
	  close(closedchan)
  }
  ```

