// This file was generated by counterfeiter
package rootfs_providerfakes

import (
	"sync"

	"code.cloudfoundry.org/garden-shed/repository_fetcher"
	"code.cloudfoundry.org/garden-shed/rootfs_provider"
	"code.cloudfoundry.org/lager"
)

type FakeLayerCreator struct {
	CreateStub        func(log lager.Logger, id string, parentImage *repository_fetcher.Image, spec rootfs_provider.Spec) (string, []string, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		log         lager.Logger
		id          string
		parentImage *repository_fetcher.Image
		spec        rootfs_provider.Spec
	}
	createReturns struct {
		result1 string
		result2 []string
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLayerCreator) Create(log lager.Logger, id string, parentImage *repository_fetcher.Image, spec rootfs_provider.Spec) (string, []string, error) {
	fake.createMutex.Lock()
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		log         lager.Logger
		id          string
		parentImage *repository_fetcher.Image
		spec        rootfs_provider.Spec
	}{log, id, parentImage, spec})
	fake.recordInvocation("Create", []interface{}{log, id, parentImage, spec})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(log, id, parentImage, spec)
	} else {
		return fake.createReturns.result1, fake.createReturns.result2, fake.createReturns.result3
	}
}

func (fake *FakeLayerCreator) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeLayerCreator) CreateArgsForCall(i int) (lager.Logger, string, *repository_fetcher.Image, rootfs_provider.Spec) {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].log, fake.createArgsForCall[i].id, fake.createArgsForCall[i].parentImage, fake.createArgsForCall[i].spec
}

func (fake *FakeLayerCreator) CreateReturns(result1 string, result2 []string, result3 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 string
		result2 []string
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeLayerCreator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeLayerCreator) recordInvocation(key string, args []interface{}) {
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

var _ rootfs_provider.LayerCreator = new(FakeLayerCreator)
