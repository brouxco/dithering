# Dithering

[![GoDoc](https://godoc.org/github.com/brouxco/dithering?status.svg)](https://godoc.org/github.com/brouxco/dithering) [![Go Report Card](https://goreportcard.com/badge/github.com/brouxco/dithering)](https://goreportcard.com/report/github.com/brouxco/dithering) [![GitHub license](https://img.shields.io/github/license/brouxco/dithering.svg)](https://github.com/brouxco/dithering/blob/master/LICENSE.md)

> Image dithering in go

This go library provides a general purpose dithering algorithm implementation.

The color palette and the error diffusion matrix are customizable.

### Install

In order to use this module run:
```bash
go get github.com/brouxco/dithering 
```
> Note: this may not be necessary if you use Go 1.11 or later: see [Go Modules](https://github.com/golang/go/wiki/Modules)

In your code don't forget the import:
```go
import "github.com/brouxco/dithering"
```

### License 

[MIT](https://github.com/brouxco/dithering/blob/master/LICENSE.md)
