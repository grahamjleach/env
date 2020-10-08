package env

import (
	"errors"
	"os"
	"reflect"
	"strconv"
)

var (
	errorNonPointer = errors.New("load target must be a pointer")
	errorNonStruct  = errors.New("unmarshal target must be a struct")
)

func Unmarshal(i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		return errorNonPointer
	}

	return unmarshalToStruct(v)
}

func unmarshalToStruct(v reflect.Value) error {
	iv := reflect.Indirect(v)

	if iv.Kind() != reflect.Struct {
		return errorNonStruct
	}

	ivType := iv.Type()

	for i := 0; i < iv.NumField(); i++ {
		field := iv.Field(i)
		key := ivType.Field(i).Tag.Get("env")
		def := ivType.Field(i).Tag.Get("default")

		switch field.Kind() {
		case reflect.Struct, reflect.Ptr:
			if err := unmarshalToStruct(field); err != nil {
				return err
			}
		case reflect.String:
			setString(key, def, field)
		case reflect.Bool:
			setBool(key, def, field)
		case reflect.Int:
			setInt(key, def, field)
		}
	}

	return nil
}

func getEnv(key, def string) string {
	if key != "" {
		if env := os.Getenv(key); env != "" {
			return env
		}
	}

	return def
}

func setString(key, def string, field reflect.Value) {
	if env := getEnv(key, def); env != "" {
		field.Set(reflect.ValueOf(env))
	}
}

func setBool(key, def string, field reflect.Value) {
	if env := getEnv(key, def); env != "" {
		b, err := strconv.ParseBool(env)
		if err != nil {
			panic(err)
		}
		field.Set(reflect.ValueOf(b))
	}
}

func setInt(key, def string, field reflect.Value) {
	if env := getEnv(key, def); env != "" {
		i, err := strconv.ParseInt(env, 10, 64)
		if err != nil {
			panic(err)
		}
		field.Set(reflect.ValueOf(int(i)))
	}
}
