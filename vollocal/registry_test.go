package vollocal_test

import (
	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/cloudfoundry-incubator/volman/vollocal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/tedsuo/ifrit/ginkgomon"
)

var _ = FDescribe("Registry", func() {
	var (
		logger      = lagertest.NewTestLogger("RegistryTest")
		registryMap map[string]string
	)

	Context("When there are no drivers on Path", func() {
		BeforeEach(func() {
			registryMap, _ = vollocal.SetDrivers(defaultPluginsDirectory)
		})
		It("should set up the registry with no drivers", func() {
			Expect(len(registryMap)).To(Equal(0))
		})
	})
	Context("When there are drivers on Path", func() {
		BeforeEach(func() {
			driverName := "fakedriver"
			err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, driverName, "http://0.0.0.0:8080")
			Expect(err).NotTo(HaveOccurred())
			fakedriverProcess = ginkgomon.Invoke(fakedriverRunner)
			registryMap, _ = vollocal.SetDrivers(defaultPluginsDirectory)
		})
		It("should set up the registry with existing drivers", func() {
			Expect(len(registryMap)).To(Equal(1))

		})
	})
})
