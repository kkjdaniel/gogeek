# GoGeek: Go Module for Simplified BoardGameGeek API Integration

[![Go Reference](https://pkg.go.dev/badge/pkg.go.dev/github.com/kkjdaniel/gogeek.svg)](https://pkg.go.dev/github.com/kkjdaniel/gogeek)
[![Go Report Card](https://goreportcard.com/badge/github.com/kkjdaniel/gogeek)](https://goreportcard.com/report/github.com/kkjdaniel/gogeek)
[![codecov](https://codecov.io/gh/kkjdaniel/gogeek/graph/badge.svg?token=W78TFFY83D)](https://codecov.io/gh/kkjdaniel/gogeek)

GoGeek is a lightweight, easy-to-use Go module designed to streamline interactions with the [BoardGameGeek API](https://boardgamegeek.com/wiki/page/BGG_XML_API2) (XML API2).

## Key Features

- **🔄 Simple Request Handling**: GoGeek abstracts the BGG API request process, allowing you to focus on utilising the data rather than managing HTTP requests.
- **📄 Data Parsing**: Automatically converts XML responses from the BGG API into Go structs, so you can work with structured data effortlessly.
- **⚠️ Error Handling**: Robust error handling for common issues like network errors, rate limiting, and unexpected response formats ensures smooth integration.

## Setup

To setup GoGeek, use the following `go get` command:

```bash
go get github.com/kkjdaniel/gogeek
```

## Usage

Getting started with GoGeek is easy. Here’s a quick example to fetch details about a specific board game:

```go
package main

import (
    "fmt"
    "log"
    "github.com/kkjdaniel/gogeek/thing"
)

func main() {
	// Catan
	game, err := thing.Query([]int{9})
	if err != nil {
		log.Fatal(err)
	}
	catan := game.Items[0]
	fmt.Printf("Name: %s\nYear Published: %d\nRating: %.2f\n", catan.Name[0].Value, catan.YearPublished.Value, catan.Statistics.AverageRating)
}
```

## Documentation

For the full documentation please see the [GoDoc here](https://pkg.go.dev/github.com/kkjdaniel/gogeek).

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue on GitHub to help improve GoGeek.
