package patcher

import "reflect"

type Patcher interface {
	PatchMethod(target, redirection interface{}) uintptr
	PatchInstanceMethodByName(target reflect.Type, methodName string, redirection interface{}) uintptr
	Unpatch(key uintptr)
	UnpatchAll()
}
