/*
Copyright 2022 The Katanomi Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"context"
	"errors"
	"reflect"

	"github.com/spf13/cobra"

	"github.com/spf13/pflag"
)

// FlagRegister flag register interface
type FlagRegister interface {
	AddFlags(flags *pflag.FlagSet)
}

// SetupRegister setup register interface
type SetupRegister interface {
	Setup(ctx context.Context, cmd *cobra.Command, args []string) (err error)
}

var (
	registerTypeError = errors.New("register must be a pointer to a struct")
)

// RegisterFlags register flags
func RegisterFlags(obj interface{}, flags *pflag.FlagSet) {
	v := reflect.ValueOf(obj)
	if !isStructPtr(v) {
		return
	}

	elem := v.Elem()
	t := elem.Type()
	for i := 0; i < t.NumField(); i++ {
		if isPtr(elem.Field(i)) {
			continue
		}
		p := elem.Field(i).Addr()
		if p.Type().Implements(reflect.TypeOf((*FlagRegister)(nil)).Elem()) {
			p.Interface().(FlagRegister).AddFlags(flags)
		}
	}
}

// RegisterSetup register setup function
func RegisterSetup(obj interface{}, ctx context.Context, cmd *cobra.Command, args []string) (err error) {
	v := reflect.ValueOf(obj)
	if !isStructPtr(v) {
		return registerTypeError
	}

	elem := v.Elem()
	t := elem.Type()
	for i := 0; i < t.NumField(); i++ {
		if isPtr(elem.Field(i)) {
			continue
		}
		p := elem.Field(i).Addr()
		if !p.Type().Implements(reflect.TypeOf((*SetupRegister)(nil)).Elem()) {
			continue
		}
		if err = p.Interface().(SetupRegister).Setup(ctx, cmd, args); err != nil {
			return err
		}
	}

	return nil
}

func isPtr(v reflect.Value) bool {
	return v.Type().Kind() == reflect.Pointer
}

func isStruct(v reflect.Value) bool {
	return v.Type().Kind() == reflect.Struct
}

func isStructPtr(v reflect.Value) bool {
	return isPtr(v) && isStruct(v.Elem())
}
