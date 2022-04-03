/*
Copyright © 2021 Dirk Lembke <dirk@lembke.nz>

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

package cmd

import (
	"reflect"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zauberhaus/42/logger"
)

func BindCmdFlag(flags *pflag.FlagSet, names ...string) {
	if len(names) == 0 {
		logger.Error("No source or target")
		return
	}

	source := names[0]
	target := source

	if len(names) > 1 {
		target = names[1]
	}

	flag := flags.Lookup(source)
	if flag == nil {
		logger.Errorf("Flag not found: %v", source)
		return
	}

	viper.BindPFlag(target, flag)
}

func AutoBindEnv(config interface{}) {
	parseTags(viper.GetViper(), reflect.ValueOf(config).Type(), []string{}, []string{})
}

func parseTags(viper *viper.Viper, fieldType reflect.Type, path []string, envpath []string) {
	if fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()
	}

	for i := 0; i < fieldType.NumField(); i++ {
		f := fieldType.Field(i)
		tag := f.Tag.Get("env")
		t := f.Type

		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}

		switch t.Kind() {
		case reflect.Struct:
			subEnvPath := envpath
			if tag != "" {
				subEnvPath = append(envpath, strings.ToUpper(tag))
			}

			subPath := append(path, f.Name)
			parseTags(viper, t, subPath, subEnvPath)
		default:
			tmp := append(path, f.Name)
			name := strings.Join(tmp, ".")

			if tag != "" && tag != "skip" {
				tmp = append(envpath, strings.ToUpper(tag))
				tag := strings.Join(tmp, "_")

				viper.BindEnv(name, tag)
			} else {
				envVar := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(name, "-", "_"), ".", "_"))
				viper.BindEnv(name, envVar)
			}
		}
	}
}
