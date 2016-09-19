package main

import "math/rand"

type UserInfo struct {
	level  uint32
	region uint32
}

type RegionInfo struct {
	region  uint32
	userNum uint32
}

type UserMap map[uint32]UserInfo
type LevelMap map[uint32][]uint32
type RegionMap map[uint32]LevelMap

type ServerState struct {
	userMap   UserMap
	regionMap RegionMap
}

func newServerState() *ServerState {
	return &ServerState{userMap: make(UserMap), regionMap: make(RegionMap)}
}

func (serverState *ServerState) GetOrCreateLevelMap(region uint32) LevelMap {
	m, exists := serverState.regionMap[region]
	if !exists {
		m = make(LevelMap)
		serverState.regionMap[region] = m
	}
	return m
}

func RemoveElementByIndex(uins []uint32, index int) []uint32 {
	if index < 0 || index >= len(uins) {
		panic("Invalid index")
	}
	uins[index] = uins[len(uins)-1]
	uins = uins[:len(uins)-1]
	return uins
}

func RemoveElementByValue(uins []uint32, uin uint32) []uint32 {
	index := -1
	for i, v := range uins {
		if v == uin {
			index = i
			break
		}
	}

	if index != -1 {
		uins = RemoveElementByIndex(uins, index)
	}
	return uins
}

func (levelMap LevelMap) DeleteUser(uin uint32, level uint32) bool {
	uins, exists := levelMap[level]
	if exists {
		levelMap[level] = RemoveElementByValue(uins, uin)
		return true
	}
	return false
}

func (levelMap LevelMap) AppendUser(uin uint32, level uint32) {
	uins, exists := levelMap[level]
	if !exists {
		uins = make([]uint32, 0)
	}
	levelMap[level] = append(uins, uin)
}

func (levelMap LevelMap) RandomPickFromLevel(level uint32,
	expectNum int) (result []uint32) {
	if expectNum <= 0 {
		panic("Invalid expectNum")
	}
	uins, exists := levelMap[level]
	if !exists {
		return result
	}

	if len(uins) <= expectNum {
		return uins
	}

	var choose map[uint32]bool
	for len(choose) != expectNum {
		choose[uint32(rand.Intn(len(uins)))] = true
	}

	for index, _ := range choose {
		result = append(result, uins[index])
	}

	return result
}

func (levelMap LevelMap) GetTotalUinsCount() uint32 {
	var result int = 0
	for _, v := range levelMap {
		result += len(v)
	}
	return uint32(result)
}

func (serverState *ServerState) UpdateUser(uin uint32, level uint32,
	region uint32) {
	// 删除可能存在的旧的数据
	info, exists := serverState.userMap[uin]
	if exists {
		levelMap := serverState.GetOrCreateLevelMap(info.region)
		levelMap.DeleteUser(uin, info.level)
	}

	// 更新为新数据
	info.level = level
	info.region = region
	serverState.userMap[uin] = info

	levelMap := serverState.GetOrCreateLevelMap(region)
	levelMap.AppendUser(uin, level)
}

func (serverState *ServerState) PickUser(selfUin uint32, selfLevel uint32,
	expectRegion uint32) (uins []uint32) {
	levelMap, exists := serverState.regionMap[expectRegion]
	if !exists {
		return uins
	}

	var prev uint32 = 0
	if selfLevel > 0 {
		prev = selfLevel - 1
	}
	next := selfLevel + 1
	const kMaxSearchLevel uint32 = 200
	const kMaxUinNum int = 10

	uins = levelMap.RandomPickFromLevel(selfLevel, kMaxUinNum)
	for len(uins) < kMaxUinNum && !(prev <= 0 && next > kMaxSearchLevel) {
		if prev > 0 {
			result := levelMap.RandomPickFromLevel(prev, kMaxUinNum-len(uins))
			RemoveElementByValue(result, selfUin)
			uins = append(uins, result...)
			prev -= 1
		}

		if len(uins) >= kMaxUinNum {
			break
		}

		if next <= kMaxSearchLevel {
			result := levelMap.RandomPickFromLevel(next, kMaxUinNum-len(uins))
			RemoveElementByValue(result, selfUin)
			uins = append(uins, result...)
			next += 1
		}
	}

	return uins
}

func (serverState *ServerState) GetRegionStatus() (result []RegionInfo) {
	for k, v := range serverState.regionMap {
		result = append(result,
			RegionInfo{region: k, userNum: v.GetTotalUinsCount()})
	}
	return result
}
