package fakedriver_test

import (
	"github.com/cloudfoundry-incubator/volman/fakedriver"
	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/cloudfoundry-incubator/volman/volmanfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager"
	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("Local Driver", func() {
	var logger lager.Logger
	var fakeFileSystem *volmanfakes.FakeFileSystem

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")
		fakeFileSystem = &volmanfakes.FakeFileSystem{}
	})

	Describe("Unmount", func() {

	})

	Describe("Mount", func() {
		Context("when the volume has been created", func() {
			It("mounts the volume on the local filesystem", func() {
				volumeName := "test-volume-name"
				volumeID := "test-volume-id"

				localDriver := fakedriver.NewLocalDriver(fakeFileSystem)
				createResponse := localDriver.Create(logger, voldriver.CreateRequest{
					Name: volumeName,
					Opts: map[string]interface{}{
						"volume_id": volumeID,
					},
				})

				Expect(createResponse.Err).To(Equal(""))

				fakeFileSystem.TempDirReturns("/some/temp/dir/")

				mountResponse := localDriver.Mount(logger, voldriver.MountRequest{
					Name: volumeName,
				})

				Expect(mountResponse.Err).To(Equal(""))
				Expect(mountResponse.Mountpoint).To(Equal("/some/temp/dir/_fakedriver/test-volume-id"))

				Expect(fakeFileSystem.TempDirCallCount()).To(Equal(1))
				Expect(fakeFileSystem.MkdirAllCallCount()).To(Equal(1))
				createdDir, permissions := fakeFileSystem.MkdirAllArgsForCall(0)
				Expect(createdDir).To(Equal("/some/temp/dir/_fakedriver/test-volume-id"))
				Expect(permissions).To(BeEquivalentTo(0777))
			})

			It("returns the mount point on a /VolumeDriver.Get response", func() {
				volumeName := "test-volume-name"
				volumeID := "test-volume-id"

				localDriver := fakedriver.NewLocalDriver(fakeFileSystem)
				createResponse := localDriver.Create(logger, voldriver.CreateRequest{
					Name: volumeName,
					Opts: map[string]interface{}{
						"volume_id": volumeID,
					},
				})
				Expect(createResponse.Err).To(Equal(""))

				fakeFileSystem.TempDirReturns("/some/temp/dir/")

				mountResponse := localDriver.Mount(logger, voldriver.MountRequest{
					Name: volumeName,
				})
				Expect(mountResponse.Err).To(Equal(""))

				getResponse := localDriver.Get(logger, voldriver.GetRequest{
					Name: volumeName,
				})

				Expect(getResponse.Err).To(Equal(""))
				Expect(getResponse.Volume.Mountpoint).To(Equal("/some/temp/dir/_fakedriver/test-volume-id"))
			})
		})

		Context("when the volume has not been created", func() {
			It("returns an error", func() {
				localDriver := fakedriver.NewLocalDriver(fakeFileSystem)
				mountResponse := localDriver.Mount(logger, voldriver.MountRequest{
					Name: "bla",
				})
				Expect(mountResponse.Err).To(Equal("Volume 'bla' must be created before being mounted"))
			})
		})
	})

	Describe("Create", func() {
		Context("when a volume ID is not provided", func() {
			It("returns an error", func() {
				localDriver := fakedriver.NewLocalDriver(fakeFileSystem)
				createResponse := localDriver.Create(logger, voldriver.CreateRequest{
					Name: "volume",
					Opts: map[string]interface{}{
						"nonsense": "bla",
					},
				})

				Expect(createResponse.Err).To(Equal("Missing mandatory 'volume_id' field in 'Opts'"))
			})
		})

		Context("when a second create is called with the same volume ID", func() {
			Context("with the same opts", func() {
				It("does nothing", func() {
					localDriver := fakedriver.NewLocalDriver(fakeFileSystem)
					createResponse := localDriver.Create(logger, voldriver.CreateRequest{
						Name: "volume",
						Opts: map[string]interface{}{
							"volume_id": "bla",
						},
					})

					Expect(createResponse.Err).To(Equal(""))

					createResponse = localDriver.Create(logger, voldriver.CreateRequest{
						Name: "volume",
						Opts: map[string]interface{}{
							"volume_id": "bla",
						},
					})

					Expect(createResponse.Err).To(Equal(""))
				})
			})

			Context("with a different opts", func() {
				It("returns an error", func() {
					localDriver := fakedriver.NewLocalDriver(fakeFileSystem)
					createResponse := localDriver.Create(logger, voldriver.CreateRequest{
						Name: "volume",
						Opts: map[string]interface{}{
							"volume_id": "bla",
						},
					})

					Expect(createResponse.Err).To(Equal(""))

					createResponse = localDriver.Create(logger, voldriver.CreateRequest{
						Name: "volume",
						Opts: map[string]interface{}{
							"volume_id": "foo",
						},
					})

					Expect(createResponse.Err).To(Equal("Volume 'volume' already exists with a different volume ID"))
				})
			})
		})

		Context("when a second create is called with the same volume ID", func() {
			It("does nothing", func() {
				localDriver := fakedriver.NewLocalDriver(fakeFileSystem)
				createResponse := localDriver.Create(logger, voldriver.CreateRequest{
					Name: "volume",
					Opts: map[string]interface{}{
						"nonsense": "bla",
					},
				})

				Expect(createResponse.Err).To(Equal("Missing mandatory 'volume_id' field in 'Opts'"))
			})
		})
	})

	Describe("Get", func() {
		Context("when the volume has been created", func() {
			It("returns the volume name", func() {
				volumeName := "test-volume"

				localDriver := fakedriver.NewLocalDriver(fakeFileSystem)
				createResponse := localDriver.Create(logger, voldriver.CreateRequest{
					Name: volumeName,
					Opts: map[string]interface{}{
						"volume_id": "test",
					},
				})

				Expect(createResponse.Err).To(Equal(""))

				getResponse := localDriver.Get(logger, voldriver.GetRequest{
					Name: volumeName,
				})

				Expect(getResponse.Err).To(Equal(""))
				Expect(getResponse.Volume.Name).To(Equal(volumeName))
			})
		})

		Context("when the volume has not been created", func() {
			It("returns an error", func() {
				volumeName := "test-volume"
				localDriver := fakedriver.NewLocalDriver(fakeFileSystem)
				getResponse := localDriver.Get(logger, voldriver.GetRequest{
					Name: volumeName,
				})

				Expect(getResponse.Err).To(Equal("Volume not found"))
				Expect(getResponse.Volume.Name).To(Equal(""))
			})
		})
	})

})
