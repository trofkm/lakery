package lakery_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLakery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lakery Suite")
}
