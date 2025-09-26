package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/trofkm/lakery"
)

// Demonstrates lakery:"each={min=0,max=23,credential}" usage on a slice of strings.

const (
	minTag        = "min"
	maxTag        = "max"
	credentialTag = "credential"
)

// minValidator checks that string length is >= provided value
func minValidator(lk *lakery.Value) error {
	val := lk.String()
	minLen, err := strconv.Atoi(lk.Param())
	if err != nil {
		return err
	}
	if len(val) < minLen {
		return fmt.Errorf("min len should be at least %d", minLen)
	}
	return nil
}

// maxValidator checks that string length is <= provided value
func maxValidator(lk *lakery.Value) error {
	val := lk.String()
	maxLen, err := strconv.Atoi(lk.Param())
	if err != nil {
		return err
	}
	if len(val) > maxLen {
		return fmt.Errorf("max len should be less equal then %d", maxLen)
	}
	return nil
}

// credentialValidator enforces a basic credential format: letters, digits, underscore, dash only
func credentialValidator(lk *lakery.Value) error {
	val := lk.String()
	// simple pattern: 1-64 of [a-zA-Z0-9_-]
	var re = regexp.MustCompile(`^[a-zA-Z0-9_-]*$`)
	if !re.MatchString(val) {
		return fmt.Errorf("should contain only letters, numbers, '_' or '-'")
	}
	return nil
}

type Payload struct {
	Credentials []string `lakery:"each={min=0,max=23,credential}"`
}

func main() {
	val := lakery.NewValidator()
	val.RegisterTag(minTag, minValidator)
	val.RegisterTag(maxTag, maxValidator)
	val.RegisterTag(credentialTag, credentialValidator)

	good := Payload{Credentials: []string{"john_doe", "user-123", "a"}}
	if err := val.Validate(good); err != nil {
		fmt.Println("good payload validation error:", err)
	} else {
		fmt.Println("good payload passed validation")
	}

	bad := Payload{Credentials: []string{"", "this_credential_is_way_too_long_beyond_twenty_three", "no spaces"}}
	if err := val.Validate(bad); err != nil {
		fmt.Println("bad payload validation error:", err)
	} else {
		fmt.Println("bad payload passed validation (unexpected)")
	}
}
