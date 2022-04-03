package cmd_test

import (
	"os"
	"testing"

	lookup "github.com/mcuadros/go-lookup"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/zauberhaus/42/cmd"
	"github.com/zauberhaus/42/generator"
)

type TestConfig struct {
	Integer        int          `env:"INTEGER" default:"123456"`
	String         string       `env:"STRING" default:"test"`
	Bool           bool         `env:"BOOL" default:"true"`
	Options        TestOptions  `env:"OPTIONS"`
	OptionsCopy    TestOptions  `env:"OPTIONS"`
	OptionsPointer *TestOptions `env:"OPTIONS_POINTER"`
}

type TestOptions struct {
	Path string `json:"path" env:"PATH"`
}

func TestEnvBind(t *testing.T) {
	cfg := TestConfig{}

	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: t.Name(),
			Short: "Test program",
			Run: func(c *cobra.Command, args []string) {
				bindings := generator.EnvBindings()
				wanted := map[string][]string{
					"integer": {
						"INTEGER",
					},
					"string": {
						"STRING",
					},
					"bool": {
						"BOOL",
					},
					"options.path": {
						"OPTIONS_PATH",
					},
					"optionscopy.path": {
						"OPTIONS_PATH",
					},
					"optionspointer.path": {
						"OPTIONS_POINTER_PATH",
					},
				}

				assert.Equal(t, wanted, bindings)

				value, err := lookup.LookupStringI(cfg, "Integer")
				assert.NoError(t, err)
				assert.Equal(t, 9987, value.Interface())

				value, err = lookup.LookupStringI(cfg, "String")
				assert.NoError(t, err)
				assert.Equal(t, "teststring", value.Interface())

				value, err = lookup.LookupStringI(cfg, "Bool")
				assert.NoError(t, err)
				assert.Equal(t, false, value.Interface())

				value, err = lookup.LookupStringI(cfg, "Options.Path")
				assert.NoError(t, err)
				assert.Equal(t, "/test/path", value.Interface())

				value, err = lookup.LookupStringI(cfg, "OptionsCopy.Path")
				assert.NoError(t, err)
				assert.Equal(t, "/test/path", value.Interface())

				value, err = lookup.LookupStringI(cfg, "OptionsPointer.Path")
				assert.NoError(t, err)
				assert.Equal(t, "/test/path2", value.Interface())
			},
		}, &cfg,
	)

	os.Setenv("INTEGER", "9987")
	os.Setenv("STRING", "teststring")
	os.Setenv("BOOL", "false")
	os.Setenv("OPTIONS_PATH", "/test/path")
	os.Setenv("OPTIONS_POINTER_PATH", "/test/path2")

	err := rootCmd.Execute()
	assert.NoError(t, err)
}

func TestDefaults(t *testing.T) {
	cfg := TestConfig{}

	rootCmd := cmd.NewRootCmd(
		&cobra.Command{Use: t.Name(),
			Short: "Test program",
			Run: func(cmd *cobra.Command, args []string) {
				value, err := lookup.LookupStringI(cfg, "Integer")
				assert.NoError(t, err)
				assert.Equal(t, 123456, value.Interface())

				value, err = lookup.LookupStringI(cfg, "String")
				assert.NoError(t, err)
				assert.Equal(t, "test", value.Interface())

				value, err = lookup.LookupStringI(cfg, "Bool")
				assert.NoError(t, err)
				assert.Equal(t, true, value.Interface())
			},
		}, &cfg,
	)

	err := rootCmd.Execute()
	assert.NoError(t, err)
}
