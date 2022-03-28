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
	"fmt"
	"os"
	"reflect"

	"github.com/zauberhaus/42/logger"

	"github.com/fsnotify/fsnotify"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AddFunc func(*RootCommand)
type InitFunc func(*RootCommand)

type RootCommand struct {
	cobra.Command
	configFile        string
	defaultConfigFile string
	version           *Version

	config interface{}
}

func NewRootCmd(cmd *cobra.Command) *RootCommand {
	var rootCmd *RootCommand

	rootCmd = &RootCommand{
		Command: *cmd,
	}

	rootCmd.init()

	return rootCmd
}

func (r *RootCommand) SetConfig(config interface{}) {
	if reflect.ValueOf(config).Type().Kind() != reflect.Pointer {
		logger.Fatal("Configuration has to be a pointer")
	}

	r.config = config
}

func (r *RootCommand) SetDefaultConfigFile(configFile string) {
	r.defaultConfigFile = configFile
}

func (r *RootCommand) SetVersion(version *Version) {
	r.version = version
}

func (r *RootCommand) WithInit(init InitFunc) {
	init(r)
}

func (r *RootCommand) WithSubCommands(commands ...AddFunc) {
	for _, c := range commands {
		c(r)
	}
}

func (r *RootCommand) Execute() error {
	return r.Command.Execute()
}

func (r *RootCommand) init() {
	old := r.PersistentPreRunE
	r.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		r.initializeConfig(cmd)

		if old != nil {
			return old(cmd, args)
		}

		return nil
	}

	if r.defaultConfigFile == "" {
		r.defaultConfigFile = r.Command.Name()
	}

	r.PersistentFlags().StringVar(&r.configFile, "config", "", "Config file (default is $HOME/"+r.configFile+".yaml)")
}

func (r *RootCommand) initializeConfig(cmd *cobra.Command) error {
	if r.configFile != "" {
		viper.SetConfigFile(r.configFile)
	} else {
		tmp := os.Getenv("CONFIG")
		if tmp != "" {
			r.configFile = tmp
			viper.SetConfigFile(r.configFile)
		} else {
			home, err := homedir.Dir()
			if err != nil {
				logger.Error("Get homedir: %v", err)
				os.Exit(1)
			}

			viper.AddConfigPath(home)
			viper.SetConfigName(r.defaultConfigFile)
		}
	}

	if err := viper.ReadInConfig(); err == nil {
		logger.Info(fmt.Sprintf("Using config file: %v", viper.ConfigFileUsed()))

		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			logger.Info(fmt.Sprintf("Config file changed: %v", e.Name))
		})
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	err := viper.Unmarshal(r.config)
	return err
}

func (r *RootCommand) GetVersion() *Version {
	return r.version
}

func (r *RootCommand) GetConfig() interface{} {
	return r.config
}
