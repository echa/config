// Copyright (c) 2013-2021 KIDTSUNAMI
// Author: alex@kidtsunami.com

package config

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var stringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

func toString(t interface{}) string {
	val := reflect.Indirect(reflect.ValueOf(t))
	if !val.IsValid() {
		return ""
	}
	if val.Type().Implements(stringerType) {
		return t.(fmt.Stringer).String()
	}
	if s, err := toRawString(val.Interface()); err == nil {
		return s
	}
	return fmt.Sprintf("%v", val.Interface())
}

func toRawString(t interface{}) (string, error) {
	val := reflect.Indirect(reflect.ValueOf(t))
	if !val.IsValid() {
		return "", nil
	}
	typ := val.Type()
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, val.Type().Bits()), nil
	case reflect.String:
		return val.String(), nil
	case reflect.Bool:
		return strconv.FormatBool(val.Bool()), nil
	case reflect.Array:
		if typ.Elem().Kind() != reflect.Uint8 {
			b := strings.Builder{}
			for i, l := 0, val.Len(); i < l; i++ {
				v, err := toRawString(val.Index(i).Interface())
				if err != nil {
					return "", err
				}
				if b.Len() > 0 {
					b.WriteByte(',')
				}
				b.WriteString(v)
			}
			return b.String(), nil
		}
		// [...]byte
		var b []byte
		if val.CanAddr() {
			b = val.Slice(0, val.Len()).Bytes()
		} else {
			b = make([]byte, val.Len())
			reflect.Copy(reflect.ValueOf(b), val)
		}
		return hex.EncodeToString(b), nil
	case reflect.Slice:
		if typ.Elem().Kind() != reflect.Uint8 {
			b := strings.Builder{}
			for i, l := 0, val.Len(); i < l; i++ {
				v, err := toRawString(val.Index(i).Interface())
				if err != nil {
					return "", err
				}
				if b.Len() > 0 {
					b.WriteByte(',')
				}
				b.WriteString(v)
			}
			return b.String(), nil
		}
		// []byte
		b := val.Bytes()
		return hex.EncodeToString(b), nil
	case reflect.Map:
		b := strings.Builder{}
		for _, e := range val.MapKeys() {
			k, err := toRawString(e.Interface())
			if err != nil {
				return "", err
			}
			v := val.MapIndex(e)
			vv, err := toRawString(v.Interface())
			if err != nil {
				return "", err
			}
			if b.Len() > 0 {
				b.WriteByte(',')
			}
			b.WriteString(k)
			b.WriteByte('=')
			b.WriteString(vv)
		}
		return b.String(), nil
	}
	return "", fmt.Errorf("no method for converting type %s (%v) to string", typ.String(), val.Kind())
}
