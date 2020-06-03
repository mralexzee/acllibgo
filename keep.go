// Copyright 2020 Alexander Zherdev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package acllibgo

import (
	"errors"
	"reflect"
)

type StructField struct {
	Name   string        `json:"name,omitempty"`
	Fields []StructField `json:"fields,omitempty"`
}

// Keep retains the value of the properties provided, other properties are set to defaults
func Keep(item interface{}, fields []StructField) error {
	if item == nil {
		return errors.New("fields: nil item")
	}
	if fields == nil {
		return errors.New("fields: nil fields")
	}

	itemValue := reflect.ValueOf(item)
	if !itemValue.IsValid() {
		return nil
	}

	// Supports a pointer to a struct, an array of pointers to struct, and map of pointers to struct
	switch itemValue.Kind() {
	case reflect.Slice, reflect.Array:
		arrItemType := reflect.TypeOf(item).Elem()
		if arrItemType.Kind() != reflect.Ptr {
			return errors.New("fields: expecting pointer for slice or array elements")
		}

		for i := 0; i < itemValue.Len(); i++ {
			item := itemValue.Index(i)
			Keep(item.Interface(), fields)
		}

		return nil
	case reflect.Map:
		arrItemType := reflect.TypeOf(item).Elem()
		if arrItemType.Kind() != reflect.Ptr {
			return errors.New("fields: expecting pointer for map values")
		}

		for _, mKey := range itemValue.MapKeys() {
			mapValue := itemValue.MapIndex(mKey)
			if mapValue.Kind() == reflect.Ptr && !mapValue.IsNil() {
				Keep(mapValue.Interface(), fields)
			}
		}

		return nil
	case reflect.Ptr:
		if itemValue.IsNil() {
			return errors.New("fields: nil " + itemValue.Type().String())
		}

	default:
		return errors.New("fields: expecting pointer, slice, or map")
	}

	elemValue := itemValue.Elem()
	if !elemValue.IsValid() {
		return nil
	}

	elemTypeInfo := getTypeInfo(elemValue.Type())

	// Ensure we have a struct
	if elemTypeInfo.Kind != reflect.Struct {
		return errors.New("fields: expecting struct, got " + elemTypeInfo.ToStringValue)
	}

	// Identify which properties to clear
	for i := 0; i < len(elemTypeInfo.Field); i++ {
		itemFieldInfo := elemTypeInfo.Field[i]

		hasDefaultValue := false
		found := false
		fieldFields := []StructField{}
		for _, k := range fields {
			if k.Name == itemFieldInfo.Name || k.Name == "*" {
				found = true
				fieldFields = k.Fields
				break
			}
		}

		if !found {
			hasDefaultValue = true
			ev := elemValue.Field(i)
			setToDefault(ev)
		}

		// scrub field if we have not set it to default and it's a supports type
		if !hasDefaultValue {
			switch itemFieldInfo.Kind {
			case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map:
				ev := elemValue.Field(i)
				Keep(ev.Interface(), fieldFields)
			}
		}
	}

	return nil
}
