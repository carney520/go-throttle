# go-throttle

[![GoDoc](https://godoc.org/github.com/carney520/go-throttle?status.svg)](https://godoc.org/github.com/carney520/go-throttle) [![Build Status](https://travis-ci.org/carney520/go-throttle.svg?branch=master)](https://travis-ci.org/carney520/go-throttle)

go-throttle create a Go version throttle.
When the passed function invoked repeatedly, will only actually call the
original function at most once per every wait `time.Duration`.
Useful for rate-limiting events that occur faster than you can keep up with.

## Usage

* Go get

```shell
go get github.com/carney520/go-throttle
```

* Example

```go
package main

import (
  "github.com/carney520/go-throttle"
  "time"
)

func main() {
  db := throttle.New(500 * time.Millisecond, func() {
    rebuild()
  })

  onFileChange(db.Trigger)
  <-exit
  db.Stop()
}
```
