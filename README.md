## rotatewriter
A no-frills io.Writer interface for automatically rotating files. Separated from [auralog](https://github.com/SaintWish/auralog) since it's being deprecated due to slog.

## Documentation
[![GoDoc](https://godoc.org/github.com/SaintWish/rotatewriter?status.svg)](https://godoc.org/github.com/SaintWish/rotatewriter)

## Getting Started
``go get github.com/saintwish/rotatewriter`` Will get the latest release.

## Example
```go
package main

import(
  "time"

  "github.com/saintwish/rotatewriter"
)

func main() {
  fileio := &rotatewriter.RotateWriter{
    Dir: "./runtime/logs/", //Required
    Filename: "log.log", //Required
    ExpireTime: 24 * time.Hour, //Required if you want daily log rotation.
    MaxSize: 5 * rotatewriter.Megabyte, //Can also use auralog.Kilobyte, Not required
  }
}
```

## Contributing
If you have any additions or contributions you would like to make please do. Just keep the below in mind.
* Try to match current naming conventions as closely as possible.
* Create a Pull Request with your changes against the master branch.

## License
It's licensed under BSD-2-Clause License to stay compatible auralogs' and default go packages.
