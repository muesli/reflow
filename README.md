# reflow

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/muesli/reflow)
[![Build Status](https://travis-ci.org/muesli/reflow.svg?branch=master)](https://travis-ci.org/muesli/reflow)
[![Coverage Status](https://coveralls.io/repos/github/muesli/reflow/badge.svg?branch=master)](https://coveralls.io/github/muesli/reflow?branch=master)
[![Go ReportCard](http://goreportcard.com/badge/muesli/reflow)](http://goreportcard.com/report/muesli/reflow)

A collection of methods and `io.Writers` helping you to transform blocks of
text.

## Word-Wrapping

The `wordwrap` package lets you word-wrap strings or entire blocks of text.
Conveniently it follows the `io.Writer` / `io.WriteCloser` interface and
supports ANSI escape sequences. This means you can still style your terminal
output without affecting the word-wrapping algorithm.

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
