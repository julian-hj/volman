package vollocal_test

import (
	"encoding/json"
	"time"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mfakes "code.cloudfoundry.org/diego-logging-client/testhelpers"
	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/voldriver"
	"code.cloudfoundry.org/volman/voldiscoverers"
	"code.cloudfoundry.org/volman/vollocal"
	"code.cloudfoundry.org/volman/volmanfakes"

	"code.cloudfoundry.org/clock/fakeclock"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"
	"code.cloudfoundry.org/voldriver/voldriverfakes"
	"code.cloudfoundry.org/volman"
	"github.com/onsi/gomega/gbytes"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/ginkgomon"
)

var _ = Describe("Volman", func() {
	var (
		logger *lagertest.TestLogger

		fakeDriverFactory *volmanfakes.FakeDockerDriverFactory
		fakeDriver        *voldriverfakes.FakeDriver
		fakeClock         *fakeclock.FakeClock
		fakeMetronClient  *mfakes.FakeIngressClient

		scanInterval time.Duration

		driverRegistry         volman.PluginRegistry
		dockerDriverDiscoverer volman.Discoverer
		durationMetricMap      map[string]time.Duration
		counterMetricMap       map[string]int

		process ifrit.Process
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("client-test")

		fakeDriverFactory = new(volmanfakes.FakeDockerDriverFactory)
		fakeClock = fakeclock.NewFakeClock(time.Unix(123, 456))

		scanInterval = 1 * time.Second

		driverRegistry = vollocal.NewPluginRegistry()
		durationMetricMap = make(map[string]time.Duration)
		counterMetricMap = make(map[string]int)

		fakeMetronClient = new(mfakes.FakeIngressClient)
		fakeMetronClient.SendDurationStub = func(name string, value time.Duration, opts ...loggregator.EmitGaugeOption) error {
			durationMetricMap[name] = value
			return nil
		}
		fakeMetronClient.IncrementCounterStub = func(name string) error {
			value, ok := counterMetricMap[name]
			if ok {
				counterMetricMap[name] = value + 1
			} else {
				counterMetricMap[name] = 1
			}
			return nil
		}

	})

	Describe("ListDrivers", func() {
		BeforeEach(func() {
			dockerDriverDiscoverer = voldiscoverers.NewDockerDriverDiscovererWithDriverFactory(logger, driverRegistry, []string{"/somePath"}, fakeDriverFactory)
			client = vollocal.NewLocalClient(logger, driverRegistry, fakeMetronClient, fakeClock)

			syncer := vollocal.NewSyncer(logger, driverRegistry, []volman.Discoverer{dockerDriverDiscoverer}, scanInterval, fakeClock)
			process = ginkgomon.Invoke(syncer.Runner())
		})

		It("should report empty list of drivers", func() {
			drivers, err := client.ListDrivers(logger)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(drivers.Drivers)).To(Equal(0))
		})

		Context("has no drivers in location", func() {

			BeforeEach(func() {
				fakeDriverFactory = new(volmanfakes.FakeDockerDriverFactory)
			})

			It("should report empty list of drivers", func() {
				drivers, err := client.ListDrivers(logger)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(drivers.Drivers)).To(Equal(0))
			})

			AfterEach(func() {
				ginkgomon.Kill(process)
			})

		})

		Context("has driver in location", func() {
			BeforeEach(func() {
				err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, "fakedriver", "spec", []byte("http://0.0.0.0:8080"))
				Expect(err).NotTo(HaveOccurred())

				dockerDriverDiscoverer = voldiscoverers.NewDockerDriverDiscovererWithDriverFactory(logger, driverRegistry, []string{defaultPluginsDirectory}, fakeDriverFactory)
				client = vollocal.NewLocalClient(logger, driverRegistry, fakeMetronClient, fakeClock)

				fakeDriver := new(voldriverfakes.FakeDriver)
				fakeDriverFactory.DockerDriverReturns(fakeDriver, nil)

				fakeDriver.ActivateReturns(voldriver.ActivateResponse{Implements: []string{"VolumeDriver"}})
			})

			It("should report empty list of drivers", func() {
				drivers, err := client.ListDrivers(logger)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(drivers.Drivers)).To(Equal(0))
			})

			Context("after running drivers discovery", func() {
				BeforeEach(func() {
					syncer := vollocal.NewSyncer(logger, driverRegistry, []volman.Discoverer{dockerDriverDiscoverer}, scanInterval, fakeClock)
					process = ginkgomon.Invoke(syncer.Runner())
				})

				AfterEach(func() {
					ginkgomon.Kill(process)
				})

				It("should report fakedriver", func() {
					drivers, err := client.ListDrivers(logger)
					Expect(err).NotTo(HaveOccurred())
					Expect(len(drivers.Drivers)).ToNot(Equal(0))
					Expect(drivers.Drivers[0].Name).To(Equal("fakedriver"))
				})

			})
		})
	})

	Describe("Mount and Unmount", func() {
		var (
			volumeId string
		)
		BeforeEach(func() {
			volumeId = "fake-volume"
		})
		Context("when given a driver", func() {
			BeforeEach(func() {
				fakeDriverFactory = new(volmanfakes.FakeDockerDriverFactory)
				fakeDriver = new(voldriverfakes.FakeDriver)
				fakeDriverFactory.DockerDriverReturns(fakeDriver, nil)

				drivers := make(map[string]voldriver.Driver)
				drivers["fakedriver"] = fakeDriver

				err := voldriver.WriteDriverSpec(logger, defaultPluginsDirectory, "fakedriver", "spec", []byte(fmt.Sprintf("http://0.0.0.0:%d", fakeDriver)))
				Expect(err).NotTo(HaveOccurred())

				fakeDriver.ActivateReturns(voldriver.ActivateResponse{Implements: []string{"VolumeDriver"}})

				dockerDriverDiscoverer = voldiscoverers.NewDockerDriverDiscovererWithDriverFactory(logger, driverRegistry, []string{defaultPluginsDirectory}, fakeDriverFactory)
				client = vollocal.NewLocalClient(logger, driverRegistry, fakeMetronClient, fakeClock)

			})

			JustBeforeEach(func() {
				syncer := vollocal.NewSyncer(logger, driverRegistry, []volman.Discoverer{dockerDriverDiscoverer}, scanInterval, fakeClock)
				process = ginkgomon.Invoke(syncer.Runner())
			})

			AfterEach(func() {
				ginkgomon.Kill(process)
			})

			Context("mount", func() {
				BeforeEach(func() {
					mountResponse := voldriver.MountResponse{Mountpoint: "/var/vcap/data/mounts/" + volumeId}
					fakeDriver.MountReturns(mountResponse)
				})

				It("should be able to mount without warning", func() {
					mountPath, err := client.Mount(logger, "fakedriver", volumeId, map[string]interface{}{"volume_id": volumeId})
					Expect(err).NotTo(HaveOccurred())
					Expect(mountPath).NotTo(Equal(""))
					Expect(logger.Buffer()).NotTo(gbytes.Say("Invalid or dangerous mountpath"))
				})

				It("should not be able to mount if mount fails", func() {
					mountResponse := voldriver.MountResponse{Err: "an error"}
					fakeDriver.MountReturns(mountResponse)

					_, err := client.Mount(logger, "fakedriver", volumeId, map[string]interface{}{"volume_id": volumeId})
					Expect(err).To(HaveOccurred())
					_, isVolmanSafeError := err.(volman.SafeError)
					Expect(isVolmanSafeError).To(Equal(false))

				})

				It("should wrap voldriver safeError to volman safeError", func() {
					voldriverSafeError := voldriver.SafeError{SafeDescription: "safe-badness"}
					safeErrBytes, err := json.Marshal(voldriverSafeError)
					Expect(err).NotTo(HaveOccurred())
					mountResponse := voldriver.MountResponse{Err: string(safeErrBytes[:])}
					fakeDriver.MountReturns(mountResponse)

					_, err = client.Mount(logger, "fakedriver", volumeId, map[string]interface{}{"volume_id": volumeId})
					Expect(err).To(HaveOccurred())
					_, isVolmanSafeError := err.(volman.SafeError)
					Expect(isVolmanSafeError).To(Equal(true))
				})

				Context("with bad mount path", func() {
					var err error
					BeforeEach(func() {
						mountResponse := voldriver.MountResponse{Mountpoint: "/var/tmp"}
						fakeDriver.MountReturns(mountResponse)
					})

					JustBeforeEach(func() {
						_, err = client.Mount(logger, "fakedriver", volumeId, map[string]interface{}{"volume_id": volumeId})
					})

					It("should return a warning in the log", func() {
						Expect(err).NotTo(HaveOccurred())
						Expect(logger.Buffer()).To(gbytes.Say("Invalid or dangerous mountpath"))
					})
				})

				Context("with metrics", func() {
					It("should emit mount time on successful mount", func() {

						client.Mount(logger, "fakedriver", volumeId, map[string]interface{}{"volume_id": volumeId})

						Eventually(durationMetricMap).Should(HaveKeyWithValue("VolmanMountDuration", Not(BeZero())))
						Eventually(durationMetricMap).Should(HaveKeyWithValue("VolmanMountDurationForfakedriver", Not(BeZero())))
					})

					It("should increment error count on mount failure", func() {
						Expect(counterMetricMap).ShouldNot(HaveKey("VolmanMountErrors"))
						mountResponse := voldriver.MountResponse{Err: "an error"}
						fakeDriver.MountReturns(mountResponse)

						client.Mount(logger, "fakedriver", volumeId, map[string]interface{}{"volume_id": volumeId})
						Expect(counterMetricMap).Should(HaveKeyWithValue("VolmanMountErrors", 1))
					})
				})
			})

			Context("umount", func() {
				It("should be able to unmount", func() {
					err := client.Unmount(logger, "fakedriver", volumeId)
					Expect(err).NotTo(HaveOccurred())
					Expect(fakeDriver.UnmountCallCount()).To(Equal(1))
					Expect(fakeDriver.RemoveCallCount()).To(Equal(0))
				})

				It("should not be able to unmount when driver unmount fails", func() {
					fakeDriver.UnmountReturns(voldriver.ErrorResponse{Err: "unmount failure"})
					err := client.Unmount(logger, "fakedriver", volumeId)
					Expect(err).To(HaveOccurred())

					_, isVolmanSafeError := err.(volman.SafeError)
					Expect(isVolmanSafeError).To(Equal(false))
				})

				It("should wrap voldriver safeError to volman safeError", func() {
					voldriverSafeError := voldriver.SafeError{SafeDescription: "safe-badness"}
					safeErrBytes, err := json.Marshal(voldriverSafeError)
					Expect(err).NotTo(HaveOccurred())
					unmountResponse := voldriver.ErrorResponse{Err: string(safeErrBytes[:])}
					fakeDriver.UnmountReturns(unmountResponse)

					err = client.Unmount(logger, "fakedriver", volumeId)
					Expect(err).To(HaveOccurred())
					_, isVolmanSafeError := err.(volman.SafeError)
					Expect(isVolmanSafeError).To(Equal(true))
				})

				Context("with metrics", func() {
					It("should emit unmount time on successful unmount", func() {
						client.Unmount(logger, "fakedriver", volumeId)

						Eventually(durationMetricMap).Should(HaveKeyWithValue("VolmanUnmountDuration", Not(BeZero())))
						Eventually(durationMetricMap).Should(HaveKeyWithValue("VolmanUnmountDurationForfakedriver", Not(BeZero())))
					})

					It("should increment error count on unmount failure", func() {
						fakeDriver.UnmountReturns(voldriver.ErrorResponse{Err: "unmount failure"})

						client.Unmount(logger, "fakedriver", volumeId)
						Expect(counterMetricMap).Should(HaveKeyWithValue("VolmanUnmountErrors", 1))
					})

				})
			})

			Context("when driver is not found", func() {
				BeforeEach(func() {
					fakeDriverFactory.DockerDriverReturns(nil, fmt.Errorf("driver not found"))
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

			Context("when driver does not implement VolumeDriver", func() {
				BeforeEach(func() {
					fakeDriver.ActivateReturns(voldriver.ActivateResponse{Implements: []string{"nada"}})
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
		})

		Context("after creating successfully driver is not found", func() {
			BeforeEach(func() {

				fakeDriverFactory = new(volmanfakes.FakeDockerDriverFactory)
				fakeDriver = new(voldriverfakes.FakeDriver)
				mountReturn := voldriver.MountResponse{Err: "driver not found",
					Mountpoint: "",
				}
				fakeDriver.MountReturns(mountReturn)
				fakeDriverFactory.DockerDriverReturns(fakeDriver, nil)

				driverRegistry := vollocal.NewPluginRegistry()
				dockerDriverDiscoverer = voldiscoverers.NewDockerDriverDiscovererWithDriverFactory(logger, driverRegistry, []string{"/somePath"}, fakeDriverFactory)
				client = vollocal.NewLocalClient(logger, driverRegistry, fakeMetronClient, fakeClock)

				syncer := vollocal.NewSyncer(logger, driverRegistry, []volman.Discoverer{dockerDriverDiscoverer}, scanInterval, fakeClock)
				process = ginkgomon.Invoke(syncer.Runner())

				calls := 0
				fakeDriverFactory.DockerDriverStub = func(lager.Logger, string, string, string) (voldriver.Driver, error) {
					calls++
					if calls > 1 {
						return nil, fmt.Errorf("driver not found")
					}
					return fakeDriver, nil
				}
			})

			AfterEach(func() {
				ginkgomon.Kill(process)
			})

			It("should not be able to mount", func() {
				_, err := client.Mount(logger, "fakedriver", "fake-volume", map[string]interface{}{"volume_id": "fake-volume"})
				Expect(err).To(HaveOccurred())
			})

		})

		Context("after unsuccessfully creating", func() {
			BeforeEach(func() {
				localDriverProcess = ginkgomon.Invoke(localDriverRunner)
				fakeDriver = new(voldriverfakes.FakeDriver)

				fakeDriverFactory = new(volmanfakes.FakeDockerDriverFactory)
				fakeDriverFactory.DockerDriverReturns(fakeDriver, nil)

				fakeDriver.CreateReturns(voldriver.ErrorResponse{"create fails"})

				driverRegistry := vollocal.NewPluginRegistry()
				dockerDriverDiscoverer = voldiscoverers.NewDockerDriverDiscovererWithDriverFactory(logger, driverRegistry, []string{"/somePath"}, fakeDriverFactory)
				client = vollocal.NewLocalClient(logger, driverRegistry, fakeMetronClient, fakeClock)

				syncer := vollocal.NewSyncer(logger, driverRegistry, []volman.Discoverer{dockerDriverDiscoverer}, scanInterval, fakeClock)
				process = ginkgomon.Invoke(syncer.Runner())
			})

			AfterEach(func() {
				ginkgomon.Kill(process)
			})

			It("should not be able to mount", func() {
				_, err := client.Mount(logger, "fakedriver", "fake-volume", map[string]interface{}{"volume_id": "fake-volume"})
				Expect(err).To(HaveOccurred())
			})

		})
	})
})
