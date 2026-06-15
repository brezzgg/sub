package yamlreflect

import (
	"fmt"
	"reflect"
	"strings"
)

func ConvertToYAMLStruct(src any) reflect.Type {
	t := reflect.TypeOf(src)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fields := make([]reflect.StructField, 0, t.NumField())

	for i := range t.NumField() {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}
		name := strings.SplitN(jsonTag, ",", 2)[0]
		fields = append(fields, reflect.StructField{
			Name: field.Name,
			Type: field.Type,
			Tag:  reflect.StructTag(fmt.Sprintf(`yaml:"%s"`, name)),
		})
	}

	return reflect.StructOf(fields)
}

func CopyToYAMLStruct(src any) any {
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	yamlType := ConvertToYAMLStruct(srcVal.Interface())
	dst := reflect.New(yamlType).Elem()

	for i := range yamlType.NumField() {
		name := yamlType.Field(i).Name
		dst.Field(i).Set(srcVal.FieldByName(name))
	}

	return dst.Addr().Interface()
}

func CopyFromYAMLStruct(dst any, src any) {
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() == reflect.Ptr {
		dstVal = dstVal.Elem()
	}

	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	for i := range srcVal.Type().NumField() {
		name := srcVal.Type().Field(i).Name
		dstField := dstVal.FieldByName(name)
		if !dstField.IsValid() || !dstField.CanSet() {
			continue
		}
		dstField.Set(srcVal.Field(i))
	}
}
