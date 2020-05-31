package acllibgo

import (
	"reflect"
	"strings"
	"sync"
)

const tagName string = "acl"

var _typeCache map[string]typeInfo
var _cacheSync *sync.Mutex

type typeInfo struct {
	Name          string
	PrkPath       string
	ToStringValue string
	Kind          reflect.Kind
	Field         []fieldInfo
}

type fieldInfo struct {
	Name    string
	Kind    reflect.Kind
	AclTags []string
}

func init() {
	_cacheSync = new(sync.Mutex)
	_typeCache = make(map[string]typeInfo)
}

func getTypeInfo(itemType reflect.Type) typeInfo {
	if itemType == nil {
		return typeInfo{}
	}

	// type key
	k := itemType.PkgPath() + ":" + itemType.Name()

	// find in cache - thread-safe
	// we lock entire call to prevent run-on by multiple CPUs
	_cacheSync.Lock()
	defer _cacheSync.Unlock()

	// if found in cache, return copy of cached data
	if item, ok := _typeCache[k]; ok {
		return item
	}

	// reflect type information - this is slow (relative to computer world)
	rv := typeInfo{}
	rv.Name = itemType.Name()
	rv.PrkPath = itemType.PkgPath()
	rv.ToStringValue = itemType.String()
	rv.Kind = itemType.Kind()

	if rv.Kind == reflect.Struct {
		rv.Field = make([]fieldInfo, itemType.NumField())

		for x := 0; x < len(rv.Field); x++ {
			field := itemType.Field(x)
			delta := fieldInfo{}
			delta.Name = field.Name
			delta.Kind = field.Type.Kind()

			aclTag := strings.TrimSpace(field.Tag.Get(tagName))

			if len(aclTag) > 0 {
				delta.AclTags = strings.Split(aclTag, ",")
			} else {
				delta.AclTags = []string{}
			}

			rv.Field[x] = delta
		}
	} else {
		rv.Field = []fieldInfo{}
	}

	// save in cache map
	_typeCache[k] = rv

	return rv
}
