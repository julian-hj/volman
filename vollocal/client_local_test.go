package vollocal_test

import (
	"bytes"
	"io"
	"time"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/cloudfoundry-incubator/volman/vollocal"
	"github.com/cloudfoundry-incubator/volman/volmanfakes"

	"github.com/pivotal-golang/clock/fakeclock"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
)

var _ = FDescribe("Volman", func() {
	var (
		logger = lagertest.NewTestLogger("client-test")

		fakeDriverFactory *volmanfakes.FakeDriverFactory
		fakeDriver        *volmanfakes.FakeDriver
		fakeClock         *fakeclock.FakeClock

		scanInterval time.Duration

		registry                *vollocal.DriversRegistry
		validDriverInfoResponse io.ReadCloser
		runner                  ifrit.Runner
	)

	BeforeEach(func() {
		fakeDriverFactory = new(volmanfakes.FakeDriverFactory)
		fakeClock = fakeclock.NewFakeClock(time.Unix(123, 456))

		scanInterval = 10 * time.Second

		registry = vollocal.NewRegistry(logger, fakeDriverFactory, scanInterval, fakeClock)

		validDriverInfoResponse = stringCloser{bytes.NewBufferString("{\"Name\":\"fakedriver\",\"Path\":\"somePath\"}")}
	})

	Describe("ListDrivers", func() {
		Context("", func() {
			BeforeEach(func() {
				client, runner = vollocal.NewLocalClient(fakeDriverFactory, *registry)
				Expect(runner).ToNot(BeNil())
			})

			It("should report empty list of drivers", func() {
				drivers, err := client.ListDrivers(logger)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(drivers.Drivers)).To(Equal(0))
			})

			Context("has no drivers in location", func() {
				BeforeEach(func() {
					fakeDriverFactory = new(volmanfakes.FakeDriverFactory)
					fakeDriverFactory.DriversDirReturns("")

					registry = vollocal.NewRegistry(logger, fakeDriverFactory, scanInterval, fakeClock)
				})

				It("should report empty list of drivers", func() {
					drivers, err := client.ListDrivers(logger)
					Expect(err).NotTo(HaveOccurred())
					Expect(len(drivers.Drivers)).To(Equal(0))
				})
			})
		})

		Context("has driver in location", func() {
			BeforeEach(func() {
				err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, "fakedriver", "http://0.0.0.0:8080")
				Expect(err).NotTo(HaveOccurred())

				client, _ = vollocal.NewLocalClient(fakeDriverFactory, *registry)
			})

			It("should report empty list of drivers", func() {
				drivers, err := client.ListDrivers(logger)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(drivers.Drivers)).To(Equal(0))
			})

			FContext("after running drivers discovery", func() {
				BeforeEach(func() {
					fakeDriverFactory.DiscoverReturns(map[string]string{"fake-driver": "fake-driver"}, nil)
					registry.SetDrivers(logger)
					//registry.DriversMap = map[string]voldriver.Driver{"fake-driver": fakeDriver}
				})

				It("should report list of drivers", func() {
					drivers, err := client.ListDrivers(logger)
					Expect(err).NotTo(HaveOccurred())
					Expect(len(drivers.Drivers)).ToNot(Equal(0))
				})

				It("should report at least fakedriver", func() {
					drivers, err := client.ListDrivers(logger)
					Expect(err).NotTo(HaveOccurred())
					Expect(len(drivers.Drivers)).ToNot(Equal(0))
					Expect(drivers.Drivers[0].Name).To(Equal("fake-driver"))
				})
			})
		})

		Context("discovery fails", func() {
			It("it should fail", func() {
				//TODO
			})
		})
	})

	Describe("Mount and Unmount", func() {
		Context("when given valid driver", func() {
			BeforeEach(func() {
				fakedriverProcess = ginkgomon.Invoke(fakedriverRunner)
				fakeDriverFactory = new(volmanfakes.FakeDriverFactory)
				fakeDriver = new(volmanfakes.FakeDriver)
				fakeDriverFactory.DriverReturns(fakeDriver, nil)

				client, runner = vollocal.NewLocalClient(fakeDriverFactory, *registry)
				Expect(runner).ToNot(BeNil())
			})

			It("should be able to mount", func() {
				volumeId := "fake-volume"

				mountPath, err := client.Mount(logger, "fakedriver", volumeId, map[string]interface{}{"volume_id": volumeId})
				Expect(err).NotTo(HaveOccurred())
				Expect(mountPath).NotTo(Equal(""))
			})

			It("should not be able to mount if mount fails", func() {
				mountResponse := voldriver.MountResponse{Err: "an error"}
				fakeDriver.MountReturns(mountResponse)

				volumeId := "fake-volume"
				_, err := client.Mount(logger, "fakedriver", volumeId, map[string]interface{}{"volume_id": volumeId})
				Expect(err).To(HaveOccurred())
			})

			It("should be able to unmount", func() {
				volumeId := "fake-volume"

				err := client.Unmount(logger, "fakedriver", volumeId)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeDriver.UnmountCallCount()).To(Equal(1))
				Expect(fakeDriver.RemoveCallCount()).To(Equal(1))
			})

			It("should not be able to unmount when driver unmount fails", func() {
				fakeDriver.UnmountReturns(voldriver.ErrorResponse{Err: "unmount failure"})
				volumeId := "fake-volume"

				err := client.Unmount(logger, "fakedriver", volumeId)
				Expect(err).To(HaveOccurred())
			})

		})

		Context("when given invalid driver", func() {
			BeforeEach(func() {
				fakedriverProcess = ginkgomon.Invoke(fakedriverRunner)
				fakeDriverFactory = new(volmanfakes.FakeDriverFactory)
				fakeDriver = new(volmanfakes.FakeDriver)

				fakeDriverFactory.DriverReturns(fakeDriver, nil)
				fakeDriverFactory.DriverReturns(nil, fmt.Errorf("driver not found"))

				client, runner = vollocal.NewLocalClient(fakeDriverFactory, *registry)
				Expect(runner).ToNot(BeNil())
			})

			It("should not be able to mount", func() {
				_, err := client.Mount(logger, "fakedriver", "fake-volume", map[string]interface{}{"volume_id": "fake-volume"})
				Expect(err).To(HaveOccurred())
			})

			It("should not be able to unmount", func() {
				err := client.Unmount(logger, "fakedriver", "fake-volume")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("after creating successfully driver is not found", func() {
			BeforeEach(func() {
				fakedriverProcess = ginkgomon.Invoke(fakedriverRunner)
				time.Sleep(time.Millisecond * 1000)

				fakeDriverFactory = new(volmanfakes.FakeDriverFactory)
				fakeDriver = new(volmanfakes.FakeDriver)

				fakeDriverFactory.DriverReturns(fakeDriver, nil)

				client, runner = vollocal.NewLocalClient(fakeDriverFactory, *registry)
				Expect(runner).ToNot(BeNil())

				calls := 0
				fakeDriverFactory.DriverStub = func(lager.Logger, string) (voldriver.Driver, error) {
					calls++
					if calls > 1 {
						return nil, fmt.Errorf("driver not found")
					}
					return fakeDriver, nil
				}
			})

			It("should not be able to mount", func() {
				_, err := client.Mount(logger, "fakedriver", "fake-volume", map[string]interface{}{"volume_id": "fake-volume"})
				Expect(err).To(HaveOccurred())
			})

		})

		Context("after unsuccessfully creating", func() {
			BeforeEach(func() {
				fakedriverProcess = ginkgomon.Invoke(fakedriverRunner)
				fakeDriver = new(volmanfakes.FakeDriver)

				fakeDriverFactory = new(volmanfakes.FakeDriverFactory)
				fakeDriverFactory.DriverReturns(fakeDriver, nil)

				fakeDriver.CreateReturns(voldriver.ErrorResponse{"create fails"})

				client, runner = vollocal.NewLocalClient(fakeDriverFactory, *registry)
				Expect(runner).ToNot(BeNil())
			})

			It("should not be able to mount", func() {
				_, err := client.Mount(logger, "fakedriver", "fake-volume", map[string]interface{}{"volume_id": "fake-volume"})
				Expect(err).To(HaveOccurred())
			})

		})
	})

})
