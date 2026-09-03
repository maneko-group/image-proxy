---
name: go-naming-conventions
description: Go naming rules and conventions for identifiers, packages, files, methods, receivers, getters/setters, and interfaces (based on Alex Edwards' practical guide). ALWAYS apply when writing, editing, reviewing, or refactoring Go code — any .go file in apps/, libs/, cli/, or elsewhere. Triggers on Go/Golang work, new packages, new files, new types/functions/variables, code review.
metadata:
  author: Alex Edwards
  version: "1.0.0"
  source: https://www.alexedwards.net/blog/go-naming-conventions
---

# Go Naming Conventions

Condensed from [Alex Edwards — Go Naming Conventions: A Practical Guide](https://www.alexedwards.net/blog/go-naming-conventions). Apply these rules to **all** Go code you write or review.

## Identifiers

### Hard rules (compiler-enforced)

- Identifiers contain unicode letters, digits, and underscores only.
- Identifiers cannot begin with a digit.
- Go keywords (`break`, `func`, `interface`, `range`, `type`, `var`, etc.) cannot be identifiers.

### Conventions (mandatory in this repo)

- **Casing:** `camelCase` for unexported, `PascalCase` for exported. Never `snake_case`,
  `Pascal_Snake_Case`, `SCREAMING_SNAKE_CASE`, or `ALLUPPERCASE`.
- **Acronyms/initialisms use consistent case within the identifier:** `apiKey` / `APIKey` are
  correct, `ApiKey` is not. Same for `ID` as shorthand for identity/identifier: `userID`, never
  `userId`.
- **ASCII only:** `pi` not `π`, `resumeCount` not `résuméCount`.
- **No clashes with builtin types/functions:** don't name variables `int`, `bool`, `any` or
  functions `min`, `max`, `len`, `clear`, `cap`, `new`, `make`.
- **No type in the name:** `count` not `intCount`, `results` not `resultSlice`. Exception:
  distinguishing a converted value, e.g. `userID := 42; userIDStr := strconv.Itoa(userID)`.
- **Avoid stdlib package-name clashes** for packages your code actually imports: if you import
  `url` and `net/mail`, don't use `url` or `mail` as identifiers.

| Bad | Reason | Better |
|---|---|---|
| `func load-user()` | punctuation not allowed | `func loadUser()` |
| `func 2FactorAuth()` | cannot start with digit | `func twoFactorAuth()` |
| `max_value := 10` | non-standard casing | `maxValue := 10` |
| `type HttpClient struct{}` | inconsistent acronym casing | `type HTTPClient struct{}` |
| `func GetSessionId()` | `ID` must be all caps | `func GetSessionID()` |
| `func clear()` | clashes with builtin | `func clearQueue()` |
| `resultSlice := []int{}` | type in name | `results := []int{}` |
| `var log = newLogger()` | clashes with stdlib package | `var logger = newLogger()` |

## Exported vs unexported

- Capital first letter = **exported** (visible outside the package). Casing of the first letter
  changes behavior — never capitalize just because it looks nice.
- **Default to unexported.** Export only when needed. Less exporting = easier refactoring
  ("write shy code").
- `main` package identifiers should normally all be unexported.
- Exception: struct fields consumed by reflection-based packages (`encoding/json`,
  `encoding/gob`, `sqlx`) must be exported to be visible.

## Identifier length

The further an identifier is used from its declaration, the more descriptive it must be.

- Small scope (short `for`/`range` block, tiny function): short or single-letter names are fine
  (`for _, p := range people`).
- Used across a whole function or beyond: descriptive names (`count`, `sum` — not `c`, `s`,
  but also not needlessly verbose `peopleCount`, `agesSum`).

## Package names

- Lowercase ASCII letters and numbers only.
- Short, easy to type, reflects contents — often one-word nouns: `orders`, `customer`, `slug`.
- Multiple words: concatenate, no separator — `ordermanager`, never `orderManager` or
  `order_manager`.
- Abbreviations are OK if the name gets long (`strconv`, `expvar`).
- Avoid names of commonly-used stdlib packages.
- Never prefix with `.` or `_` — such packages are invisible to Go tooling.
- Never use `vendor`, `testdata`, or `internal` as package names — special meaning in Go.
- **No catch-all names** like `common`, `util`, `helpers`, `types`, `interfaces` — unclear
  scope, large blast radius, import-cycle risk. Split into focused packages instead
  (`validation`, `formatting`).

| Bad | Reason | Better |
|---|---|---|
| `package order_manager` | separators/casing | `package ordermanager` |
| `package stuff` | too vague | `package orders` |
| `package url` | stdlib clash | `package links` |
| `package internal` | special directory name | `package internalauth` |
| `package utils` | catch-all | `package validation` |

## File names

- Ideally one word, all lowercase, summarizing contents: `cookie.go`, `server.go`, `status.go`.
- Multi-word: no strong convention — concatenate (`routingindex.go`) or underscore
  (`routing_index.go`); pick one and stay consistent. Prefer concatenation and reserve `_` for
  special suffixes.
- **Special meanings — avoid unless you want the behavior:**
  - `.` or `_` prefix → ignored by Go tooling.
  - `_test.go` → only built by `go test`.
  - `_linux.go`, `_windows.go`, `_darwin.go`, ... → OS-specific builds.
  - `_amd64.go`, `_arm64.go`, `_wasm.go`, ... → architecture-specific builds.

## Avoiding chatter

Don't repeat the package name inside exported identifiers — it stutters at the call site.

| Bad | Better |
|---|---|
| `customer.NewCustomer()` | `customer.New()` |
| `customer.CustomerOrders()` | `customer.Orders()` |
| `customer.CustomerAddress` | `customer.Address` |

- Acceptable exception: an exported type sharing the package name when renaming would reduce
  clarity — `time.Time`, `context.Context`, `regexp.Regexp`.
- Same for methods: on a `Token` type, prefer `Validate()` over `ValidateToken()`, `IsExpired()`
  over `IsTokenExpired()`.

## Method receivers

- Short: 1–3 characters, usually an abbreviation of the type: `c` or `cus` for `Customer`,
  `hs` for `HighScore`.
- Never `this`, `self`, or `me`.
- **Consistent across all methods of the same type** — don't mix `c` and `cus`.

```go
// Good
func (o *Order) Validate() bool { return o.Items > 0 }

// Bad: too long
func (order *Order) Validate() bool { return order.Items > 0 }

// Bad: generic name
func (self *Order) Validate() bool { return self.Items > 0 }
```

## Getters and setters

- Prefer direct field access; only add getter/setter methods for **unexported** fields that
  must be reachable from outside the package.
- Setter gets the `Set` prefix; getter does **not** get a `Get` prefix.

```go
type Customer struct {
	address string
}

func (c *Customer) Address() string        { return c.address }
func (c *Customer) SetAddress(addr string) { c.address = addr }
```

## Interfaces

- Single-method interfaces: method name + `-er` (or similar): `Speaker` for `Speak()`,
  `Authorizer` for `Authorize()`, like stdlib `io.Reader`, `io.Writer`, `fmt.Stringer`.
- No type-in-name: never `UserInterface` or `OrderInterface`.

## Breaking conventions

Rarely acceptable (e.g. mirroring an external system's exact identifiers in a private codebase
when it clarifies intent). Default: follow the conventions — they make code predictable,
consistent, and reduce bugs.
