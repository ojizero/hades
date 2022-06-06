# gracefully

> Handling your Go servers with grace.

## Usage

Install the package

    go get github.com/ojizero/hades/gracefully

Replace your application's startup point with the following

```go
import (
  "net/http"

  "github.com/ojizero/hades/gracefully"
)

func main() {
  var h http.Handler // ... build your http handler
  var addr string    // set the address to listen to
  gracefully.ServeHandler(h, addr)
}
```

And in case you have an `http.Server` instance you can use

```go
import (
  "net/http"

  "github.com/ojizero/hades/gracefully"
)

func main() {
  var srv http.Server // ... build your http server
  gracefully.Serve(srv)
}
```
