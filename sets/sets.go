// Copyright 2009 The Go Authors. All rights reserved.
// Copyright 2022 The pM1ng.
// slice 运算

package sets

import (
	"errors"
	"math/rand"
	"reflect"
)

// Contain 元素是否包含在集合內
func Contain(tar, arr interface{}) (exist bool, err error) {
	s := reflect.ValueOf(arr)
	if s.Kind() != reflect.Slice {
		err = errors.New("arr not slice")
		return
	}
	for i := 0; i < s.Len(); i++ {
		item := s.Index(i).Interface()
		if reflect.DeepEqual(tar, item) {
			exist = true
			return
		}
	}
	return
}

//GetRepeatItem 获得重复的元素
func GetRepeatItem(set []int32) (repeats []int32) {
	if len(set) == 0 {
		return
	}
	m := map[int32]int{}
	for _, item := range set {
		m[item]++
	}
	for item, cnt := range m {
		if cnt > 1 {
			repeats = append(repeats, item)
		}
	}
	return
}

// Intersect 求交集
// struct 类型slice 支持根据字段名求交集
func Intersect(aSet, bSet interface{}, key string) (iArr []interface{}) {
	aMap := make(map[interface{}]interface{}, 0)
	bMap := make(map[interface{}]interface{}, 0)
	if key == "" {
		aMap = arr2map(aSet)
		bMap = arr2map(bSet)
	} else {
		aMap = arr2mapWithKey(aSet, key)
		bMap = arr2mapWithKey(bSet, key)
	}
	iArr = IntersectMap(aMap, bMap)
	return
}

// Union 求并集
func Union(aSet, bSet interface{}, key string) (uArr []interface{}) {
	aMap := make(map[interface{}]interface{}, 0)
	bMap := make(map[interface{}]interface{}, 0)
	if key == "" {
		aMap = arr2map(aSet)
		bMap = arr2map(bSet)
	} else {
		aMap = arr2mapWithKey(aSet, key)
		bMap = arr2mapWithKey(bSet, key)
	}
	uArr = UnionMap(aMap, bMap)
	return
}

// Sub 求差集
func Sub(aSet, bSet interface{}, key string) (aSub []interface{}) {
	aMap := make(map[interface{}]interface{}, 0)
	bMap := make(map[interface{}]interface{}, 0)
	if key == "" {
		aMap = arr2map(aSet)
		bMap = arr2map(bSet)
	} else {
		aMap = arr2mapWithKey(aSet, key)
		bMap = arr2mapWithKey(bSet, key)
	}
	aSub = SubMap(aMap, bMap)
	return
}

// DiffSlice aSub：a相较于b的差集；bSub：b相较于a的差集；
// altA, altB: 未指定key时没有意义，指定key时，字段相同值不相同的结果
// a,b: struct slice 可以通过key指定struct中的字段，进行计算结果
func DiffSlice(a, b interface{}, key string) (aSub, bSub, altA, altB []interface{}) {
	aMap := make(map[interface{}]interface{}, 0)
	bMap := make(map[interface{}]interface{}, 0)
	if key == "" {
		aMap = arr2map(a)
		bMap = arr2map(b)
	} else {
		aMap = arr2mapWithKey(a, key)
		bMap = arr2mapWithKey(b, key)
	}
	aSub = SubMap(aMap, bMap)
	bSub = SubMap(bMap, aMap)
	altA, altB = CompareMap(aMap, bMap)
	return
}

// arr2map 切片转map
func arr2map(arr interface{}) (aMap map[interface{}]interface{}) {
	aMap = make(map[interface{}]interface{}, 0)
	s := reflect.ValueOf(arr)
	if s.Kind() != reflect.Slice {
		// todo log
		return
	}
	for i := 0; i < s.Len(); i++ {
		value := s.Index(i).Interface()
		aMap[value] = value
	}
	return
}

// arr2mapWithKey
// struct切片转map 带key（即指定字段名，注意是字段名）
func arr2mapWithKey(arr interface{}, key string) (aMap map[interface{}]interface{}) {
	aMap = make(map[interface{}]interface{}, 0)
	aSlice := reflect.ValueOf(arr)
	if aSlice.Kind() != reflect.Slice {
		// todo log
		return
	}
	for i := 0; i < aSlice.Len(); i++ {
		value := aSlice.Index(i)
		var iv interface{}
		if value.Kind() == reflect.Ptr {
			iv = value.Elem().Interface()
		} else {
			iv = value.Interface()
		}
		if reflect.ValueOf(iv).Kind() != reflect.Struct {
			// todo log
			return
		}
		ikey := reflect.ValueOf(iv).FieldByName(key).Interface()
		aMap[ikey] = iv
	}
	return
}

// GetItemRandomly 随机获取slice一个元素
func GetItemRandomly(aSet interface{}) (value interface{}, err error) {
	aSlice := reflect.ValueOf(aSet)
	if aSlice.Kind() != reflect.Slice {
		err = errors.New("not slice")
		return
	}
	aLen := aSlice.Len()
	if aLen == 0 {
		err = errors.New("slice is empty")
		return
	}
	value = aSlice.Index(rand.Intn(aLen)).Interface()
	return
}

// ConvInt64s 转换成[]int64
func ConvInt64s(src []interface{}) (target []int64, err error) {
	for _, iv := range src {
		v, ok := iv.(int64)
		if !ok {
			err = errors.New("conv to int64 error")
			return
		}
		target = append(target, v)
	}
	return
}

// ConvInt32s 转换成[]int32
func ConvInt32s(src []interface{}) (target []int32, err error) {
	for _, iv := range src {
		v, ok := iv.(int32)
		if !ok {
			err = errors.New("conv to int32 error")
			return
		}
		target = append(target, v)
	}
	return
}

// NMxX 生成组合数
func NMxX(aSet, bSet, cSet, dSet []int32) (aCom [][]int32) {
	if len(aSet) == 0 {
		return
	}
	for _, arg1 := range aSet {
		if len(bSet) == 0 {
			aCom = append(aCom, []int32{arg1})
			continue
		}
		for _, arg2 := range bSet {
			if len(cSet) == 0 {
				aCom = append(aCom, []int32{arg1, arg2})
				continue
			}
			for _, arg3 := range cSet {
				if len(dSet) == 0 {
					aCom = append(aCom, []int32{arg1, arg2, arg3})
					continue
				}
				for _, arg4 := range dSet {
					aCom = append(aCom, []int32{arg1, arg2, arg3, arg4})
				}
			}
		}
	}
	return
}

// EqualInt32s 判断两个 []int32 是否相等
func EqualInt32s(a, b []int32) (isEqual bool) {
	if len(a) != len(b) {
		return
	}
	for i, item := range a {
		if b[i] != item {
			return
		}
	}
	isEqual = true
	return
}
