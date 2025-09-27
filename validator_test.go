package lakery_test

import (
	"errors"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/trofkm/lakery"
)

var _ = Describe("Validator", func() {
	Context("construction", func() {
		It("auto-registers builtins", func() {
			v := lakery.NewValidator()
			validators := v.ListValidators()
			Expect(validators).To(ContainElements("min", "max", "required"))
		})
	})

	Context("min and max", func() {
		type S struct {
			Name string `lakery:"min=2,max=4"`
		}

		It("passes when length within bounds", func() {
			v := lakery.NewValidator()
			s := S{Name: "john"}
			Expect(v.Validate(s)).To(Succeed())
		})

		It("fails when length below min", func() {
			v := lakery.NewValidator()
			s := S{Name: "J"}
			err := v.Validate(s)
			Expect(err).To(HaveOccurred())
		})

		It("fails when length above max", func() {
			v := lakery.NewValidator()
			s := S{Name: "johnny"}
			err := v.Validate(s)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("required", func() {
		type S struct {
			Name string `lakery:"required"`
			Ptr  *int   `lakery:"required"`
		}
		It("fails on zero values", func() {
			v := lakery.NewValidator()
			s := S{}
			Expect(v.Validate(s)).To(HaveOccurred())
		})
		It("passes on non-zero", func() {
			v := lakery.NewValidator()
			x := 10
			s := S{Name: "ok", Ptr: &x}
			Expect(v.Validate(s)).To(Succeed())
		})
	})
	Context("each with invalid", func() {
		It("empty parameters", func() {
			type S struct {
				Creds []string `lakery:"each={}"`
			}
			v := lakery.NewValidator()
			s := S{Creds: []string{"a", "bb", "ccc"}}
			Expect(v.Validate(s)).To(Succeed())
		})
		It("unclosed parantheses with no params", func() {
			type S struct {
				Creds []string `lakery:"each={"`
			}
			v := lakery.NewValidator()
			s := S{Creds: []string{"a", "bb", "ccc"}}
			Expect(v.Validate(s)).To(MatchError(ContainSubstring("unclosed braces")))
		})
		It("unopened parantheses with no params", func() {
			type S struct {
				Creds []string `lakery:"each=}"`
			}
			v := lakery.NewValidator()
			s := S{Creds: []string{"a", "bb", "ccc"}}
			Expect(v.Validate(s)).To(MatchError(ContainSubstring("unopened braces")))
		})
		It("unclosed parantheses with params", func() {
			type S struct {
				Creds []string `lakery:"each={min=100"`
			}
			v := lakery.NewValidator()
			s := S{Creds: []string{"a", "bb", "ccc"}}
			Expect(v.Validate(s)).To(MatchError(ContainSubstring("unclosed braces")))
		})
		It("invalid colon syntax should be ignored (unknown validator)", func() {
			type S struct {
				Creds []string `lakery:"each:{}"`
			}
			v := lakery.NewValidator()
			s := S{Creds: []string{"a", "bb", "ccc"}}
			// Runtime validator ignores unknown validators, but build-time should catch this
			Expect(v.Validate(s)).To(Succeed())
		})
	})

	Context("each for string slice", func() {
		type S struct {
			Creds []string `lakery:"each={min=1,max=5}"`
		}
		It("passes when all elements are valid", func() {
			v := lakery.NewValidator()
			s := S{Creds: []string{"a", "bb", "ccc"}}
			Expect(v.Validate(s)).To(Succeed())
		})
		It("fails when any element invalid", func() {
			v := lakery.NewValidator()
			s := S{Creds: []string{"", "bb"}}
			Expect(v.Validate(s)).To(HaveOccurred())
		})
		It("errors when used on non-slice", func() {
			type T struct {
				Name string `lakery:"each={min=1}"`
			}
			v := lakery.NewValidator()
			s := T{Name: "aa"}
			err := v.Validate(s)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("custom error formatter", func() {
		It("wraps underlying error", func() {
			old := lakery.CurrentErrorFormatFunc
			defer func() { lakery.CurrentErrorFormatFunc = old }()
			lakery.CurrentErrorFormatFunc = func(sf reflect.StructField, rv reflect.Value, err error) error {
				return errors.Join(errors.New("wrapped"), err)
			}
			type S struct {
				Name string `lakery:"min=3"`
			}
			v := lakery.NewValidator()
			s := S{Name: "aa"}
			err := v.Validate(s)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("wrapped"))
		})
	})
})
