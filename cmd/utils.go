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

	BindFlag(target, flag)
}

func BindFlag(target string, flag *pflag.Flag) {
	viper.BindPFlag(target, flag)
	envVar := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(flag.Name, "-", "_"), ".", "_"))
	viper.BindEnv(target, envVar)
}
