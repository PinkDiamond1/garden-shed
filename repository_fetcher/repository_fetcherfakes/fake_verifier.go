// This file was generated by counterfeiter
package repository_fetcherfakes

import (
	"io"
	"sync"

	"code.cloudfoundry.org/garden-shed/repository_fetcher"
	"github.com/docker/distribution/digest"
)

type FakeVerifier struct {
	VerifyStub        func(io.Reader, digest.Digest) (io.ReadCloser, error)
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		arg1 io.Reader
		arg2 digest.Digest
	}
	verifyReturns struct {
		result1 io.ReadCloser
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeVerifier) Verify(arg1 io.Reader, arg2 digest.Digest) (io.ReadCloser, error) {
	fake.verifyMutex.Lock()
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		arg1 io.Reader
		arg2 digest.Digest
	}{arg1, arg2})
	fake.recordInvocation("Verify", []interface{}{arg1, arg2})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(arg1, arg2)
	} else {
		return fake.verifyReturns.result1, fake.verifyReturns.result2
	}
}

func (fake *FakeVerifier) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *FakeVerifier) VerifyArgsForCall(i int) (io.Reader, digest.Digest) {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.verifyArgsForCall[i].arg1, fake.verifyArgsForCall[i].arg2
}

func (fake *FakeVerifier) VerifyReturns(result1 io.ReadCloser, result2 error) {
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 io.ReadCloser
		result2 error
	}{result1, result2}
}

func (fake *FakeVerifier) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeVerifier) recordInvocation(key string, args []interface{}) {
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

var _ repository_fetcher.Verifier = new(FakeVerifier)
