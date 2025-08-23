package iutils

import "strings"

func GetDistinctFields[T any, Y int | int64 | int32 | string](entities []T, getField func(T) Y) []Y {
	m := make(map[Y]struct{}, len(entities)) // 转换为切片用于查询
	res := make([]Y, 0, len(m))
	for _, re := range entities {
		m[getField(re)] = struct{}{}
	}
	if len(m) == 0 {
		// 没有有效的 ID，直接返回空 map 即可
		return res
	}

	for id := range m {
		res = append(res, id)
	}
	return res
}

// 计算两个集合中哪些是新增，修改和删除
func DiffEntities[T1, T2 any, Y comparable](
	dbList []T1,
	inputList []T2,
	getId1 func(T1) Y,
	getId2 func(T2) Y,
	mapFunc func(T1, T2) T1,
) (toCreate, toUpdate []T1, toDelete []Y) {
	dbMap := make(map[Y]T1, len(dbList))
	inputMap := make(map[Y]T2, len(inputList))

	// 构建 map
	for _, dbEnt := range dbList {
		dbMap[getId1(dbEnt)] = dbEnt
	}
	for _, inputEnt := range inputList {
		id := getId2(inputEnt)
		if id == zero[Y]() {
			var nilT1 T1
			toCreate = append(toCreate, mapFunc(nilT1, inputEnt))
		} else {
			inputMap[getId2(inputEnt)] = inputEnt
		}
	}
	if len(dbMap) > 0 {
		// 查找需要更新和删除的
		for id, dbEnt := range dbMap {
			if inputEnt, exists := inputMap[id]; exists {
				toUpdate = append(toUpdate, mapFunc(dbEnt, inputEnt))
				delete(inputMap, id) // 已处理，避免重复判断
			} else {
				// 数据库中有，输入中无 → 删除
				toDelete = append(toDelete, getId1(dbEnt))
			}
		}
	}
	return toCreate, toUpdate, toDelete
}

func zero[T comparable]() T {
	var t T
	return t
}

func FindCreateAndUpdate[T1, T2 any](
	dbList []T1,
	inputList []T2,
	compareFunc func(T1, T2) bool,
	mapFunc func(T1, T2) T1,
) (toCreate, toUpdate []T1) {

	for _, input := range inputList {
		found := false
		for _, db := range dbList {
			if compareFunc(db, input) {
				// 匹配成功：更新
				toUpdate = append(toUpdate, mapFunc(db, input))
				found = true
				break
			}
		}
		if !found {
			// 没有匹配项：新增
			var nilT1 T1
			toCreate = append(toCreate, mapFunc(nilT1, input))
		}
	}

	return toCreate, toUpdate
}

// OrderToLetter 将从 0 开始的整数转换为 A/B/C/D...
func OrderToLetter(order int32) string {
	if order < 0 {
		return ""
	}
	return string(rune('A' + order))
}

// LetterToOrder 将 A/B/C/D... 转换为从 0 开始的整数
func LetterToOrder(letter string) int32 {
	if letter == "" {
		return -1
	}
	letter = strings.ToUpper(letter)
	r := rune(letter[0])
	if r < 'A' || r > 'Z' {
		return -1
	}
	return int32(r - 'A')
}
