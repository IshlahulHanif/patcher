package patcher

import (
	"github.com/undefinedlabs/go-mpatch"
	"reflect"
	"sync"
	"testing"
)

func New(t *testing.T) *patcher {
	return &patcher{
		data: make(map[uintptr]*mpatch.Patch),
		lock: sync.Mutex{},
		test: t,
	}
}

func (p *patcher) PatchMethod(target, redirection interface{}) uintptr {
	if target == nil || redirection == nil {
		p.test.Errorf("[Error][Patcher][PatchMethod]: target or function redirection must not be nil\n")
	}

	t := reflect.ValueOf(target)

	patch, err := mpatch.PatchMethod(target, redirection)
	if err != nil {
		p.test.Errorf("[Error][Patcher][PatchMethod]: Patch method failed: %v\n", err)
	}

	p.assignPatches(t, patch)

	return t.Pointer()
}

func (p *patcher) PatchInstanceMethodByName(target reflect.Type, methodName string, redirection interface{}) uintptr {
	if target == nil || redirection == nil {
		p.test.Errorf("[Error][Patcher][PatchInstanceMethodByName]: target or function redirection must not be nil\n")
	}

	m, ok := target.MethodByName(methodName)
	if !ok {
		p.test.Errorf("[Error][Patcher][PatchInstanceMethodByName]: unknown method %s\n", methodName)
	}

	patch, err := mpatch.PatchInstanceMethodByName(target, methodName, redirection)
	if err != nil {
		p.test.Errorf("[Error][Patcher][PatchInstanceMethodByName]: Patch instance method by name failed: %v\n", err)
	}

	p.assignPatches(m.Func, patch)

	return m.Func.Pointer()
}

func (p *patcher) Unpatch(key uintptr) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if patch, ok := p.data[key]; ok {
		if err := patch.Unpatch(); err != nil {
			p.test.Errorf("[Error][Patcher][Unpatch]: error unpatch key: %v, reason: %s\n", key, err)
		}
		delete(p.data, key)
	}
}

func (p *patcher) UnpatchAll() {
	p.lock.Lock()
	defer p.lock.Unlock()

	for target, patch := range p.data {
		if err := patch.Unpatch(); err != nil {
			p.test.Errorf("[Error][Patcher][UnpatchAll]: error unpatch %s\n", err)
		}
		delete(p.data, target)
	}
}

func (p *patcher) assignPatches(target reflect.Value, patch *mpatch.Patch) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.data[target.Pointer()] = patch
}
