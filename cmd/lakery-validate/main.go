package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// Built-in validator tags that are always available
var builtinValidators = map[string]bool{
	"min":      true,
	"max":      true,
	"required": true,
	"each":     true, // special tag
	"dive":     true, // special tag
}

// TagError represents a validation error in a lakery tag
type TagError struct {
	File    string
	Line    int
	Column  int
	Field   string
	Tag     string
	Message string
}

func (e TagError) Error() string {
	return fmt.Sprintf("%s:%d:%d: lakery tag error in field %q: %s (tag: %q)",
		e.File, e.Line, e.Column, e.Field, e.Message, e.Tag)
}

// TagValidator validates lakery struct tags at build time
type TagValidator struct {
	fset             *token.FileSet
	customValidators map[string]bool
	errors           []TagError
}

func NewTagValidator() *TagValidator {
	return &TagValidator{
		fset:             token.NewFileSet(),
		customValidators: make(map[string]bool),
		errors:           []TagError{},
	}
}

func (tv *TagValidator) addError(pos token.Pos, field, tag, message string) {
	position := tv.fset.Position(pos)
	tv.errors = append(tv.errors, TagError{
		File:    position.Filename,
		Line:    position.Line,
		Column:  position.Column,
		Field:   field,
		Tag:     tag,
		Message: message,
	})
}

// ValidatePackage validates all lakery tags in the given package directory
func (tv *TagValidator) ValidatePackage(pkgDir string) error {
	// Parse all Go files in the package, excluding test files
	pkgs, err := parser.ParseDir(tv.fset, pkgDir, func(info os.FileInfo) bool {
		// Skip test files as they often contain intentionally invalid code
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse package: %w", err)
	}

	// First pass: collect custom validator registrations
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			tv.findCustomValidators(file)
		}
	}

	// Second pass: validate lakery tags
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			tv.validateFile(file)
		}
	}

	if len(tv.errors) > 0 {
		for _, err := range tv.errors {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		return fmt.Errorf("found %d lakery tag validation errors", len(tv.errors))
	}

	return nil
}

// findCustomValidators searches for RegisterTag calls to discover custom validators
func (tv *TagValidator) findCustomValidators(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		// Look for method calls like validator.RegisterTag("name", func)
		if call, ok := n.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "RegisterTag" && len(call.Args) >= 1 {
					// Extract the validator name from the first argument
					if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
						name, err := strconv.Unquote(lit.Value)
						if err == nil {
							tv.customValidators[name] = true
						}
					}
				}
			}
		}
		return true
	})
}

// validateFile validates all lakery tags in a single Go file
func (tv *TagValidator) validateFile(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		// Look for struct type definitions
		if ts, ok := n.(*ast.TypeSpec); ok {
			if st, ok := ts.Type.(*ast.StructType); ok {
				tv.validateStruct(ts.Name.Name, st)
			}
		}
		return true
	})
}

// validateStruct validates lakery tags in a struct definition
func (tv *TagValidator) validateStruct(structName string, st *ast.StructType) {
	for _, field := range st.Fields.List {
		if field.Tag != nil {
			tagValue := field.Tag.Value
			// Remove quotes from tag string
			if len(tagValue) >= 2 && tagValue[0] == '`' && tagValue[len(tagValue)-1] == '`' {
				tagValue = tagValue[1 : len(tagValue)-1]
			}

			// Parse struct tag
			tag := reflect.StructTag(tagValue)
			lakeryTag := tag.Get("lakery")

			if lakeryTag != "" {
				fieldName := ""
				if len(field.Names) > 0 {
					fieldName = field.Names[0].Name
				} else {
					fieldName = "<embedded>"
				}

				tv.validateLakeryTag(field.Pos(), fieldName, lakeryTag)
			}
		}
	}
}

// validateLakeryTag validates the syntax and semantics of a lakery tag
func (tv *TagValidator) validateLakeryTag(pos token.Pos, fieldName, tag string) {
	// First, validate the overall syntax by parsing top-level comma-separated parts
	parts, err := tv.splitTopLevelByComma(tag)
	if err != nil {
		tv.addError(pos, fieldName, tag, err.Error())
		return
	}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		tv.validateTagPart(pos, fieldName, tag, part)
	}
}

// validateTagPart validates a single tag part (like "min=5" or "each={min=1,max=10}")
func (tv *TagValidator) validateTagPart(pos token.Pos, fieldName, fullTag, part string) {
	// Split by = to get key and value
	kv := strings.SplitN(part, "=", 2)
	key := strings.TrimSpace(kv[0])

	if key == "" {
		tv.addError(pos, fieldName, fullTag, "empty validator name")
		return
	}

	// Special handling for 'each' tag
	if key == "each" {
		if len(kv) == 2 {
			value := strings.TrimSpace(kv[1])
			tv.validateEachTag(pos, fieldName, fullTag, value)
		}
		return
	}

	// Check if the validator exists
	if !tv.validatorExists(key) {
		// Special case: check if user meant "each=" instead of "each:"
		if strings.HasPrefix(key, "each:") {
			tv.addError(pos, fieldName, fullTag, fmt.Sprintf("invalid syntax %q - did you mean \"each=%s\"?", key, key[5:]))
		} else {
			tv.addError(pos, fieldName, fullTag, fmt.Sprintf("unknown validator %q", key))
		}
		return
	}

	// Validate parameter format for specific validators
	if len(kv) == 2 {
		value := strings.TrimSpace(kv[1])
		tv.validateValidatorParam(pos, fieldName, fullTag, key, value)
	}
}

// validateEachTag validates the content inside each={...}
func (tv *TagValidator) validateEachTag(pos token.Pos, fieldName, fullTag, value string) {
	// Remove surrounding braces if present
	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		value = strings.TrimSpace(value[1 : len(value)-1])
	}

	if value == "" {
		return // Empty each tag is valid
	}

	// Parse the inner validators
	innerParts, err := tv.splitTopLevelByComma(value)
	if err != nil {
		tv.addError(pos, fieldName, fullTag, fmt.Sprintf("invalid syntax in each tag: %s", err.Error()))
		return
	}

	for _, innerPart := range innerParts {
		innerPart = strings.TrimSpace(innerPart)
		if innerPart == "" {
			continue
		}
		tv.validateTagPart(pos, fieldName, fullTag, innerPart)
	}
}

// validateValidatorParam validates parameters for specific validators
func (tv *TagValidator) validateValidatorParam(pos token.Pos, fieldName, fullTag, validator, param string) {
	switch validator {
	case "min", "max":
		if _, err := strconv.Atoi(param); err != nil {
			tv.addError(pos, fieldName, fullTag, fmt.Sprintf("%s expects integer parameter, got %q", validator, param))
		}
	case "required":
		if param != "" {
			tv.addError(pos, fieldName, fullTag, "required validator does not accept parameters")
		}
	}
}

// validatorExists checks if a validator is available (builtin or custom)
func (tv *TagValidator) validatorExists(name string) bool {
	return builtinValidators[name] || tv.customValidators[name]
}

// splitTopLevelByComma splits a string by commas, ignoring commas inside curly braces
// This duplicates the logic from the main lakery package to ensure consistency
func (tv *TagValidator) splitTopLevelByComma(s string) ([]string, error) {
	var parts []string
	depth := 0
	last := 0

	for i, r := range s {
		switch r {
		case '{':
			depth++
		case '}':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, s[last:i])
				last = i + 1
			}
		}
	}
	parts = append(parts, s[last:])

	if depth < 0 {
		return nil, fmt.Errorf("unopened braces in %q", s)
	} else if depth > 0 {
		return nil, fmt.Errorf("unclosed braces in %q", s)
	}

	return parts, nil
}

func main() {
	var (
		pkgDir = flag.String("package", ".", "Package directory to validate")
		help   = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nValidates lakery struct tags at build time.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Convert relative path to absolute
	absPath, err := filepath.Abs(*pkgDir)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	validator := NewTagValidator()
	if err := validator.ValidatePackage(absPath); err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	fmt.Println("All lakery tags are valid")
}
