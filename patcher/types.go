package patcher

import (
	"github.com/undefinedlabs/go-mpatch"
	"sync"
	"testing"
)

type patcher struct {
	data map[uintptr]*mpatch.Patch
	lock sync.Mutex
	test *testing.T
}
