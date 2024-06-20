package typejuggle

import (
	"reflect"
)

func FillFields[D interface{}](src interface{}, dest *D, convert ...bool) {
	assignRecursive(reflect.ValueOf(src), reflect.ValueOf(dest), convert...)
}

func assignRecursive(srcVal reflect.Value, destVal reflect.Value, convert ...bool) {
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	if destVal.Kind() == reflect.Ptr {
		if destVal.IsNil() {
			destVal.Set(reflect.New(destVal.Type().Elem()))
		}

		destVal = destVal.Elem()
	}

	if srcVal.Kind() == reflect.Struct {
		if destVal.Kind() == reflect.Struct {
			assignStructFieldsRecursive(srcVal, destVal, convert...)
		}
	} else if srcVal.Kind() == reflect.Slice {
		assignSliceFields(srcVal, destVal, convert...)
	} else {
		if srcVal.IsNil() {
			return
		}
		if srcVal.Type().AssignableTo(destVal.Type()) {
			destVal.Set(srcVal)
		} else if len(convert) > 0 && convert[0] { // WARNING: this can overflow integers
			if srcVal.Type().ConvertibleTo(destVal.Type()) {
				destVal.Set(srcVal.Convert(destVal.Type()))
			}
		}
	}
}

func assignSliceFields(srcVal reflect.Value, destVal reflect.Value, convert ...bool) {
	for i := 0; i < srcVal.Len(); i++ {
		srcElem := srcVal.Index(i)
		var destElem reflect.Value

		if destVal.Type().Elem().Kind() == reflect.Ptr {
			destElem = reflect.New(destVal.Type().Elem().Elem())
		} else {
			destElem = reflect.New(destVal.Type().Elem()).Elem()
		}

		assignRecursive(srcElem, destElem, convert...)
		destVal.Set(reflect.Append(destVal, destElem))
	}
}

func assignStructFieldsRecursive(srcVal reflect.Value, destVal reflect.Value, convert ...bool) {
	for i := 0; i < srcVal.NumField(); i++ {
		srcFieldType := srcVal.Type().Field(i)
		srcFieldVal := srcVal.Field(i)

		// Handling "ID" <-> "Id"
		fieldName := srcFieldType.Name
		switch fieldName {
		case "ID":
			fieldName = "Id"
		case "Id":
			fieldName = "ID"
		}

		// Attempt to find the field in the destination struct
		destFieldVal := destVal.FieldByName(fieldName)
		if !destFieldVal.IsValid() {
			// Try using the original name if the special case didn't match
			destFieldVal = destVal.FieldByName(srcFieldType.Name)
		}

		assignRecursive(srcFieldVal, destFieldVal, convert...)
	}
}
