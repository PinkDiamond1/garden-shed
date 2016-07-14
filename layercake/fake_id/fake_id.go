// This file was generated by counterfeiter
package fake_id

import (
	"sync"

	"code.cloudfoundry.org/garden-shed/layercake"
)

type FakeID struct {
	GraphIDStub        func() string
	graphIDMutex       sync.RWMutex
	graphIDArgsForCall []struct{}
	graphIDReturns     struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeID) GraphID() string {
	fake.graphIDMutex.Lock()
	fake.graphIDArgsForCall = append(fake.graphIDArgsForCall, struct{}{})
	fake.recordInvocation("GraphID", []interface{}{})
	fake.graphIDMutex.Unlock()
	if fake.GraphIDStub != nil {
		return fake.GraphIDStub()
	} else {
		return fake.graphIDReturns.result1
	}
}

func (fake *FakeID) GraphIDCallCount() int {
	fake.graphIDMutex.RLock()
	defer fake.graphIDMutex.RUnlock()
	return len(fake.graphIDArgsForCall)
}

func (fake *FakeID) GraphIDReturns(result1 string) {
	fake.GraphIDStub = nil
	fake.graphIDReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeID) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.graphIDMutex.RLock()
	defer fake.graphIDMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeID) recordInvocation(key string, args []interface{}) {
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

var _ layercake.ID = new(FakeID)
