// Copyright 2020 Alexander Zherdev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package acllibgo

import (
	"encoding/json"
	"reflect"
)

func setToDefault(f reflect.Value) {
	if !f.IsValid() || !f.CanSet() {
		return
	}

	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f.SetInt(0)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f.SetUint(0)
	case reflect.Float32, reflect.Float64:
		f.SetFloat(0)
	case reflect.Complex64, reflect.Complex128:
		f.SetComplex(0)
	case reflect.String:
		f.SetString("")
	case reflect.Bool:
		f.SetBool(false)
	case reflect.Ptr, reflect.Array, reflect.Map, reflect.Slice, reflect.Interface:
		f.Set(reflect.Zero(f.Type()))
	}
}

func toJson(i interface{}) string {
	j, _ := json.Marshal(i)
	return string(j)
}
