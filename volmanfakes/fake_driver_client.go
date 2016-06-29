// This file was generated by counterfeiter
package volmanfakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/voldriver"
	"github.com/pivotal-golang/lager"
)

type FakeDriver struct {
	ActivateStub        func(logger lager.Logger) voldriver.ActivateResponse
	activateMutex       sync.RWMutex
	activateArgsForCall []struct {
		logger lager.Logger
	}
	activateReturns struct {
		result1 voldriver.ActivateResponse
	}
	CreateStub        func(logger lager.Logger, createRequest voldriver.CreateRequest) voldriver.ErrorResponse
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		logger        lager.Logger
		createRequest voldriver.CreateRequest
	}
	createReturns struct {
		result1 voldriver.ErrorResponse
	}
	GetStub        func(logger lager.Logger, getRequest voldriver.GetRequest) voldriver.GetResponse
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		logger     lager.Logger
		getRequest voldriver.GetRequest
	}
	getReturns struct {
		result1 voldriver.GetResponse
	}
	ListStub        func(logger lager.Logger) voldriver.ListResponse
	listMutex       sync.RWMutex
	listArgsForCall []struct {
		logger lager.Logger
	}
	listReturns struct {
		result1 voldriver.ListResponse
	}
	MountStub        func(logger lager.Logger, mountRequest voldriver.MountRequest) voldriver.MountResponse
	mountMutex       sync.RWMutex
	mountArgsForCall []struct {
		logger       lager.Logger
		mountRequest voldriver.MountRequest
	}
	mountReturns struct {
		result1 voldriver.MountResponse
	}
	PathStub        func(logger lager.Logger, pathRequest voldriver.PathRequest) voldriver.PathResponse
	pathMutex       sync.RWMutex
	pathArgsForCall []struct {
		logger      lager.Logger
		pathRequest voldriver.PathRequest
	}
	pathReturns struct {
		result1 voldriver.PathResponse
	}
	RemoveStub        func(logger lager.Logger, removeRequest voldriver.RemoveRequest) voldriver.ErrorResponse
	removeMutex       sync.RWMutex
	removeArgsForCall []struct {
		logger        lager.Logger
		removeRequest voldriver.RemoveRequest
	}
	removeReturns struct {
		result1 voldriver.ErrorResponse
	}
	UnmountStub        func(logger lager.Logger, unmountRequest voldriver.UnmountRequest) voldriver.ErrorResponse
	unmountMutex       sync.RWMutex
	unmountArgsForCall []struct {
		logger         lager.Logger
		unmountRequest voldriver.UnmountRequest
	}
	unmountReturns struct {
		result1 voldriver.ErrorResponse
	}
}

func (fake *FakeDriver) Activate(logger lager.Logger) voldriver.ActivateResponse {
	fake.activateMutex.Lock()
	fake.activateArgsForCall = append(fake.activateArgsForCall, struct {
		logger lager.Logger
	}{logger})
	fake.activateMutex.Unlock()
	if fake.ActivateStub != nil {
		return fake.ActivateStub(logger)
	} else {
		return fake.activateReturns.result1
	}
}

func (fake *FakeDriver) ActivateCallCount() int {
	fake.activateMutex.RLock()
	defer fake.activateMutex.RUnlock()
	return len(fake.activateArgsForCall)
}

func (fake *FakeDriver) ActivateArgsForCall(i int) lager.Logger {
	fake.activateMutex.RLock()
	defer fake.activateMutex.RUnlock()
	return fake.activateArgsForCall[i].logger
}

func (fake *FakeDriver) ActivateReturns(result1 voldriver.ActivateResponse) {
	fake.ActivateStub = nil
	fake.activateReturns = struct {
		result1 voldriver.ActivateResponse
	}{result1}
}

func (fake *FakeDriver) Create(logger lager.Logger, createRequest voldriver.CreateRequest) voldriver.ErrorResponse {
	fake.createMutex.Lock()
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		logger        lager.Logger
		createRequest voldriver.CreateRequest
	}{logger, createRequest})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(logger, createRequest)
	} else {
		return fake.createReturns.result1
	}
}

func (fake *FakeDriver) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeDriver) CreateArgsForCall(i int) (lager.Logger, voldriver.CreateRequest) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].logger, fake.createArgsForCall[i].createRequest
}

func (fake *FakeDriver) CreateReturns(result1 voldriver.ErrorResponse) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 voldriver.ErrorResponse
	}{result1}
}

func (fake *FakeDriver) Get(logger lager.Logger, getRequest voldriver.GetRequest) voldriver.GetResponse {
	fake.getMutex.Lock()
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		logger     lager.Logger
		getRequest voldriver.GetRequest
	}{logger, getRequest})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(logger, getRequest)
	} else {
		return fake.getReturns.result1
	}
}

func (fake *FakeDriver) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeDriver) GetArgsForCall(i int) (lager.Logger, voldriver.GetRequest) {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].logger, fake.getArgsForCall[i].getRequest
}

func (fake *FakeDriver) GetReturns(result1 voldriver.GetResponse) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 voldriver.GetResponse
	}{result1}
}

func (fake *FakeDriver) List(logger lager.Logger) voldriver.ListResponse {
	fake.listMutex.Lock()
	fake.listArgsForCall = append(fake.listArgsForCall, struct {
		logger lager.Logger
	}{logger})
	fake.listMutex.Unlock()
	if fake.ListStub != nil {
		return fake.ListStub(logger)
	} else {
		return fake.listReturns.result1
	}
}

func (fake *FakeDriver) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *FakeDriver) ListArgsForCall(i int) lager.Logger {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return fake.listArgsForCall[i].logger
}

func (fake *FakeDriver) ListReturns(result1 voldriver.ListResponse) {
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 voldriver.ListResponse
	}{result1}
}

func (fake *FakeDriver) Mount(logger lager.Logger, mountRequest voldriver.MountRequest) voldriver.MountResponse {
	fake.mountMutex.Lock()
	fake.mountArgsForCall = append(fake.mountArgsForCall, struct {
		logger       lager.Logger
		mountRequest voldriver.MountRequest
	}{logger, mountRequest})
	fake.mountMutex.Unlock()
	if fake.MountStub != nil {
		return fake.MountStub(logger, mountRequest)
	} else {
		return fake.mountReturns.result1
	}
}

func (fake *FakeDriver) MountCallCount() int {
	fake.mountMutex.RLock()
	defer fake.mountMutex.RUnlock()
	return len(fake.mountArgsForCall)
}

func (fake *FakeDriver) MountArgsForCall(i int) (lager.Logger, voldriver.MountRequest) {
	fake.mountMutex.RLock()
	defer fake.mountMutex.RUnlock()
	return fake.mountArgsForCall[i].logger, fake.mountArgsForCall[i].mountRequest
}

func (fake *FakeDriver) MountReturns(result1 voldriver.MountResponse) {
	fake.MountStub = nil
	fake.mountReturns = struct {
		result1 voldriver.MountResponse
	}{result1}
}

func (fake *FakeDriver) Path(logger lager.Logger, pathRequest voldriver.PathRequest) voldriver.PathResponse {
	fake.pathMutex.Lock()
	fake.pathArgsForCall = append(fake.pathArgsForCall, struct {
		logger      lager.Logger
		pathRequest voldriver.PathRequest
	}{logger, pathRequest})
	fake.pathMutex.Unlock()
	if fake.PathStub != nil {
		return fake.PathStub(logger, pathRequest)
	} else {
		return fake.pathReturns.result1
	}
}

func (fake *FakeDriver) PathCallCount() int {
	fake.pathMutex.RLock()
	defer fake.pathMutex.RUnlock()
	return len(fake.pathArgsForCall)
}

func (fake *FakeDriver) PathArgsForCall(i int) (lager.Logger, voldriver.PathRequest) {
	fake.pathMutex.RLock()
	defer fake.pathMutex.RUnlock()
	return fake.pathArgsForCall[i].logger, fake.pathArgsForCall[i].pathRequest
}

func (fake *FakeDriver) PathReturns(result1 voldriver.PathResponse) {
	fake.PathStub = nil
	fake.pathReturns = struct {
		result1 voldriver.PathResponse
	}{result1}
}

func (fake *FakeDriver) Remove(logger lager.Logger, removeRequest voldriver.RemoveRequest) voldriver.ErrorResponse {
	fake.removeMutex.Lock()
	fake.removeArgsForCall = append(fake.removeArgsForCall, struct {
		logger        lager.Logger
		removeRequest voldriver.RemoveRequest
	}{logger, removeRequest})
	fake.removeMutex.Unlock()
	if fake.RemoveStub != nil {
		return fake.RemoveStub(logger, removeRequest)
	} else {
		return fake.removeReturns.result1
	}
}

func (fake *FakeDriver) RemoveCallCount() int {
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	return len(fake.removeArgsForCall)
}

func (fake *FakeDriver) RemoveArgsForCall(i int) (lager.Logger, voldriver.RemoveRequest) {
	fake.removeMutex.RLock()
	defer fake.removeMutex.RUnlock()
	return fake.removeArgsForCall[i].logger, fake.removeArgsForCall[i].removeRequest
}

func (fake *FakeDriver) RemoveReturns(result1 voldriver.ErrorResponse) {
	fake.RemoveStub = nil
	fake.removeReturns = struct {
		result1 voldriver.ErrorResponse
	}{result1}
}

func (fake *FakeDriver) Unmount(logger lager.Logger, unmountRequest voldriver.UnmountRequest) voldriver.ErrorResponse {
	fake.unmountMutex.Lock()
	fake.unmountArgsForCall = append(fake.unmountArgsForCall, struct {
		logger         lager.Logger
		unmountRequest voldriver.UnmountRequest
	}{logger, unmountRequest})
	fake.unmountMutex.Unlock()
	if fake.UnmountStub != nil {
		return fake.UnmountStub(logger, unmountRequest)
	} else {
		return fake.unmountReturns.result1
	}
}

func (fake *FakeDriver) UnmountCallCount() int {
	fake.unmountMutex.RLock()
	defer fake.unmountMutex.RUnlock()
	return len(fake.unmountArgsForCall)
}

func (fake *FakeDriver) UnmountArgsForCall(i int) (lager.Logger, voldriver.UnmountRequest) {
	fake.unmountMutex.RLock()
	defer fake.unmountMutex.RUnlock()
	return fake.unmountArgsForCall[i].logger, fake.unmountArgsForCall[i].unmountRequest
}

func (fake *FakeDriver) UnmountReturns(result1 voldriver.ErrorResponse) {
	fake.UnmountStub = nil
	fake.unmountReturns = struct {
		result1 voldriver.ErrorResponse
	}{result1}
}

var _ voldriver.Driver = new(FakeDriver)
