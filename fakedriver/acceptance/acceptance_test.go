package acceptance_test

import (
	"github.com/cloudfoundry-incubator/volman/certification"
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit/ginkgomon"
)

var _ = Describe("Fake Driver Certification", func() {
	certification.CertifiyWith("Fakedriver TCP", func() (*ginkgomon.Runner, *ginkgomon.Runner, int, string, string, int, string, func() (string, map[string]interface{})) {
		volumeInfo := func() (string, map[string]interface{}) {
			uuid, err := uuid.NewV4()
			Expect(err).NotTo(HaveOccurred())
			volumeId := "fake-volume-id_" + uuid.String()
			volumeName := "fake-volume-name_" + uuid.String()
			opts := map[string]interface{}{"volume_id": volumeId}
			return volumeName, opts
		}

		return driverRunner, volmanRunner, volmanServerPort, debugServerAddress, tmpDriversPath, driverServerPort, "fakedriver", volumeInfo
	})

	certification.CertifiyWith("Fakedriver UNIX", func() (*ginkgomon.Runner, *ginkgomon.Runner, int, string, string, int, string, func() (string, map[string]interface{})) {
		volumeInfo := func() (string, map[string]interface{}) {
			uuid, err := uuid.NewV4()
			Expect(err).NotTo(HaveOccurred())
			volumeId := "fake-volume-id_" + uuid.String()
			volumeName := "fake-volume-name_" + uuid.String()
			opts := map[string]interface{}{"volume_id": volumeId}
			return volumeName, opts
		}

		return unixDriverRunner, volmanRunner, volmanServerPort, debugServerAddress, tmpDriversPath, -1, "fakedriver", volumeInfo
	})
})