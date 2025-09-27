# Lakery Build-time Validation Approaches

This document explains the different ways lakery can automatically validate struct tags at build time **without requiring users to manually run `go generate`**.

## 🚀 Approach 1: Automatic Validation on Import (Recommended)

**File**: `validate.go`

**How it works**: 
- Validation runs automatically in `init()` when lakery is imported
- Scans the entire project for files using lakery tags
- Fails compilation if invalid tags are found

**Usage**:
```go
import _ "github.com/trofkm/lakery"  // Triggers automatic validation
```

**Pros**:
- ✅ Zero configuration required
- ✅ Works out of the box
- ✅ Catches errors immediately
- ✅ Project-wide validation

**Cons**:
- ⚠️ Adds startup time during development
- ⚠️ May need to be disabled in production

**Control**:
```go
// Disable programmatically
lakery.AutoValidateOnImport = false

// Or via environment variable
export LAKERY_SKIP_VALIDATION=1
```

## 🔒 Approach 2: Compile-time Enforcement with Build Tags

**File**: `compile_check.go`

**How it works**:
- Only compiled when `lakery_validate` build tag is present
- Panics during compilation if validation fails
- Forces build to fail with invalid tags

**Usage**:
```bash
go build -tags lakery_validate ./...
```

**Pros**:
- ✅ Explicit opt-in behavior
- ✅ Integrates well with CI/CD
- ✅ No runtime overhead

**Cons**:
- ⚠️ Requires explicit build tag
- ⚠️ Less discoverable

## 📋 Approach 3: Manual Validation Tools

**Files**: `cmd/lakery-validate/main.go`, `go:generate` directive

**How it works**:
- Standalone validation tool
- Can be run via `go:generate`, Makefile, or CI/CD
- Provides detailed error reporting

**Usage**:
```bash
# Via go:generate
go generate

# Direct command
go run github.com/trofkm/lakery/cmd/lakery-validate -package .

# Via Makefile
make validate
```

**Pros**:
- ✅ Explicit control
- ✅ Detailed error messages
- ✅ Can be integrated anywhere
- ✅ No runtime impact

**Cons**:
- ⚠️ Requires manual invocation
- ⚠️ Can be forgotten

## 🛠️ Integration Examples

### CI/CD (GitHub Actions)
```yaml
- name: Validate Lakery tags
  run: go run github.com/trofkm/lakery/cmd/lakery-validate -package .
```

### Makefile
```makefile
build: validate
	go build ./...

validate:
	go run github.com/trofkm/lakery/cmd/lakery-validate -package .
```

### Pre-commit Hook
```bash
#!/bin/sh
go run github.com/trofkm/lakery/cmd/lakery-validate -package .
```

## 🎯 Validation Capabilities

All approaches catch the same types of errors:

- **Unknown validators**: `lakery:"nonexistent_validator"`
- **Syntax errors**: `lakery:"each={min=1,max=10"` (unclosed brace)
- **Wrong syntax**: `lakery:"each:{}"` (colon instead of equals)
- **Parameter validation**: `lakery:"min=not_a_number"`
- **Parameter misuse**: `lakery:"required=true"`

## 🏆 Recommendation

**For most users**: Use **Approach 1 (Automatic Validation)** as it provides the best developer experience with zero configuration.

**For CI/CD pipelines**: Add **Approach 3 (Manual Tools)** as an explicit validation step.

**For strict environments**: Use **Approach 2 (Build Tags)** with `-tags lakery_validate` in your build process.

**Best practice**: Combine approaches for comprehensive coverage:
- Automatic validation for development
- Explicit validation in CI/CD
- Build tag validation for production builds
