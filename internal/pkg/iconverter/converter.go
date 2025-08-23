package iconverter

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// 缓存结构体字段映射
var fieldCache sync.Map

type fieldMap struct {
	srcIndex []int
	dstIndex []int
}

// SmartConvert 智能转换器
func SmartConvert(src, dst interface{}) {
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() == reflect.Ptr && dstVal.IsNil() {
		fmt.Println("result 是一个指向 nil 的指针")
		return
	}
	dstVal = dstVal.Elem()
	srcVal := reflect.ValueOf(src).Elem()

	srcType := srcVal.Type()
	dstType := dstVal.Type()

	// 尝试从缓存获取字段映射
	cacheKey := srcType.String() + "->" + dstType.String()
	if cached, ok := fieldCache.Load(cacheKey); ok {
		fm := cached.(*fieldMap)
		for i, srcIdx := range fm.srcIndex {
			dstIdx := fm.dstIndex[i]
			dstField := dstVal.Field(dstIdx)
			srcField := srcVal.Field(srcIdx)

			if dstField.CanSet() && dstField.Kind() == srcField.Kind() {
				dstField.Set(srcField)
			}
		}
		return
	}

	// 没有缓存，创建新的映射
	fm := &fieldMap{}

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Type().Field(i)
		dstField, ok := dstType.FieldByName(srcField.Name)

		if !ok {
			continue
		}

		// 特殊处理时间类型
		if srcField.Type == reflect.TypeOf(time.Time{}) {
			if dstField.Type == reflect.TypeOf(time.Time{}) {
				fm.srcIndex = append(fm.srcIndex, i)
				fm.dstIndex = append(fm.dstIndex, dstField.Index[0])
			}
			continue
		}

		// 基本类型匹配
		if srcField.Type.Kind() == dstField.Type.Kind() {
			fm.srcIndex = append(fm.srcIndex, i)
			fm.dstIndex = append(fm.dstIndex, dstField.Index[0])
		}
	}

	// 保存到缓存
	fieldCache.Store(cacheKey, fm)

	// 应用转换
	for i, srcIdx := range fm.srcIndex {
		dstIdx := fm.dstIndex[i]
		dstVal.Field(dstIdx).Set(srcVal.Field(srcIdx))
	}
}
