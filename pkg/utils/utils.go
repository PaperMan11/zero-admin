package utils

import "reflect"

type TypeError struct {
	FieldName string
}

func (e *TypeError) Error() string {
	return "Type mismatch for field: " + e.FieldName
}

// MapToStruct 将 map 映射到结构体指针
func MapToStruct(data map[string]interface{}, dest interface{}) error {
	destVal := reflect.ValueOf(dest).Elem()
	destType := destVal.Type()

	// 构建字段名到字段信息的映射
	fieldMap := make(map[string]reflect.StructField)
	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		tagName := field.Tag.Get("json")
		if tagName == "" {
			tagName = field.Name
		}
		fieldMap[tagName] = field
	}

	for key, value := range data {
		structField, ok := fieldMap[key]
		if !ok {
			continue // 忽略不存在的字段
		}

		field := destVal.FieldByName(structField.Name)
		if !field.IsValid() {
			continue
		}
		if !field.CanSet() {
			continue
		}

		valueType := reflect.TypeOf(value)
		fieldType := structField.Type

		// 处理指针字段的情况
		if fieldType.Kind() == reflect.Ptr {
			elemType := fieldType.Elem()
			if valueType != elemType {
				// 类型不匹配
				return &TypeError{FieldName: key}
			}
			// 创建指针并赋值
			ptr := reflect.New(elemType)
			ptr.Elem().Set(reflect.ValueOf(value))
			field.Set(ptr)
		} else {
			if valueType != fieldType {
				//fmt.Printf("valueType=%v, fieldType=%v\n", valueType, fieldType)
				return &TypeError{FieldName: key}
			}
			field.Set(reflect.ValueOf(value))
		}
	}
	return nil
}
