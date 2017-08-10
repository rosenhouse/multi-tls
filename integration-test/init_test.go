package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestIntegrationTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Test Suite")
}

// path to binary executables used for testing
var paths struct {
	Backend string
	Router  string
}

// runs before any tests
var _ = BeforeSuite(func() {
	// build binary executables used in testing
	var err error
	paths.Backend, err = gexec.Build("github.com/rosenhouse/multi-tls/backend")
	Expect(err).NotTo(HaveOccurred())

	paths.Router, err = gexec.Build("github.com/rosenhouse/multi-tls/router")
	Expect(err).NotTo(HaveOccurred())
})
