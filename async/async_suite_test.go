package async_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAsync(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Async Suite")
}
