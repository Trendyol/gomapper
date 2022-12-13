# gomapper

[![GoDoc](https://godoc.org/github.com/Trendyol/gomapper?status.svg)](https://godoc.org/github.com/Trendyol/gomapper)
An auto mapping tool for Golang dtos.
(Originally inspired by https://github.com/stroiman/go-automapper)

Examples:

Declare your types:
```go
    type A struct {
    	Name string
    	zone []string
    }

    type X struct {
    	Ui64 uint64
    	i64  int64
    	A    A
    }

    type Y struct {
    	Ui64 uint64
    	i64  int64
    	A    *A
    }
```

Use Map function:
```go
	source := X{
		Ui64: 123,
		i64:  321,
		A: A{
			Name: "abc",
			zone: []string{"k", "l", "m"},
		},
	}

	dest := &Y{}

	if err := Map(source, dest); err != nil {
        // handle mapping error
    }

    // You will see i64 field in Y equal to zero.
    // Because private fields won't map until source and dest types are the same. (X and Y are different)
```