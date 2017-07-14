// Copyright (c) 2018 Nomad Media, Inc. All Rights Reserved.


package store

import (
	"time"

	l4g "github.com/alecthomas/log4go"
	"github.com/nomadsingles/platform/model"
)

type StoreResult struct {
	Data interface{}
	Err  *model.AppError
}

type StoreChannel chan StoreResult

func Must(sc StoreChannel) interface{} {
	r := <-sc
	if r.Err != nil {
		l4g.Close()
		time.Sleep(time.Second)
		panic(r.Err)
	}

	return r.Data
}

type Store interface {
	System() SystemStore
	MarkSystemRanUnitTests()
	Close()
	DropAllTables()
	TotalMasterDbConnections() int
	TotalReadDbConnections() int
	TotalSearchDbConnections() int
}







type SystemStore interface {
	Save(system *model.System) StoreChannel
	SaveOrUpdate(system *model.System) StoreChannel
	Update(system *model.System) StoreChannel
	Get() StoreChannel
	GetByName(name string) StoreChannel
}


