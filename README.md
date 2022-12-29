# gomapper

[![GoDoc](https://godoc.org/github.com/Trendyol/gomapper?status.svg)](https://godoc.org/github.com/Trendyol/gomapper)

An auto mapping tool for Golang dtos.
(Originally inspired by https://github.com/stroiman/go-automapper)

Examples:

Declare your types:

```go
    type Location struct {
        Name string
        zone []string
    }

    type Source struct {
        Ui64     uint64
        i64      int64
        Location Location
    }

    type Destination struct {
        Ui64     uint64
        i64      int64
        Location *Location
    }
```

Use Map function:

```go
    src := Source{
        Ui64: 123,
        i64:  321,
        Location: Location{
	        Name: "abc",
	        zone: []string{"k", "l", "m"},
	    },
    }

    dest := Destination{}

    if err := gomapper.Map(src, &dest); err != nil {
        // handle mapping error
    }

    // You will see i64 field in dest is equal to zero.
    // Because private fields won't map until source
    // and destination types are the same.
    // (src and dest variables are two different types)
```
