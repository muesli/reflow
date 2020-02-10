# reflow

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/muesli/reflow)
[![Build Status](https://travis-ci.org/muesli/reflow.svg?branch=master)](https://travis-ci.org/muesli/reflow)
[![Coverage Status](https://coveralls.io/repos/github/muesli/reflow/badge.svg?branch=master)](https://coveralls.io/github/muesli/reflow?branch=master)
[![Go ReportCard](http://goreportcard.com/badge/muesli/reflow)](http://goreportcard.com/report/muesli/reflow)

A collection of ANSI-aware methods and `io.Writers` helping you to transform
blocks of text. This means you can still style your terminal output with ANSI
escape sequences without them affecting the reflow operations & algorithms.

## Word-Wrapping

The `wordwrap` package lets you word-wrap strings or entire blocks of text.

### Usage

```go
import "github.com/muesli/reflow/wordwrap"

s := wordwrap.String("Hello World!", 5)
fmt.Println(s)
```

Result:
```
Hello
World!
```

The word-wrapping Writer is compatible with the `io.Writer` / `io.WriteCloser` interfaces:

```go
f := wordwrap.NewWriter(limit)
f.Write(b)
f.Close()

fmt.Println(f.String())
```

Customize word-wrapping behavior:

```go
f := wordwrap.NewWriter(limit)
f.Breakpoints = []rune{':', ','}
f.Newline = []rune{'\r'}
```

### ANSI Example

```go
s := wordwrap.String("I really \x1B[38;2;249;38;114mlove\x1B[0m Go!", 8)
fmt.Println(s)
```

Result:

![ANSI Example Output](https://github.com/muesli/reflow/blob/master/reflow.png)

## Indentation

The `indent` package lets you indent strings or entire blocks of text.

### Usage

```go
import "github.com/muesli/reflow/indent"

s := indent.String("Hello World!", 4)
fmt.Println(s)
```

Result: `    Hello World!`

There is also an indenting Writer, which is compatible with the `io.Writer`
interface:

```go
// indent uses spaces per default:
f := indent.NewWriter(width, nil)

// but you can also use a custom indentation function:
f = indent.NewWriter(width, func(w io.Writer) {
    w.Write([]byte("."))
})

f.Write(b)
f.Close()

fmt.Println(f.String())
```

## Padding

The `padding` package lets you pad strings or entire blocks of text.

### Usage

```go
import "github.com/muesli/reflow/padding"

s := padding.String("Hello", 8)
fmt.Println(s)
```

Result: `Hello   `

There is also a padding Writer, which is compatible with the `io.WriteCloser`
interface:

```go
// padding uses spaces per default:
f := padding.NewWriter(width, nil)

// but you can also use a custom padding function:
f = padding.NewWriter(width, func(w io.Writer) {
    w.Write([]byte("."))
})

f.Write(b)
f.Close()

fmt.Println(f.String())
```
