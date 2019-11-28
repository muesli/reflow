# reflow
Reflow lets you word-wrap strings or entire blocks of text.
It conveniently follows the `io.Writer` / `io.WriteCloser` interface and
supports ANSI escape sequences. This means you can style your terminal output
without it affecting the word-wrapping algorithm.

## Usage

```go
s := reflow.ReflowString("Hello World!", 5)
fmt.Println(s)
```

Result:
```
Hello
World!
```

You can also customize reflow's behavior:

```go
f := reflow.NewReflow(limit)
f.Breakpoints = []rune{':', ','}
f.Newline = []rune{'\r'}

f.Write(b)
f.Close()

fmt.Println(f.String())
```

## ANSI Example

```go
    s := reflow.ReflowString("I really \x1B[38;2;249;38;114mlove\x1B[0m Go!")
    fmt.Println(s)
```

Result:

![ANSI Example Output](https://github.com/muesli/reflow/blob/master/reflow.png)

## Development

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/muesli/reflow)
[![Build Status](https://travis-ci.org/muesli/reflow.svg?branch=master)](https://travis-ci.org/muesli/reflow)
[![Coverage Status](https://coveralls.io/repos/github/muesli/reflow/badge.svg?branch=master)](https://coveralls.io/github/muesli/reflow?branch=master)
[![Go ReportCard](http://goreportcard.com/badge/muesli/reflow)](http://goreportcard.com/report/muesli/reflow)
