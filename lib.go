package typejuggle

import (
	"fmt"
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

		fmt.Println(fieldName)

		// Attempt to find the field in the destination struct
		destFieldVal := destVal.FieldByName(fieldName)
		if !destFieldVal.IsValid() {
			// Try using the original name if the special case didn't match
			destFieldVal = destVal.FieldByName(srcFieldType.Name)
		}

		assignRecursive(srcFieldVal, destFieldVal, convert...)

		// if destFieldVal.IsValid() && destFieldVal.CanSet() {
		// 	if srcFieldVal.Kind() == reflect.Ptr && srcFieldVal.Elem().Kind() == reflect.Struct {
		// 		if destFieldVal.Kind() == reflect.Ptr && destFieldVal.IsNil() {
		// 			// Allocate memory for the destination field's type if its pointer is nil
		// 			destFieldVal.Set(reflect.New(destFieldVal.Type().Elem()))

		// 			// Skip if the destination field is not a struct but the source is
		// 			if destFieldVal.Elem().Kind() != reflect.Struct {
		// 				continue
		// 			}
		// 		}

		// 		assignStructFieldsRecursive(srcFieldVal.Elem(), destFieldVal.Elem())
		// 	} else if srcFieldVal.Kind() == reflect.Struct && destFieldVal.Kind() == reflect.Ptr && destFieldVal.Type().Elem().Kind() == reflect.Struct {
		// 		if destFieldVal.IsNil() {
		// 			destFieldVal.Set(reflect.New(destFieldVal.Type().Elem()))
		// 		}
		// 		assignStructFieldsRecursive(srcFieldVal, destFieldVal.Elem())
		// 	} else if srcFieldVal.Kind() == reflect.Struct && destFieldVal.Kind() == reflect.Struct {
		// 		assignStructFieldsRecursive(srcFieldVal, destFieldVal)
		// 	} else if srcFieldVal.Kind() == reflect.Slice && destFieldVal.Kind() == reflect.Slice {
		// 		assignSliceFields(srcFieldVal, destFieldVal)
		// 	} else if srcFieldVal.Type().AssignableTo(destFieldVal.Type()) {
		// 		destFieldVal.Set(srcFieldVal)
		// 	} else if srcFieldVal.Type().ConvertibleTo(destFieldVal.Type()) {
		// 		destFieldVal.Set(srcFieldVal.Convert(destFieldVal.Type()))
		// 	} else if destFieldVal.Kind() == reflect.Ptr {
		// 		if destFieldVal.IsNil() {
		// 			destFieldVal.Set(reflect.New(destFieldVal.Type().Elem()))
		// 		}

		// 		if srcFieldVal.Type().AssignableTo(destFieldVal.Type().Elem()) {
		// 			destFieldVal.Elem().Set(srcFieldVal)
		// 		} else if srcFieldVal.Type().ConvertibleTo(destFieldVal.Type().Elem()) {
		// 			destFieldVal.Elem().Set(srcFieldVal.Convert(destFieldVal.Type().Elem()))
		// 		}

		// 		if srcFieldVal.Kind() == reflect.Ptr {
		// 			if srcFieldVal.Elem().Type().AssignableTo(destFieldVal.Type().Elem()) {
		// 				destFieldVal.Elem().Set(srcFieldVal.Elem())
		// 			} else if srcFieldVal.Elem().Type().ConvertibleTo(destFieldVal.Type().Elem()) {
		// 				destFieldVal.Elem().Set(srcFieldVal.Elem().Convert(destFieldVal.Type().Elem()))
		// 			}
		// 		}
		// 	}
		// }
	}
}

func main() {
	type Inner struct {
		Name string
		ID   int
	}

	type InnerLow struct {
		Name string
		Id   int
	}

	type A struct {
		Inner  Inner
		Detail string
	}

	type B struct {
		Inner  InnerLow
		Detail string
	}

	a := B{Inner: InnerLow{Name: "Example", Id: 42}, Detail: "Detail A"}
	b := A{}

	aSlice := []B{a, a}
	var bSlice []*A

	fmt.Printf("Before assignment: %+v\n", b)
	FillFields(a, &b)
	fmt.Printf("After assignment: %+v\n", b)

	fmt.Printf("Before slice assignment: %+v\n", bSlice)
	FillFields(aSlice, &bSlice)
	fmt.Printf("After slice assignment: %+v\n", bSlice)

	fmt.Printf("Assignment: %+v\n", *bSlice[0])
	aSlice[0].Inner.Name = "Changed"
	bSlice[0].Inner.Name = "WOOP!"

	fmt.Printf("Assignment: %+v\n", aSlice[0])
	fmt.Printf("Assignment: %+v\n", *bSlice[0])
}
