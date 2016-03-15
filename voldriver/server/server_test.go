package server_test

import (
	"net"

	"github.com/cloudfoundry-incubator/volman/voldriver/fakes"
	"github.com/cloudfoundry-incubator/volman/voldriver/server"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/tedsuo/ifrit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Volman Driver Server", func() {
	var (
		logger *lagertest.TestLogger
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
	})

	Context("when passed a tcp address", func() {
		It("listens on the given address", func() {
			serverRunner := server.New(logger, new(fakes.FakeBackend), "tcp", ":60555")
			serverProcess := ifrit.Invoke(serverRunner)
			Eventually(serverProcess.Ready()).Should(BeClosed())
			_, err := net.Dial("tcp", ":60555")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
