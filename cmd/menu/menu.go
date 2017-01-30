package main

import (
	"github.com/asdine/storm"
	"github.com/tomsteele/emptynest"
)

type menu struct {
	DB          *storm.DB
	HostChanMap map[int](chan emptynest.Payload)
	PayloadMap  map[string]emptynest.PayloadPlugin
}
