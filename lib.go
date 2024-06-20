package typejuggle

import (
	"fmt"
	"reflect"
	"strconv"
)

func FillFields[D interface{}](src interface{}, dest *D, convert ...bool) {
	assignRecursive(reflect.ValueOf(src), reflect.ValueOf(dest), convert...)
}

func assignRecursive(srcVal reflect.Value, destVal reflect.Value, convert ...bool) {
	if srcVal.Kind() == reflect.Ptr {
		srcVal = deepDereference(srcVal)
	}

	if destVal.Kind() == reflect.Ptr {
		destVal = deepAllocate(destVal)
	}

	if srcVal.Kind() == reflect.Struct {
		if destVal.Kind() == reflect.Struct {
			assignStructFieldsRecursive(srcVal, destVal, convert...)
		}
	} else if srcVal.Kind() == reflect.Slice {
		assignSliceFields(srcVal, destVal, convert...)
	} else if srcVal.IsValid() && destVal.IsValid() {
		if srcVal.Type().AssignableTo(destVal.Type()) {
			destVal.Set(srcVal)
		} else if len(convert) > 0 && convert[0] { // WARNING: this can overflow integers
			if convertedVal, err := specialConversion(srcVal, destVal); err == nil {
				destVal.Set(convertedVal)
			} else if srcVal.Type().ConvertibleTo(destVal.Type()) {
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

func deepDereference(ptr reflect.Value) reflect.Value {
	if ptr.Kind() != reflect.Ptr {
		panic("dereferenceRecursive: val must be a pointer")
	}

	for ptr.Kind() == reflect.Ptr {
		ptr = ptr.Elem()
	}

	return ptr
}

func deepAllocate(ptr reflect.Value) reflect.Value {
	if ptr.Kind() != reflect.Ptr {
		panic("allocateRecursive: destPtr must be a pointer")
	}

	for ptr.Kind() == reflect.Ptr {
		if ptr.IsNil() {
			ptr.Set(reflect.New(ptr.Type().Elem()))
		}
		ptr = ptr.Elem()
	}

	return ptr
}

func specialConversion(srcVal reflect.Value, destVal reflect.Value) (reflect.Value, error) {
	switch srcVal.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch destVal.Kind() {
		case reflect.String:
			strVal := strconv.Itoa(int(srcVal.Int()))
			return reflect.ValueOf(strVal), nil
		default:
			return reflect.Value{}, fmt.Errorf("specialConversion: unsupported conversion from %s to %s", srcVal.Kind(), destVal.Kind())
		}
	case reflect.String:
		switch destVal.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.Atoi(srcVal.String()); err == nil {
				return reflect.ValueOf(intVal), nil
			}
			return reflect.Value{}, fmt.Errorf("specialConversion: failed to convert %s to %s", srcVal.Kind(), destVal.Kind())
		case reflect.Bool:
			if boolVal, err := strconv.ParseBool(srcVal.String()); err == nil {
				return reflect.ValueOf(boolVal), nil
			}
			return reflect.Value{}, fmt.Errorf("specialConversion: failed to convert %s to %s", srcVal.Kind(), destVal.Kind())
		default:
			return reflect.Value{}, fmt.Errorf("specialConversion: unsupported conversion from %s to %s", srcVal.Kind(), destVal.Kind())
		}
	case reflect.Bool:
		switch destVal.Kind() {
		case reflect.String:
			strVal := strconv.FormatBool(srcVal.Bool())
			return reflect.ValueOf(strVal), nil
		default:
			return reflect.Value{}, fmt.Errorf("specialConversion: unsupported conversion from %s to %s", srcVal.Kind(), destVal.Kind())
		}
	default:
		return reflect.Value{}, fmt.Errorf("specialConversion: unsupported conversion from %s to %s", srcVal.Kind(), destVal.Kind())
	}
}
