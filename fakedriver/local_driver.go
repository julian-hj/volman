package fakedriver

import (
	"fmt"
	"os"

	"strings"

	"github.com/cloudfoundry-incubator/volman/voldriver"
	"github.com/pivotal-golang/lager"
)

const RootDir = "_fakedriver/"

type localDriver struct { // see voldriver.resources.go
	volumes    map[string]*volume
	fileSystem FileSystem
}

type volume struct {
	volumeID   string
	mountpoint string
}

func NewLocalDriver(fileSystem FileSystem) *localDriver {
	return &localDriver{
		volumes:    map[string]*volume{},
		fileSystem: fileSystem,
	}
}

func (d *localDriver) Info(logger lager.Logger) (voldriver.InfoResponse, error) {
	return voldriver.InfoResponse{
		Name: "fakedriver",
		Path: "/fake/path",
	}, nil
}

func (d *localDriver) Mount(logger lager.Logger, mountRequest voldriver.MountRequest) voldriver.MountResponse {
	var vol *volume
	var ok bool
	if vol, ok = d.volumes[mountRequest.Name]; !ok {
		return voldriver.MountResponse{Err: fmt.Sprintf("Volume '%s' must be created before being mounted", mountRequest.Name)}
	}

	mountPath := d.mountPath(vol.volumeID)

	logger.Info("mounting-volume", lager.Data{"id": vol.volumeID, "mountpoint": mountPath})
	err := d.fileSystem.MkdirAll(mountPath, 0777)
	if err != nil {
		logger.Error("failed-creating-mountpoint", err)
		return voldriver.MountResponse{Err: fmt.Sprintf("Error mounting volume: %s", err.Error())}
	}

	vol.mountpoint = mountPath

	mountResponse := voldriver.MountResponse{Mountpoint: mountPath}
	return mountResponse
}

func (d *localDriver) Unmount(logger lager.Logger, unmountRequest voldriver.UnmountRequest) error {
	mountPath := d.mountPath(unmountRequest.VolumeId)

	exists, err := exists(mountPath)
	if err != nil {
		logger.Error("failed-retrieving-mount-info", err, lager.Data{"mountpoint": mountPath})
		return fmt.Errorf("Error establishing whether volume exists")
	}
	if !exists {
		logger.Info(fmt.Sprintf("Volume %s does not exist, nothing to do!", unmountRequest.VolumeId))
		logger.Error("mountpoint-not-found", nil, unmountRequest.VolumeId))
		return fmt.Errorf("Volume %s does not exist, nothing to do!", unmountRequest.VolumeId)
	} else {
		logger.Info(fmt.Sprintf("Removing volume path %s", mountPath))
		err := os.RemoveAll(mountPath)
		if err != nil {
			logger.Info(fmt.Sprintf("Unexpected error removing mount path %s", unmountRequest.VolumeId))
			return fmt.Errorf("Unexpected error removing mount path %s", unmountRequest.VolumeId)
		}
		logger.Info(fmt.Sprintf("Unmounted volume %s", unmountRequest.VolumeId))
	}
	return nil
}

func (d *localDriver) Create(logger lager.Logger, createRequest voldriver.CreateRequest) voldriver.ErrorResponse {
	logger = logger.Session("create")
	if id, ok := createRequest.Opts["volume_id"]; ok {
		logger.Info("creating-volume", lager.Data{"volume_name": createRequest.Name, "volume_id": id})

		if v, ok := d.volumes[createRequest.Name]; ok {
			// If a volume with the given name already exists, no-op unless the opts are different
			if v.volumeID != id {
				logger.Info("duplicate-volume", lager.Data{"volume_name": createRequest.Name})
				return voldriver.ErrorResponse{Err: fmt.Sprintf("Volume '%s' already exists with a different volume ID", createRequest.Name)}
			}

			return voldriver.ErrorResponse{}
		}

		d.volumes[createRequest.Name] = &volume{volumeID: id.(string)}
		return voldriver.ErrorResponse{}
	}

	logger.Info("missing-volume-id", lager.Data{"volume_name": createRequest.Name})
	return voldriver.ErrorResponse{Err: "Missing mandatory 'volume_id' field in 'Opts'"}
}

func (d *localDriver) Get(logger lager.Logger, getRequest voldriver.GetRequest) voldriver.GetResponse {
	if vol, ok := d.volumes[getRequest.Name]; ok {
		return voldriver.GetResponse{Volume: voldriver.VolumeInfo{Name: getRequest.Name, Mountpoint: vol.mountpoint}}
	}

	return voldriver.GetResponse{Err: "Volume not found"}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (d *localDriver) mountPath(volumeId string) string {

	tmpDir := d.fileSystem.TempDir()
	if !strings.HasSuffix(tmpDir, "/") {
		tmpDir = fmt.Sprintf("%s/", tmpDir)
	}

	return tmpDir + RootDir + volumeId
}
