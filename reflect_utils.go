package main

import (
	"fmt"
	"reflect"
)

type I interface{}

type A struct {
	Greeting string
	Message  string
	Pi       float64
}

type B struct {
	Struct    A
	Ptr       *A
	Answer    int
	Map       map[string]string
	StructMap map[string]interface{}
	Slice     []string
}

func create() I {
	// The type C is actually hidden, but reflection allows us to look inside it
	type C struct {
		String string
	}

	return B{
		Struct: A{
			Greeting: "Hello!",
			Message:  "translate this",
			Pi:       3.14,
		},
		Ptr: &A{
			Greeting: "What's up?",
			Message:  "point here",
			Pi:       3.14,
		},
		Map: map[string]string{
			"Test": "translate this as well",
		},
		StructMap: map[string]interface{}{
			"C": C{
				String: "deep",
			},
		},
		Slice: []string{
			"and one more",
		},
		Answer: 42,
	}
}

type User struct {
	name     string `json:"name-field"`
	age      int
	random   interface{}
	nickname []string
	tags     map[string]int
	addr     Addr
	laddr    []Addr
	ptr      *Addr
}

type Addr struct {
	address string
	street  string
	pincode int
}

func main() {
	user := User{"John Doe The Fourth", 20, "Random String", []string{"Jack", "Master"}, map[string]int{"jo": 1, "john": 2}, Addr{"Palo", "142", 96132}, []Addr{Addr{"Alto", "131", 96132}, Addr{"Stanford", "432", 96132}}, &Addr{"University", "445", 96132}}

	//res := Validate(reflect.ValueOf(user).Elem())
	res := PrintDetails(&user)

	fmt.Println(res)
	fmt.Println(GetStructFields(res))
	fmt.Println(GetTag(res, "name"))
	fmt.Println(GetValue(res, "addr"))
	fmt.Println(GetType(res, "addr"))

	str := "Hello world"
	fmt.Println(PrintDetails(&str))

	mapStr := map[string]interface{}{"hello": "hi", "here": "world"}
	res = PrintDetails(&mapStr)
	fmt.Println(res)
	fmt.Println(GetValue(res, "hello"))

	created := create()
	fmt.Println(PrintDetails(&created))
	//fmt.Println("name", res["addr"]["street"])
}

func GetValue(arg interface{}, key string) interface{} {
	mapArg := arg.(map[string]interface{})
	if IsStruct(arg) {
		field := mapArg[key].(map[string]interface{})
		return field["value"]
	} else {
		field := mapArg["value"].(map[interface{}]interface{})
		return field[key]
	}
}

func GetType(arg interface{}, key string) interface{} {
	mapArg := arg.(map[string]interface{})
	if IsStruct(arg) {
		field := mapArg[key].(map[string]interface{})
		return field["type"]
	} else {
		fmt.Println("Argument is not a struct")
		return nil
		//field := mapArg["type"]
		//return field
	}

}
func GetTag(arg interface{}, key string) interface{} {
	mapArg := arg.(map[string]interface{})
	if IsStruct(arg) {
		field := mapArg[key].(map[string]interface{})
		return field["tags"]
	} else {
		fmt.Println("Argument is not a struct")
		return nil
	}
}

func IsStruct(arg interface{}) bool {
	mapArg := arg.(map[string]interface{})
	if _, ok := mapArg["type"]; ok {
		_, ok = mapArg["type"].(map[string]interface{})
		if ok {
			return true
		} else {
			return false
		}
	}
	return true
}

func GetStructFields(arg interface{}) []string {
	if !IsStruct(arg) {
		fmt.Println("Argument is not struct")
	}
	mapArg := arg.(map[string]interface{})
	keys := []string{}
	for k := range mapArg {
		keys = append(keys, k)
	}
	return keys
}

func PrintDetails(arg interface{}) interface{} {

	return ParseFields(reflect.ValueOf(arg).Elem())
}

func ParseFields(uValue reflect.Value) map[string]interface{} {
	res := map[string]interface{}{}
	//innermost := 1
	switch uValue.Kind() {
	case reflect.Struct:
		//innermost = 0
		//res["innermost"] = innermost
		for i := 0; i < uValue.NumField(); i++ {
			field := uValue.Field(i)
			name := uValue.Type().Field(i).Name

			fieldType := map[string]interface{}{}
			fieldType["type"] = field.Kind()
			if len(uValue.Type().Field(i).Tag) != 0 {
				fieldType["tags"] = uValue.Type().Field(i).Tag
			}
			//fieldType["type"] = uValue.Type().Field(i).Type.Name()
			fieldType["value"] = PopulateFieldValues(field)
			res[name] = fieldType

			//fmt.Println("Field Name ", name)
		}
	default:
		res["type"] = uValue.Kind()
		res["value"] = PopulateFieldValues(uValue)
	}
	return res
}

func PopulateFieldValues(field reflect.Value) interface{} {
	switch field.Kind() {
	case reflect.String:
		return field.String()
	case reflect.Bool:
		return field.Bool()
		//fmt.Println("String ", field.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int()
		//fmt.Println("Int ", field.Int())
	case reflect.Float32, reflect.Float64:
		return field.Float()
	case reflect.Map:
		//fmt.Println(field.Type())
		mapType := map[interface{}]interface{}{}
		for _, key := range field.MapKeys() {
			val := field.MapIndex(key)
			mapType[PopulateFieldValues(key)] = PopulateFieldValues(val)
		}
		return mapType
	case reflect.Slice:
		sliceType := []interface{}{}
		for i := 0; i < field.Len(); i += 1 {
			val := field.Index(i)
			sliceType = append(sliceType, PopulateFieldValues(val))
		}
		return sliceType
	case reflect.Struct:
		//fmt.Println(Validate(field))
		//innermost = 0
		return ParseFields(field)
	case reflect.Ptr:
		if !field.IsValid() {
			return nil
		}
		return PopulateFieldValues(field.Elem())
		//return fmt.Sprintf("0x%x", field.Pointer())
	case reflect.Interface:
		return PopulateFieldValues(field.Elem())

	}
	return nil
}
