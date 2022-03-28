package cmd_test

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/zauberhaus/42/cmd"
)

var (
	config Config
)

type Config struct {
	Name  string
	Value int
}

func TestRunConfigFile(t *testing.T) {
	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: "root-test",
			Short: "A tool to watch and export Kubernetes logs to sentry",
			Run: func(cmd *cobra.Command, args []string) {
				assert.Equal(t, "test", config.Name)
				assert.Equal(t, 1, config.Value)
			},
		},
	)

	rootCmd.SetConfig(&config)
	rootCmd.SetArgs([]string{"--config", "./testdata/config.yaml"})

	err := rootCmd.Execute()
	assert.NoError(t, err)
}

func TestRunFlags(t *testing.T) {
	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: "root-test",
			Short: "A tool to watch and export Kubernetes logs to sentry",
			Run: func(cmd *cobra.Command, args []string) {
				assert.Equal(t, t.Name(), config.Name)
				assert.Equal(t, 99, config.Value)
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
	rootCmd.SetArgs([]string{"-n", t.Name(), "-v", "99"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
}

func TestRunEnv(t *testing.T) {
	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: "root-test",
			Short: "A tool to watch and export Kubernetes logs to sentry",
			Run: func(cmd *cobra.Command, args []string) {
				assert.Equal(t, t.Name(), config.Name)
				assert.Equal(t, 7365, config.Value)
			},
		},
	)

	rootCmd.WithInit(func(rc *cmd.RootCommand) {
		rc.Flags().StringP("name", "n", "", "Name")
		cmd.BindCmdFlag(rc.Flags(), "name")

		rc.Flags().IntP("value", "v", 0, "Value")
		cmd.BindCmdFlag(rc.Flags(), "value")
	})

	os.Setenv("NAME", t.Name())
	os.Setenv("VALUE", "7365")

	rootCmd.SetConfig(&config)
	err := rootCmd.Execute()
	assert.NoError(t, err)

	os.Unsetenv("NAME")
	os.Unsetenv("VALUE")
}
