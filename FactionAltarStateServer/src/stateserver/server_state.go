package main

import "math/rand"

type StateMap map[uint32]bool
type GroupIndexMap map[uint32]StateMap

type ServerState struct {
	m GroupIndexMap
}

func newServerState() *ServerState {
	var serverState ServerState
	serverState.m = make(GroupIndexMap)
	return &serverState
}

func (serverState *ServerState) GetOrCreateStateMap(groupIndex uint32) StateMap {
	m, exists := serverState.m[groupIndex]
	if !exists {
		m = make(StateMap)
		serverState.m[groupIndex] = m
	}
	return m
}

func (serverState *ServerState) UpdateState(groupIndex uint32, factionID uint32,
	underProtection bool) {
	stateMap := serverState.GetOrCreateStateMap(groupIndex)
	if underProtection {
		delete(stateMap, factionID)
	} else {
		stateMap[factionID] = underProtection
	}
}

func (serverState *ServerState) PickRandomFaction(groupIndex uint32,
	selfFactionID uint32) uint32 {
	if stateMap, exists := serverState.m[groupIndex]; exists && len(stateMap) > 0 {
		const kMaxKeysNum int = 100
		keys := make([]uint32, 0, kMaxKeysNum)
		for k := range stateMap {
			if k != selfFactionID {
				keys = append(keys, k)
			}
			if len(keys) >= kMaxKeysNum {
				break
			}
		}
		if len(keys) == 0 {
			return 0
		} else {
			r := rand.Intn(len(keys))
			return keys[r]
		}
	}
	return 0
}
