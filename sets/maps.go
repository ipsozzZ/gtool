package sets

import (
	"reflect"
)

// IntersectMap 求交集
func IntersectMap(aMap, bMap map[interface{}]interface{}) (inSet []interface{}) {
	for aKey, aValue := range aMap {
		if _, ok := bMap[aKey]; ok {
			inSet = append(inSet, aValue)
		}
	}
	return
}

// UnionMap 求并集
func UnionMap(aMap, bMap map[interface{}]interface{}) (retSet []interface{}) {
	for _, aValue := range aMap {
		retSet = append(retSet, aValue)
	}
	for bKey, bValue := range bMap {
		if _, ok := aMap[bKey]; !ok {
			retSet = append(retSet, bValue)
		}
	}
	return
}

// SubMap 求差集
func SubMap(aMap, bMap map[interface{}]interface{}) (retSet []interface{}) {
	for aKey, aValue := range aMap {
		if _, ok := bMap[aKey]; !ok {
			retSet = append(retSet, aValue)
		}
	}
	return
}

// CompareMap 相同key，value不相同
func CompareMap(aMap, bMap map[interface{}]interface{}) (a, b []interface{}) {
	for aKey, aValue := range aMap {
		if bValue, ok := bMap[aKey]; ok {
			if !reflect.DeepEqual(bValue, aValue) {
				a = append(a, aValue)
				b = append(b, bValue)
			}
		}
	}
	return
}

// DiffMap 求变动
func DiffMap(oldMap, curMap interface{}) (dels, adds, altOld, altCur []interface{}) {
	old := reflect.ValueOf(oldMap)
	if old.Kind() != reflect.Map {
		panic("toslice arr not map")
	}
	cur := reflect.ValueOf(curMap)
	if cur.Kind() != reflect.Map {
		panic("toslice arr not map")
	}
	// old - cur
	oldKeys := old.MapKeys()
	for _, k := range oldKeys {
		v := cur.MapIndex(k)
		if v.Kind() == reflect.Invalid {
			dels = append(dels, old.MapIndex(k).Interface())
		}
	}
	// cur - old
	curKeys := cur.MapKeys()
	for _, k := range curKeys {
		v := old.MapIndex(k)
		if v.Kind() == reflect.Invalid {
			adds = append(adds, cur.MapIndex(k).Interface())
		}
	}
	// compare
	for _, k := range oldKeys {
		v := cur.MapIndex(k)
		if v.Kind() != reflect.Invalid {
			curV := cur.MapIndex(k).Interface()
			oldV := old.MapIndex(k).Interface()
			if !reflect.DeepEqual(curV, oldV) {
				altOld = append(altOld, oldV)
				altCur = append(altCur, curV)
			}
		}
	}
	return dels, adds, altOld, altCur
}
