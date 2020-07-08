package basemysql

// //Fields retorna []string campos da struct com tag field
// func Fields(values interface{}) []string {
// 	v := reflect.ValueOf(values)
// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}
// 	fields := []string{}
// 	if v.Kind() == reflect.Struct {
// 		for i := 0; i < v.NumField(); i++ {
// 			field := v.Type().Field(i).Tag.Get("field")
// 			if field != "" {
// 				fields = append(fields, field)
// 			}
// 		}
// 		return fields
// 	}
// 	panic(fmt.Errorf("DBFields requires a struct or a map, found: %s", v.Kind().String()))
// }

//Values retorna []interface{} valores
// func Values(values interface{}) []interface{} {
// 	v := reflect.ValueOf(values)
// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}
// 	var fields []interface{}
// 	if v.Kind() == reflect.Struct {
// 		for i := 0; i < v.NumField(); i++ {
// 			v.Type().Field(i).Tag
// 			field := v.Type().Field(i).Tag.Get("field")
// 			if field != "" {
// 				fields = append(fields, field)
// 			}
// 		}
// 		return fields
// 	}
// 	panic(fmt.Errorf("DBFields requires a struct or a map, found: %s", v.Kind().String()))
// }
