package ami

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func marshalAction(i interface{}, id int) ([]byte, error) {
	b := new(bytes.Buffer)
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Struct {
		return nil, errors.New("can't marshal non-struct types")
	}
	_, err := fmt.Fprintf(b, "Action: %s\r\n", v.Type().Name())
	if err != nil {
		return nil, err
	}
	_, err = fmt.Fprintf(b, "ActionID: %d\r\n", id)
	if err != nil {
		return nil, err
	}
	for i := 0; i < v.Type().NumField(); i++ {
		f := v.Field(i)
		sf := v.Type().Field(i)
		name := sf.Name
		//TODO: should we clean name? Invalid characters?
		value := ""
		valueValid := false
		switch f.Kind() {
		case reflect.String:
			value = f.String()
			if value != "" || hasTag(sf, "empty") {
				valueValid = true
			}
			//TODO: should we clean value? Invalid characters?
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(f.Int(), 10)
			if value != "0" || hasTag(sf, "zero") {
				valueValid = true
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = strconv.FormatUint(f.Uint(), 10)
			if value != "0" || hasTag(sf, "zero") {
				valueValid = true
			}
		default:
			return nil, errors.New("can't marshal field type other than string or (u)int")
		}
		if valueValid {
			_, err = fmt.Fprintf(b, "%s: %s\r\n", name, value)
		}
	}
	return b.Bytes(), nil
}

func hasTag(sf reflect.StructField, tag string) bool {
	for _, t := range strings.Split(sf.Tag.Get("ami"), ",") {
		if t == tag {
			return true
		}
	}
	return false
}
