package ami

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	errMarshalType = errors.New("unsupported Action type")
)

func marshalAction(i interface{}, id ActionID) ([]byte, error) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Struct {
		return nil, errMarshalType
	}
	b := new(bytes.Buffer)
	_, _ = fmt.Fprintf(b, "Action: %s\r\n", v.Type().Name())
	if id != 0 {
		_, _ = fmt.Fprintf(b, "ActionID: %d\r\n", uint64(id))
	}
	for i := 0; i < v.Type().NumField(); i++ {
		f := v.Field(i)
		sf := v.Type().Field(i)
		name := sf.Name
		tags := sf.Tag.Get("ami")
		if f.Kind() == reflect.Array || f.Kind() == reflect.Slice {
			for j := 0; j < f.Len(); j++ {
				value, valid := getValue(f.Index(j), tags)
				if valid {
					_, _ = fmt.Fprintf(b, "%s: %s\r\n", name, value)
				}
			}
		} else {
			value, valid := getValue(f, tags)
			if valid {
				_, _ = fmt.Fprintf(b, "%s: %s\r\n", name, value)
			}
		}
	}
	_, _ = fmt.Fprint(b, "\r\n")
	return b.Bytes(), nil
}

func getValue(v reflect.Value, tags string) (value string, valid bool) {
	switch v.Kind() {
	case reflect.String:
		value = v.String()
		if value != "" || hasTag(tags, "empty") {
			valid = true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = strconv.FormatInt(v.Int(), 10)
		if value != "0" || hasTag(tags, "zero") {
			valid = true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = strconv.FormatUint(v.Uint(), 10)
		if value != "0" || hasTag(tags, "zero") {
			valid = true
		}
	}
	return
}

func hasTag(tags, tag string) bool {
	for _, t := range strings.Split(tags, ",") {
		if t == tag {
			return true
		}
	}
	return false
}
