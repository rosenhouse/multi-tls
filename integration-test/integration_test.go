package integration_test

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Integration tests", func() {
	var routerSession *gexec.Session
	var backendSessions []*gexec.Session

	const routerListenPort = 2000
	const backendStartPort = 10000
	const numBackends = 5

	BeforeEach(func() {
		var err error
		routerCmd := exec.Command(paths.Router,
			"-listenPort", fmt.Sprintf("%d", routerListenPort),
		)
		routerSession, err = gexec.Start(routerCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		backendSessions = make([]*gexec.Session, numBackends)
		for i := 0; i < numBackends; i++ {
			backendCmd := exec.Command(paths.Backend,
				"-listenPort", fmt.Sprintf("%d", backendStartPort+i))
			backendSessions[i], err = gexec.Start(backendCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	AfterEach(func() {
		routerSession.Interrupt()

		for i := 0; i < numBackends; i++ {
			backendSessions[i].Interrupt()
		}
		Eventually(routerSession).Should(gexec.Exit())
		for i := 0; i < numBackends; i++ {
			Eventually(backendSessions[i]).Should(gexec.Exit())
		}
	})

	Describe("server lifecycle", func() {
		It("stays alive", func() {
			Consistently(routerSession).ShouldNot(gexec.Exit())

			for i := 0; i < numBackends; i++ {
				Consistently(backendSessions[i]).ShouldNot(gexec.Exit())
			}
		})
	})
})
