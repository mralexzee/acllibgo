// Copyright 2020 Alexander Zherdev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package acllibgo

import (
	"errors"
	"reflect"
	"strings"
)

// Scrub sets structure's fields to default value based on optional 'acl' field tag
// Tag 'acl' on a field has the following effect on Scrub:
//  - <not defined> : Field is not altered
//  - acl:"" : Field is not altered
//  - acl:"*" : Field is not altered as long as Scrub acl has some value
//  - acl:admin : Field is not altered as long as Scrub acl has an array containing "admin" element
//  - acl:admin,user : Field is not altered as long as Scrub acl has an array containing "admin" or "user" element
func Scrub(item interface{}, acl []string) error {
	if item == nil {
		return errors.New("scrub: nil item")
	}
	if acl == nil {
		return errors.New("scrub: nil acl")
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
			return errors.New("scrub: expecting pointer for slice or array elements")
		}

		for i := 0; i < itemValue.Len(); i++ {
			item := itemValue.Index(i)
			Scrub(item.Interface(), acl)
		}

		return nil
	case reflect.Map:
		arrItemType := reflect.TypeOf(item).Elem()
		if arrItemType.Kind() != reflect.Ptr {
			return errors.New("scrub: expecting pointer for map values")
		}

		for _, mKey := range itemValue.MapKeys() {
			mapValue := itemValue.MapIndex(mKey)
			if mapValue.Kind() == reflect.Ptr && !mapValue.IsNil() {
				Scrub(mapValue.Interface(), acl)
			}
		}

		return nil
	case reflect.Ptr:
		if itemValue.IsNil() {
			return errors.New("scrub: nil " + itemValue.Type().String())
		}

	default:
		return errors.New("scrub: expecting pointer, slice, or map")
	}

	elemValue := itemValue.Elem()
	if !elemValue.IsValid() {
		return nil
	}

	elemTypeInfo := getTypeInfo(elemValue.Type())

	// Ensure we have a struct
	if elemTypeInfo.Kind != reflect.Struct {
		return errors.New("scrub: expecting struct, got " + elemTypeInfo.ToStringValue)
	}

	// Identify which properties to clear
	for i := 0; i < len(elemTypeInfo.Field); i++ {

		itemFieldInfo := elemTypeInfo.Field[i]
		hasDefaultValue := false
		if len(itemFieldInfo.AclTags) > 0 {
			found := false

			for _, providedAcl := range acl {
				for _, tagAcl := range itemFieldInfo.AclTags {
					if tagAcl == "*" || strings.EqualFold(tagAcl, providedAcl) {
						found = true
					}
				}
			}

			if !found {
				hasDefaultValue = true
				ev := elemValue.Field(i)
				setToDefault(ev)
			}
		}

		// scrub field if we have not set it to default and it's a supports type
		if !hasDefaultValue {
			switch itemFieldInfo.Kind {
			case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map:
				ev := elemValue.Field(i)
				Scrub(ev.Interface(), acl)
			}
		}
	}

	return nil
}
