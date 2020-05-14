package acllibgo

import (
	"errors"
	"reflect"
	"strings"
)

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
			return errors.New("scrub: expecting pointer for elements")
		}

		for i := 0; i < itemValue.Len(); i++ {
			item := itemValue.Index(i)
			Scrub(item.Interface(), acl)
		}

		return nil
	case reflect.Map:
		arrItemType := reflect.TypeOf(item).Elem()
		if arrItemType.Kind() != reflect.Ptr {
			return errors.New("scrub: expecting pointer for elements")
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

	typeOfV := elemValue.Type()

	// Ensure we have a struct
	if typeOfV.Kind() != reflect.Struct {
		return errors.New("scrub: expecting struct, got " + typeOfV.String())
	}

	// Identify which properties to clear
	for i := 0; i < elemValue.NumField(); i++ {
		itemField := typeOfV.Field(i)
		ev := elemValue.Field(i)
		if !ev.IsValid() || !ev.CanSet() {
			continue
		}

		aclTag := strings.TrimSpace(itemField.Tag.Get("acl"))
		if len(aclTag) > 0 {
			found := false

			if aclTag == "*" || strings.Contains(aclTag, "*") {
				found = true
			} else {
				aclTagArray := strings.Split(aclTag, ",")

				for _, providedAcl := range acl {
					for _, tagAcl := range aclTagArray {
						if tagAcl == "" || tagAcl == "*" || tagAcl == providedAcl {
							found = true
						}
					}
				}
			}

			if !found {
				setToDefault(ev)
			}
		}

		Scrub(ev.Interface(), acl)
	}

	return nil
}

func setToDefault(f reflect.Value) {
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
