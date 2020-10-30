package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func main() {
	content := `{
		"aa":1,
		"bb":"2",
		"cc":{
			"dd":4,
			"ee":"5"
		}
	}`
	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(content), &m); err != nil {
		panic(err)
	}
	out := map[string]string{}
	walk(reflect.ValueOf(m), "", out)
	fmt.Println(m)
	fmt.Println(out)
}

func walk(v reflect.Value, path string, out map[string]string) {
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			walk(v.MapIndex(k), fmt.Sprintf("%s/%s", path, k), out)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out[path] = fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out[path] = fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		out[path] = fmt.Sprintf("%f", v.Float())
	case reflect.String:
		out[path] = v.String()
	default:
		out[path] = fmt.Sprintf("%v", v)
	}
}
