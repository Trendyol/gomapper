# gomapper

[![GoDoc](https://godoc.org/github.com/Trendyol/gomapper?status.svg)](https://godoc.org/github.com/Trendyol/gomapper)

An auto mapping tool for Golang DTOs. Basically it maps DTOs by field names.\
(Originally inspired by https://github.com/stroiman/go-automapper)

### Important!
Please be very careful when using this tool for deep copying of objects. This tool is not intended for deep copying. While moving the source to the destination, it carries the same reference (pointer) types only by reference. It would only make sense to use it as a data transfer object (DTO) between software layers.

### Example:
```go
    // Declare your types:

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

    // Declare your variables:
    
    src := Source{
        Ui64: 123,
        i64:  321,
        Location: Location{
	        Name: "abc",
	        zone: []string{"k", "l", "m"},
	    },
    }

    dest := Destination{}

    // Use Map function:

    if err := gomapper.Map(src, &dest); err != nil {
        // handle mapping error
    }

    // You will see i64 field in the dest object is equal to zero.
    // Because private fields won't map until source
    // and destination types are the same.
    // (src and dest variables are two different types)
```
