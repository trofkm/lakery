# ✅ Lakery — Tiny Tag-Based Validator for Go

Lakery is a tiny, dependency-free, tag-based validation library for Go. Define validation rules right in struct tags, register your own validators, and validate with a single call.

## 🌟 Features

- **Zero dependencies** — pure Go
- **Built-in tags** — `min`, `max`, `required`
- **Collection rules** — `each={...}` applies validators to every element of a slice/array
- **Pluggable validators** — register custom tags easily
- **Custom error formatting** — control how validation errors are presented

## 🚀 Quick Start

```go
package main

import (
	"fmt"
	"github.com/trofkm/lakery"
)

type User struct {
	Name      string   `lakery:"required,min=2,max=32"`
	Nicknames []string `lakery:"each={min=1,max=10}"`
}

func main() {
	v := lakery.NewValidator() // built-ins are auto-registered

	u := User{
		Name:      "John",
		Nicknames: []string{"jj", "doe"},
	}
	if err := v.Validate(u); err != nil {
		fmt.Println("validation error:", err)
	}
}
```

## 🧩 Tags and Syntax

- **Simple tags**: `lakery:"required"`, `lakery:"min=1,max=10"`
- **Each for collections**: `lakery:"each={min=0,max=23,credential}"`
	- Curly braces contain a comma-separated list of validators applied to every element

### Built-in Tags

- `required` — value must be non-zero (non-empty string, non-nil pointer/slice/map, non-zero numbers, etc.)
- `min` — for strings/slices/arrays/maps checks length ≥ N; for numbers checks value ≥ N
- `max` — for strings/slices/arrays/maps checks length ≤ N; for numbers checks value ≤ N

### Custom Tags (Example)

```go
const credential = "credential"

func credentialValidator(lk *lakery.Value) error {
	val := lk.String()
	var re = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)
	if !re.MatchString(val) {
		return fmt.Errorf("should contain only letters, numbers, '_' or '-'")
	}
	return nil
}

func main() {
	v := lakery.NewValidator()
	v.RegisterTag(credential, credentialValidator)

	type Payload struct {
		Credentials []string `lakery:"each={min=1,max=23,credential}"`
	}

	p := Payload{Credentials: []string{"john_doe", "user-123"}}
	_ = v.Validate(p)
}
```

## 🧪 Examples

- Minimal custom tag: `_example/simple/main.go`
- Custom validator (credential): `_example/custom_validator/main.go`
- Error formatting (i18n): `_example/error_fmt/main.go`
- Collections with `each={...}`: `_example/each/main.go`

Run any example, for example:

```bash
cd _example/each && go run .
```

## 🧰 API Reference

```go
// Create a validator (built-ins auto-registered)
func NewValidator() *Validator

// Register custom tag validators
type TagValidationFunc = func(*Value) error
func (v *Validator) RegisterTag(tag string, fn TagValidationFunc)

// Inspect registered tags
func (v *Validator) ListValidators() []string

// Validate a struct value
func (v *Validator) Validate(s any) error

// Customize error formatting
type ErrorFormatFunc = func(fieldType reflect.StructField, fieldValue reflect.Value, err error) error
var CurrentErrorFormatFunc ErrorFormatFunc
```

### Value Helpers (for validator authors)

```go
type Value struct {
	// underlying reflect.Value and tag param are encapsulated
}

func (v *Value) String() string   // returns underlying string value
func (v *Value) Interface() any   // returns underlying interface value
func (v *Value) Param() string    // returns tag parameter (e.g., "10" for min=10)
```

## 🧭 Behavior Notes

- Lakery validates only structs passed to `Validate`.
- The `each={...}` tag is special-cased and applies included validators to every element of a slice/array.
- Tag parsing supports comma-separated lists and ignores commas inside `{ ... }` blocks.
- Built-ins are registered automatically in `NewValidator`.

## 🧪 Tests

This project uses Ginkgo and Gomega.

```bash
go test ./...
```

The suite covers:
- Built-in registration (`min`, `max`, `required`)
- `min`/`max` for string length
- `required` for strings and pointers
- `each={...}` validation on string slices
- Custom error formatting

## 📦 Installation

```bash
go get github.com/trofkm/lakery
```

## 🛣️ Roadmap

- [x] Simple tag validation (e.g., credential, email via custom tags)
- [x] Validation expressions (e.g., `min=0,max=255`)
- [x] Collection validation (`each={...}`)
- [ ] Dive into nested structs with `dive`
- [ ] More tests

## 📄 License

MIT License — see `LICENSE` for details.

---

"Small, readable, and extensible validation for Go"
