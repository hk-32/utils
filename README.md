# utils
This is a lightweight replacement for some std-lib packages, reduces the binary size by roughly 400kb in a hello world program.

*Please note: This package achieves smaller binaries primarily by not relying on std-lib reflection. So if your program does use reflection itself then this won't benefit you much.*

## Getting Started

### Installing

To start using utils, install Go and run `go get`:

```sh
$ go get -u github.com/hk-32/utils
```

This will retrieve the library. Specifically the v1.0.0 right now. Works perfecly fine with modules.

### Examples
Hello World:

```go
package main

import "github.com/hk-32/utils/out"

func main() {
    out.Println("Hello World")
}
```

Or get some user input:

```go
package main

import (
    "github.com/hk-32/utils/out"
    "github.com/hk-32/utils/in"
)

func main() {
    // name := in.ReadLine()
    name := in.Input("Please enter your name: ")
    out.WriteLine("Hello", name)
}
```

Or check the type of `interface{}` values:

```go
package main

import (
    "github.com/hk-32/utils/out"
    "github.com/hk-32/utils/kind"
)

func main() {
    var x interface{} = 50
    
    if kind.Of(x) == kind.Int {
        out.WriteLine("Yes!")
    } else {
        out.WriteLine("No!")
    }
}
```

Maybe even pick random values:

```go
package main

import (
    "github.com/hk-32/utils/out"
    "github.com/hk-32/utils/random"
)

func main() {
    out.WriteLine(random.Pick(0, 1, 2, 3, 4, 5, 6, 7, 8, 9))
    out.WriteLine(random.Pick(0, 1, 2, 3, 4, 5, 6, 7, 8, 9))
    out.WriteLine(random.Pick(0, 1, 2, 3, 4, 5, 6, 7, 8, 9))
    out.WriteLine(random.Pick(0, 1, 2, 3, 4, 5, 6, 7, 8, 9))
}
```

## Contact

Hassan Khan: HK.32@outlook.com

## License
`utils` source code is available under the MIT [License](/LICENSE).
