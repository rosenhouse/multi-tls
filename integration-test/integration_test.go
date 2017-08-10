package integration_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Integration tests", func() {
	var routerSession *gexec.Session
	var backendSessions []*gexec.Session

	const routerPort = 2000
	const backendStartPort = 10000
	const numBackends = 5

	BeforeEach(func() {
		var err error
		routerCmd := exec.Command(paths.Router,
			"-listenPort", fmt.Sprintf("%d", routerPort),
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

		for i := 0; i < numBackends; i++ {
			Eventually(backendSessions[i].Out).Should(gbytes.Say("will listen on port"))
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

	Describe("router", func() {
		It("forwards requests to the correct backend", func() {
			routerURL := fmt.Sprintf("http://127.0.0.1:%d", routerPort)
			req, _ := http.NewRequest("GET", routerURL, nil)
			req.Host = "backend-10003"

			resp, err := http.DefaultClient.Do(req)
			Expect(err).NotTo(HaveOccurred())

			respBody, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			Expect(respBody).To(ContainSubstring("Hello from backend listening on 10003"))
		})

	})
})
