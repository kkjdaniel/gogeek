<p align="center">
  <img src="gogeek-logo.png" width="350" alt="GoGeek Logo">
</p>

<h1 align="center">GoGeek: BoardGameGeek API for Go</h1>

<p align="center">
  <a href="https://pkg.go.dev/github.com/kkjdaniel/gogeek/v3"><img src="https://pkg.go.dev/badge/github.com/kkjdaniel/gogeek/v3.svg" alt="Go Reference"></a>
  <a href="https://github.com/kkjdaniel/gogeek/actions/workflows/lint.yml"><img src="https://github.com/kkjdaniel/gogeek/actions/workflows/lint.yml/badge.svg" alt="Lint"></a>
  <a href="https://codecov.io/gh/kkjdaniel/gogeek"><img src="https://codecov.io/gh/kkjdaniel/gogeek/graph/badge.svg?token=W78TFFY83D" alt="codecov"></a>
  <a href="https://github.com/kkjdaniel/gogeek/actions/workflows/contract-tests.yml"><img src="https://github.com/kkjdaniel/gogeek/actions/workflows/contract-tests.yml/badge.svg" alt="Contract Tests"></a>
</p>

GoGeek is a lightweight, easy-to-use Go module designed to streamline interactions with the [BoardGameGeek API](https://boardgamegeek.com/wiki/page/BGG_XML_API2) (XML API2).

## Key Features

- **🔄 Simple Request Handling**: GoGeek abstracts the BGG API request process, allowing you to focus on utilising the data rather than managing HTTP requests.
- **🔐 Authentication Support**: Built-in support for API key and cookie-based authentication to access authenticated endpoints.
- **📄 Data Parsing**: Automatically converts and normalises XML responses from the BGG API into Go structs, so you can work with structured data effortlessly.
- **⚠️ Error Handling**: Robust error handling for common issues like network errors, rate limiting, queued requests and unexpected response formats.
- **✅ Full API Coverage**: All BGG XML API2 endpoints are supported, with automated contract tests that run against the live API weekly to detect any structural changes or drift.

## Setup

To setup GoGeek, use the following `go get` command:

```bash
go get github.com/kkjdaniel/gogeek/v3
```

## Usage

### Basic Usage

Getting started with GoGeek is easy. First, create a client, then use it to make API requests:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kkjdaniel/gogeek/v3"
	"github.com/kkjdaniel/gogeek/v3/thing"
)

func main() {
	// Create a client (all BGG endpoints require authentication,
	// see the Authentication section below)
	client := gogeek.NewClient(gogeek.APIKey("your-api-key"))

	// Query board games by BGG ID
	games, err := thing.Query(context.Background(), client, []int{13, 12, 3})
	if err != nil {
		log.Fatal(err)
	}

	for _, game := range games.Items {
		fmt.Printf("Name: %s\nYear Published: %d\n", game.Name[0].Value, game.YearPublished.Value)
	}
}
```

```
Name: CATAN
Year Published: 1995
Name: Ra
Year Published: 1999
Name: Samurai
Year Published: 1998
```

### Authentication

All BGG endpoints require authorisation — anonymous access is no longer supported by the API, so every client is constructed with credentials.

**API Key Authentication** (recommended):

```go
client := gogeek.NewClient(gogeek.APIKey("your-api-key"))
collection, err := collection.Query(ctx, client, "username")
```

Note: To get an API key you can request one via the [application form here.](https://boardgamegeek.com/applications)

**Cookie Authentication** (can access private collection data an API key may not):

```go
cookie := "bggusername=user; bggpassword=pass; SessionID=xyz"
client := gogeek.NewClient(gogeek.Cookie(cookie))
user, err := user.Query(ctx, client, "username")
```

### Rate Limiting & Retries

All clients enforce a rate limit of **2 requests per second** to comply with BoardGameGeek's API guidelines, time out requests after 30 seconds, and automatically retry queued (202) or throttled (429/503) responses with exponential backoff. All of this is configurable via client options, e.g. `gogeek.WithRateLimit(1)`, `gogeek.WithRetry(3, time.Second)`, `gogeek.WithHTTPClient(custom)`.

### Notes

- Every query takes a `context.Context`, so requests can be cancelled or given deadlines
- The `thing` query accepts any number of IDs — batches of more than 20 are split into multiple rate-limited requests automatically
- Query filters use per-package options (e.g. `plays.WithPage(2)`, `collection.WithOwned(true)`); invalid option values return an error wrapping `gogeek.ErrInvalidOption` before any request is made

## Documentation

For the full documentation please see the [GoDoc here](https://pkg.go.dev/github.com/kkjdaniel/gogeek/v3). Details on how to use each query function as well as the interfaces for each of the APIs can be found within their respective packages.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue on GitHub to help improve GoGeek.
