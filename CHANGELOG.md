# Changelog

## v3.0.0

Module path is now `github.com/kkjdaniel/gogeek/v3`.

### Breaking changes

- **Context support**: every `Query` function takes a `context.Context` as its first parameter. Migration: `thing.Query(client, ids)` → `thing.Query(ctx, client, ids)`.
- **Mandatory authentication**: `NewClient` requires an `Auth` argument built with `gogeek.APIKey(...)` or `gogeek.Cookie(...)`. Migration: `gogeek.NewClient(gogeek.WithAPIKey(k))` → `gogeek.NewClient(gogeek.APIKey(k))`. `AuthNone`, `WithAPIKey`, and `WithCookie` are removed.
- **Credential getters removed**: `Client.APIKey()`, `CookieString()`, `AuthMode()`, and `Limiter()` are gone. Migration: use `Client.Prepare(ctx, req)` to apply auth and rate limiting to a raw request, and `Client.HTTPClient()` to send it.
- **Search signature**: the variadic exact flag is replaced by options. Migration: `search.Query(client, q, true)` → `search.Query(ctx, client, q, search.WithExact())`.
- **Option types renamed and validated**: `collection.CollectionOption` → `collection.Option` (same for `forum`), and options now return an error — out-of-range values (e.g. `collection.WithMinRating(11)`) fail with `gogeek.ErrInvalidOption` instead of being silently ignored.
- **`thing.Query` auto-chunks**: more than 20 IDs are split into sequential batches of 20 and merged; `thing.ErrTooManyIDs` is removed.
- **`request` and `testutils` moved under `internal/`**: they were only ever intended for GoGeek's own use.
- **Invalid `hot.Query` item types** now return `gogeek.ErrInvalidOption` instead of hitting the API.

### New

- Injectable HTTP client with a 30-second default timeout (`gogeek.WithHTTPClient`).
- Configurable rate limit and retry behaviour (`gogeek.WithRateLimit`, `gogeek.WithRetry`).
- 429/503 responses are retried with `Retry-After` support and exponential backoff with jitter; other non-200 statuses fail fast.
- New endpoint parameters, all verified against the live API:
  - `plays`: `WithPage`, `WithMinDate`, `WithMaxDate`, `WithID`, `WithType`, `WithSubtype`
  - `guild`: `WithMembers`, `WithPage`, `WithSort` (plus the `Members` roster on the model)
  - `thread`: `WithMinArticleID`, `WithMinArticleDate`, `WithCount`
  - `search`: `WithType`
  - `user`: `WithoutBuddies`, `WithoutGuilds`, `WithoutTop`, `WithHot`, `WithDomain`, `WithPage` (plus the `Hot` list on the model)
- Doc comments on all exported types; the revive `exported` lint rule is enforced.

### Removed dependencies

- `clbanning/mxj` — the fallback XML re-parser never fired against the live API.
- `uber-go/ratelimit` and `benbjohnson/clock` — replaced by `golang.org/x/time/rate`, whose `Wait(ctx)` is context-aware.
