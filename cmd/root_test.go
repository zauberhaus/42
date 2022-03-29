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

package cmd_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/zauberhaus/42/cmd"
	"github.com/zauberhaus/42/logger"
	"gopkg.in/yaml.v3"
)

var (
	config Config
)

type Config struct {
	Name  string
	Value int
}

func TestRunYamlConfigFile(t *testing.T) {
	expected, err := readConfig("./testdata/config.yaml")
	assert.NoError(t, err)

	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: t.Name(),
			Short: "Test program",
			Run: func(cmd *cobra.Command, args []string) {
				assert.Equal(t, expected.Name, config.Name)
				assert.Equal(t, expected.Value, config.Value)
			},
		},
	)

	rootCmd.SetConfig(&config)
	rootCmd.SetArgs([]string{"--config", "./testdata/config.yaml"})

	err = rootCmd.Execute()
	assert.NoError(t, err)
}

func TestRunTomlConfigFile(t *testing.T) {
	expected, err := readConfig("./testdata/config.yaml")
	assert.NoError(t, err)

	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: t.Name(),
			Short: "Test program",
			Run: func(cmd *cobra.Command, args []string) {
				assert.Equal(t, expected.Name, config.Name)
				assert.Equal(t, expected.Value, config.Value)
			},
		},
	)

	rootCmd.SetConfig(&config)
	rootCmd.SetArgs([]string{"--config", "./testdata/config.toml"})

	err = rootCmd.Execute()
	assert.NoError(t, err)
}

func TestRunJsonConfigFile(t *testing.T) {
	expected, err := readConfig("./testdata/config.yaml")
	assert.NoError(t, err)

	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: t.Name(),
			Short: "Test program",
			Run: func(cmd *cobra.Command, args []string) {
				assert.Equal(t, expected.Name, config.Name)
				assert.Equal(t, expected.Value, config.Value)
			},
		},
	)

	rootCmd.SetConfig(&config)
	rootCmd.SetArgs([]string{"--config", "./testdata/config.json"})

	err = rootCmd.Execute()
	assert.NoError(t, err)
}

func TestRunEnvConfigFile(t *testing.T) {
	expected, err := readConfig("./testdata/config.yaml")
	assert.NoError(t, err)

	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: t.Name(),
			Short: "Test program",
			Run: func(cmd *cobra.Command, args []string) {
				assert.Equal(t, expected.Name, config.Name)
				assert.Equal(t, expected.Value, config.Value)
			},
		},
	)

	rootCmd.SetConfig(&config)
	os.Setenv("CONFIG", "./testdata/config.yaml")

	err = rootCmd.Execute()
	assert.NoError(t, err)

	os.Unsetenv("CONFIG")
}

func TestRunFlags(t *testing.T) {
	name := t.Name()
	value := len(t.Name())

	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: t.Name(),
			Short: "Test program",
			Run: func(cmd *cobra.Command, args []string) {
				assert.Equal(t, name, config.Name)
				assert.Equal(t, value, config.Value)
			},
		},
	)

	rootCmd.WithInit(func(rc *cmd.RootCommand) {
		rc.Flags().StringP("name", "n", "", "Name")
		cmd.BindCmdFlag(rc.Flags(), "name")

		rc.Flags().IntP("value", "v", 0, "Value")
		cmd.BindCmdFlag(rc.Flags(), "value")
	})

	rootCmd.SetConfig(&config)
	rootCmd.SetArgs([]string{"-n", name, "-v", fmt.Sprintf("%v", value)})
	err := rootCmd.Execute()
	assert.NoError(t, err)
}

func TestRunEnv(t *testing.T) {
	name := t.Name()
	value := len(t.Name())

	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: t.Name(),
			Short: "Test program",
			Run: func(cmd *cobra.Command, args []string) {
				assert.Equal(t, name, config.Name)
				assert.Equal(t, value, config.Value)
			},
		},
	)

	rootCmd.WithInit(func(rc *cmd.RootCommand) {
		rc.Flags().StringP("name", "n", "", "Name")
		cmd.BindCmdFlag(rc.Flags(), "name")

		rc.Flags().IntP("value", "v", 0, "Value")
		cmd.BindCmdFlag(rc.Flags(), "value")
	})

	os.Setenv("NAME", name)
	os.Setenv("VALUE", fmt.Sprintf("%v", value))

	rootCmd.SetConfig(&config)
	err := rootCmd.Execute()
	assert.NoError(t, err)

	os.Unsetenv("NAME")
	os.Unsetenv("VALUE")
}

func TestRunSubCommand(t *testing.T) {
	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: t.Name(),
			Short: "Test program",
			Run: func(cmd *cobra.Command, args []string) {
				assert.Equal(t, t.Name(), config.Name)
				assert.Equal(t, 7365, config.Value)
			},
		},
	)

	version := cmd.NewVersion("today", "123456", "v1.1.1", "dirty")
	rootCmd.SetVersion(version)

	rootCmd.WithInit(func(rc *cmd.RootCommand) {
		rc.Flags().StringP("name", "n", "", "Name")
		cmd.BindCmdFlag(rc.Flags(), "name")

		rc.Flags().IntP("value", "v", 0, "Value")
		cmd.BindCmdFlag(rc.Flags(), "value")
	})

	rootCmd.WithSubCommands(func(rc *cmd.RootCommand) {
		var versionCmd = &cobra.Command{
			Use:   "version",
			Short: "Show the version info",
			Run: func(cmd *cobra.Command, args []string) {

				data, err := yaml.Marshal(rc.GetVersion())
				if err != nil {
					logger.Error("Invalid version: %v", err)
				}

				fmt.Fprintln(cmd.OutOrStderr(), string(data))
			},
		}

		rc.AddCommand(versionCmd)
	})

	rootCmd.SetConfig(&config)
	rootCmd.SetArgs([]string{
		"version",
	})

	output := &bytes.Buffer{}
	rootCmd.SetOut(output)
	err := rootCmd.Execute()
	if assert.NoError(t, err) {
		var info cmd.Version
		err := yaml.Unmarshal(output.Bytes(), &info)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, version, &info, "Invalid value")
	}
}

func TestLogLevel(t *testing.T) {
	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: t.Name(),
			Short: "Test program",
			Run: func(cmd *cobra.Command, args []string) {
				assert.Equal(t, int8(-1), logger.GetLogger().GetLevel())
			},
		},
	)

	rootCmd.SetConfig(&config)
	rootCmd.SetArgs([]string{"-l", "debug"})

	err := rootCmd.Execute()
	assert.NoError(t, err)
}

func readConfig(file string) (*Config, error) {
	var result Config

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
