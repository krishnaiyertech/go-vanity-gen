// SPDX-FileCopyrightText: Copyright 2025 Krishna Iyer <www.krishnaiyer.tech>
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"reflect"

	"github.com/spf13/pflag"
)

// genFlags reads the configuration and uses reflection to generate a corresponding flagset.
// Takes an input pointer to bind flags directly to the element.
func genFlags[T any](c *T) (*pflag.FlagSet, error) {
	fs := pflag.NewFlagSet("config", pflag.ContinueOnError)

	v := reflect.ValueOf(c)

	// Ensure we have a pointer to a struct
	if v.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("expected pointer to struct, got %s", v.Kind())
	}

	v = v.Elem() // Dereference the pointer

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", v.Kind())
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip un-settable fields
		if !fieldValue.CanSet() {
			continue
		}

		// Get the required tag values
		name := field.Tag.Get("name")
		short := field.Tag.Get("short")
		description := field.Tag.Get("description")

		// Skip fields without name tag
		if name == "" {
			continue
		}

		// Get pointer to the field for *Var methods
		fieldPtr := fieldValue.Addr().Interface()

		switch fieldValue.Kind() {
		case reflect.String:
			if short != "" {
				fs.StringVarP(fieldPtr.(*string), name, short, fieldValue.String(), description)
			} else {
				fs.StringVar(fieldPtr.(*string), name, fieldValue.String(), description)
			}
		case reflect.Int:
			if short != "" {
				fs.IntVarP(fieldPtr.(*int), name, short, int(fieldValue.Int()), description)
			} else {
				fs.IntVar(fieldPtr.(*int), name, int(fieldValue.Int()), description)
			}
		case reflect.Int8:
			if short != "" {
				fs.Int8VarP(fieldPtr.(*int8), name, short, int8(fieldValue.Int()), description)
			} else {
				fs.Int8Var(fieldPtr.(*int8), name, int8(fieldValue.Int()), description)
			}
		case reflect.Int16:
			if short != "" {
				fs.Int16VarP(fieldPtr.(*int16), name, short, int16(fieldValue.Int()), description)
			} else {
				fs.Int16Var(fieldPtr.(*int16), name, int16(fieldValue.Int()), description)
			}
		case reflect.Int32:
			if short != "" {
				fs.Int32VarP(fieldPtr.(*int32), name, short, int32(fieldValue.Int()), description)
			} else {
				fs.Int32Var(fieldPtr.(*int32), name, int32(fieldValue.Int()), description)
			}
		case reflect.Int64:
			if short != "" {
				fs.Int64VarP(fieldPtr.(*int64), name, short, fieldValue.Int(), description)
			} else {
				fs.Int64Var(fieldPtr.(*int64), name, fieldValue.Int(), description)
			}
		case reflect.Uint:
			if short != "" {
				fs.UintVarP(fieldPtr.(*uint), name, short, uint(fieldValue.Uint()), description)
			} else {
				fs.UintVar(fieldPtr.(*uint), name, uint(fieldValue.Uint()), description)
			}
		case reflect.Uint8:
			if short != "" {
				fs.Uint8VarP(fieldPtr.(*uint8), name, short, uint8(fieldValue.Uint()), description)
			} else {
				fs.Uint8Var(fieldPtr.(*uint8), name, uint8(fieldValue.Uint()), description)
			}
		case reflect.Uint16:
			if short != "" {
				fs.Uint16VarP(fieldPtr.(*uint16), name, short, uint16(fieldValue.Uint()), description)
			} else {
				fs.Uint16Var(fieldPtr.(*uint16), name, uint16(fieldValue.Uint()), description)
			}
		case reflect.Uint32:
			if short != "" {
				fs.Uint32VarP(fieldPtr.(*uint32), name, short, uint32(fieldValue.Uint()), description)
			} else {
				fs.Uint32Var(fieldPtr.(*uint32), name, uint32(fieldValue.Uint()), description)
			}
		case reflect.Uint64:
			if short != "" {
				fs.Uint64VarP(fieldPtr.(*uint64), name, short, fieldValue.Uint(), description)
			} else {
				fs.Uint64Var(fieldPtr.(*uint64), name, fieldValue.Uint(), description)
			}
		case reflect.Bool:
			if short != "" {
				fs.BoolVarP(fieldPtr.(*bool), name, short, fieldValue.Bool(), description)
			} else {
				fs.BoolVar(fieldPtr.(*bool), name, fieldValue.Bool(), description)
			}
		case reflect.Float32:
			if short != "" {
				fs.Float32VarP(fieldPtr.(*float32), name, short, float32(fieldValue.Float()), description)
			} else {
				fs.Float32Var(fieldPtr.(*float32), name, float32(fieldValue.Float()), description)
			}
		case reflect.Float64:
			if short != "" {
				fs.Float64VarP(fieldPtr.(*float64), name, short, fieldValue.Float(), description)
			} else {
				fs.Float64Var(fieldPtr.(*float64), name, fieldValue.Float(), description)
			}
		case reflect.Slice:
			if fieldValue.Type().Elem().Kind() == reflect.String {
				defaultValue := make([]string, fieldValue.Len())
				for j := 0; j < fieldValue.Len(); j++ {
					defaultValue[j] = fieldValue.Index(j).String()
				}
				if short != "" {
					fs.StringSliceVarP(fieldPtr.(*[]string), name, short, defaultValue, description)
				} else {
					fs.StringSliceVar(fieldPtr.(*[]string), name, defaultValue, description)
				}
			} else {
				return nil, fmt.Errorf("unsupported slice type %s for field %s", fieldValue.Type(), field.Name)
			}
		default:
			return nil, fmt.Errorf("unsupported field type %s for field %s", fieldValue.Kind(), field.Name)
		}
	}

	return fs, nil
}
