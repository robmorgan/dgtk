package main

import (
	"github.com/dynport/dgtk/vmware"
)

func init() {
	router.Register("stop", &StopAction{})
}

type StopAction struct {
	Name string `cli:"type=arg required=true"`
}

func (action *StopAction) Run() error {
	vm, e := findFirst(action.Name)
	if e != nil {
		return e
	}
	return vmware.Stop(vm.Path)
}