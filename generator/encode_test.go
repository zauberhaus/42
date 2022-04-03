package generator_test

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/zauberhaus/42/cmd"
	"github.com/zauberhaus/42/generator"
)

type TestConfig struct {
	Integer        int
	String         string
	Bool           bool
	Options        TestOptions
	OptionsPointer *TestOptions
}

type TestOptions struct {
	Path string
}

func TestDecode(t *testing.T) {
	dir, err := ioutil.TempDir("", "config")
	if !assert.NoError(t, err) {
		return
	}
	defer os.RemoveAll(dir)

	wanted := TestConfig{
		Integer: 674741956,
		String:  "hdgkjFGAHfkjhakj",
		Bool:    false,
		Options: TestOptions{
			Path: "gsdgdgafgjd",
		},
		OptionsPointer: &TestOptions{
			Path: "ughjkgFA",
		},
	}

	formats := []string{"json", "yaml", "toml", "hcl"}

	for _, f := range formats {
		t.Run(t.Name()+"_"+f, func(t *testing.T) {
			cfg := TestConfig{}

			filename := filepath.Join(dir, "config.json")
			data, err := generator.Marshal(wanted, filename)
			if !assert.NoError(t, err) {
				return
			}

			err = ioutil.WriteFile(filename, data, fs.ModePerm)
			if !assert.NoError(t, err) {
				return
			}

			rootCmd := cmd.NewRootCmd(
				&cobra.Command{Use: t.Name(),
					Short: "Test program",
					Run: func(cmd *cobra.Command, args []string) {
						assert.Equal(t, wanted.Integer, cfg.Integer)
						assert.Equal(t, wanted.String, cfg.String)
						assert.Equal(t, wanted.Bool, cfg.Bool)

						assert.Equal(t, wanted.Options.Path, cfg.Options.Path)

						if assert.NotNil(t, cfg.OptionsPointer) {
							assert.Equal(t, wanted.OptionsPointer.Path, cfg.OptionsPointer.Path)
						}

					},
				}, &cfg,
			)

			rootCmd.SetArgs([]string{"--config", filename})

			err = rootCmd.Execute()
			assert.NoError(t, err)

		})
	}
}
