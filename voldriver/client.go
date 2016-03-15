package voldriver

import "github.com/pivotal-golang/lager"

//go:generate counterfeiter . Client

type Client interface {
	Create(logger lager.Logger, name string, opts Opts) error
	Remove(logger lager.Logger, name string) error
	Mount(logger lager.Logger, name string) (string, error)
	Path(logger lager.Logger, name string) (string, error)
	Unmount(logger lager.Logger, name string) error
	Get(logger lager.Logger, name string) (Volume, error)
	List(logger lager.Logger, name string) (Volumes, error)
}

//go:generate counterfeiter . Backend

type Backend interface {
	Client
}

type Opts map[string]string

type Volume struct {
	Name       string
	Mountpoint string
}

type Volumes []Volume
