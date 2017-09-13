// Code generated by counterfeiter. DO NOT EDIT.
package volmanfakes

import (
	"sync"

	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/volman"
)

type FakePlugin struct {
	ListVolumesStub        func(logger lager.Logger) ([]string, error)
	listVolumesMutex       sync.RWMutex
	listVolumesArgsForCall []struct {
		logger lager.Logger
	}
	listVolumesReturns struct {
		result1 []string
		result2 error
	}
	listVolumesReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	MountStub        func(logger lager.Logger, volumeId string, config map[string]interface{}) (volman.MountResponse, error)
	mountMutex       sync.RWMutex
	mountArgsForCall []struct {
		logger   lager.Logger
		volumeId string
		config   map[string]interface{}
	}
	mountReturns struct {
		result1 volman.MountResponse
		result2 error
	}
	mountReturnsOnCall map[int]struct {
		result1 volman.MountResponse
		result2 error
	}
	UnmountStub        func(logger lager.Logger, volumeId string) error
	unmountMutex       sync.RWMutex
	unmountArgsForCall []struct {
		logger   lager.Logger
		volumeId string
	}
	unmountReturns struct {
		result1 error
	}
	unmountReturnsOnCall map[int]struct {
		result1 error
	}
	MatchesStub        func(lager.Logger, volman.PluginSpec) bool
	matchesMutex       sync.RWMutex
	matchesArgsForCall []struct {
		arg1 lager.Logger
		arg2 volman.PluginSpec
	}
	matchesReturns struct {
		result1 bool
	}
	matchesReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePlugin) ListVolumes(logger lager.Logger) ([]string, error) {
	fake.listVolumesMutex.Lock()
	ret, specificReturn := fake.listVolumesReturnsOnCall[len(fake.listVolumesArgsForCall)]
	fake.listVolumesArgsForCall = append(fake.listVolumesArgsForCall, struct {
		logger lager.Logger
	}{logger})
	fake.recordInvocation("ListVolumes", []interface{}{logger})
	fake.listVolumesMutex.Unlock()
	if fake.ListVolumesStub != nil {
		return fake.ListVolumesStub(logger)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.listVolumesReturns.result1, fake.listVolumesReturns.result2
}

func (fake *FakePlugin) ListVolumesCallCount() int {
	fake.listVolumesMutex.RLock()
	defer fake.listVolumesMutex.RUnlock()
	return len(fake.listVolumesArgsForCall)
}

func (fake *FakePlugin) ListVolumesArgsForCall(i int) lager.Logger {
	fake.listVolumesMutex.RLock()
	defer fake.listVolumesMutex.RUnlock()
	return fake.listVolumesArgsForCall[i].logger
}

func (fake *FakePlugin) ListVolumesReturns(result1 []string, result2 error) {
	fake.ListVolumesStub = nil
	fake.listVolumesReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakePlugin) ListVolumesReturnsOnCall(i int, result1 []string, result2 error) {
	fake.ListVolumesStub = nil
	if fake.listVolumesReturnsOnCall == nil {
		fake.listVolumesReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.listVolumesReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakePlugin) Mount(logger lager.Logger, volumeId string, config map[string]interface{}) (volman.MountResponse, error) {
	fake.mountMutex.Lock()
	ret, specificReturn := fake.mountReturnsOnCall[len(fake.mountArgsForCall)]
	fake.mountArgsForCall = append(fake.mountArgsForCall, struct {
		logger   lager.Logger
		volumeId string
		config   map[string]interface{}
	}{logger, volumeId, config})
	fake.recordInvocation("Mount", []interface{}{logger, volumeId, config})
	fake.mountMutex.Unlock()
	if fake.MountStub != nil {
		return fake.MountStub(logger, volumeId, config)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.mountReturns.result1, fake.mountReturns.result2
}

func (fake *FakePlugin) MountCallCount() int {
	fake.mountMutex.RLock()
	defer fake.mountMutex.RUnlock()
	return len(fake.mountArgsForCall)
}

func (fake *FakePlugin) MountArgsForCall(i int) (lager.Logger, string, map[string]interface{}) {
	fake.mountMutex.RLock()
	defer fake.mountMutex.RUnlock()
	return fake.mountArgsForCall[i].logger, fake.mountArgsForCall[i].volumeId, fake.mountArgsForCall[i].config
}

func (fake *FakePlugin) MountReturns(result1 volman.MountResponse, result2 error) {
	fake.MountStub = nil
	fake.mountReturns = struct {
		result1 volman.MountResponse
		result2 error
	}{result1, result2}
}

func (fake *FakePlugin) MountReturnsOnCall(i int, result1 volman.MountResponse, result2 error) {
	fake.MountStub = nil
	if fake.mountReturnsOnCall == nil {
		fake.mountReturnsOnCall = make(map[int]struct {
			result1 volman.MountResponse
			result2 error
		})
	}
	fake.mountReturnsOnCall[i] = struct {
		result1 volman.MountResponse
		result2 error
	}{result1, result2}
}

func (fake *FakePlugin) Unmount(logger lager.Logger, volumeId string) error {
	fake.unmountMutex.Lock()
	ret, specificReturn := fake.unmountReturnsOnCall[len(fake.unmountArgsForCall)]
	fake.unmountArgsForCall = append(fake.unmountArgsForCall, struct {
		logger   lager.Logger
		volumeId string
	}{logger, volumeId})
	fake.recordInvocation("Unmount", []interface{}{logger, volumeId})
	fake.unmountMutex.Unlock()
	if fake.UnmountStub != nil {
		return fake.UnmountStub(logger, volumeId)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.unmountReturns.result1
}

func (fake *FakePlugin) UnmountCallCount() int {
	fake.unmountMutex.RLock()
	defer fake.unmountMutex.RUnlock()
	return len(fake.unmountArgsForCall)
}

func (fake *FakePlugin) UnmountArgsForCall(i int) (lager.Logger, string) {
	fake.unmountMutex.RLock()
	defer fake.unmountMutex.RUnlock()
	return fake.unmountArgsForCall[i].logger, fake.unmountArgsForCall[i].volumeId
}

func (fake *FakePlugin) UnmountReturns(result1 error) {
	fake.UnmountStub = nil
	fake.unmountReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePlugin) UnmountReturnsOnCall(i int, result1 error) {
	fake.UnmountStub = nil
	if fake.unmountReturnsOnCall == nil {
		fake.unmountReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.unmountReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePlugin) Matches(arg1 lager.Logger, arg2 volman.PluginSpec) bool {
	fake.matchesMutex.Lock()
	ret, specificReturn := fake.matchesReturnsOnCall[len(fake.matchesArgsForCall)]
	fake.matchesArgsForCall = append(fake.matchesArgsForCall, struct {
		arg1 lager.Logger
		arg2 volman.PluginSpec
	}{arg1, arg2})
	fake.recordInvocation("Matches", []interface{}{arg1, arg2})
	fake.matchesMutex.Unlock()
	if fake.MatchesStub != nil {
		return fake.MatchesStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.matchesReturns.result1
}

func (fake *FakePlugin) MatchesCallCount() int {
	fake.matchesMutex.RLock()
	defer fake.matchesMutex.RUnlock()
	return len(fake.matchesArgsForCall)
}

func (fake *FakePlugin) MatchesArgsForCall(i int) (lager.Logger, volman.PluginSpec) {
	fake.matchesMutex.RLock()
	defer fake.matchesMutex.RUnlock()
	return fake.matchesArgsForCall[i].arg1, fake.matchesArgsForCall[i].arg2
}

func (fake *FakePlugin) MatchesReturns(result1 bool) {
	fake.MatchesStub = nil
	fake.matchesReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakePlugin) MatchesReturnsOnCall(i int, result1 bool) {
	fake.MatchesStub = nil
	if fake.matchesReturnsOnCall == nil {
		fake.matchesReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.matchesReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakePlugin) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listVolumesMutex.RLock()
	defer fake.listVolumesMutex.RUnlock()
	fake.mountMutex.RLock()
	defer fake.mountMutex.RUnlock()
	fake.unmountMutex.RLock()
	defer fake.unmountMutex.RUnlock()
	fake.matchesMutex.RLock()
	defer fake.matchesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePlugin) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ volman.Plugin = new(FakePlugin)