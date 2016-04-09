package vollocal_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-golang/clock/fakeclock"
	"github.com/pivotal-golang/lager/lagertest"

	"github.com/cloudfoundry-incubator/volman/vollocal"
	volmanfakes "github.com/cloudfoundry-incubator/volman/volmanfakes"
)

var _ = Describe("Registry", func() {
	var (
		logger *lagertest.TestLogger

		scanInterval time.Duration

		fakeClock         *fakeclock.FakeClock
		fakeDriverFactory *volmanfakes.FakeDriverFactory

		registry *vollocal.DriversRegistry
	)

	BeforeEach(func() {

		logger = lagertest.NewTestLogger("RegistryTest")
		fakeClock = fakeclock.NewFakeClock(time.Unix(123, 456))
		fakeDriverFactory = new(volmanfakes.FakeDriverFactory)

		scanInterval = 10 * time.Second

		registry = vollocal.NewRegistry(logger, fakeDriverFactory, scanInterval, fakeClock)
	})

	Describe("#SetDrivers", func() {
		Context("when there are no drivers", func() {
			It("should have no drivers in registry map", func() {
				Expect(len(registry.DriversMap)).To(Equal(0))

				Expect(fakeDriverFactory.DiscoverCallCount()).To(Equal(0))
				Expect(fakeDriverFactory.DriverCallCount()).To(Equal(0))
			})

		})
	})

	Describe("#ifrit.RunFunc", func() {
		// Context("When RegistryRunner is run", func() {
		// 	BeforeEach(func() {
		// 		// driverName := "fakedriver"
		// 		// err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, driverName, "http://0.0.0.0:8080")
		// 		// Expect(err).NotTo(HaveOccurred())
		// 		// fakedriverProcess = ginkgomon.Invoke(fakedriverRunner)
		// 		// registryMap, _ = vollocal.SetDrivers(defaultPluginsDirectory)
		// 	})
		// 	It("should set up the registry with existing drivers", func() {
		// 		Expect(len(registryMap)).To(Equal(1))

		// 	})
		// })
	})
})
