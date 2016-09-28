package main

import (
	"sync"
	"time"
)

type ServerState struct {
	m          sync.RWMutex
	salesBegin time.Time
	salesEnd   time.Time
	userGroups map[uint32]uint32
}

func newServerState() *ServerState {
	now := time.Now()
	return &ServerState{
		salesBegin: now,
		salesEnd:   now,
		userGroups: make(map[uint32]uint32),
	}
}

func (state *ServerState) GetUserGroup(uin uint32) uint32 {
	group, exists := state.userGroups[uin]
	if exists {
		return group
	} else {
		return 0
	}
}
