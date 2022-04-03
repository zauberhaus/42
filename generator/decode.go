package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"github.com/zauberhaus/42/logger"
	"gopkg.in/yaml.v3"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
)

func Decode(cfg interface{}, file string) error {

	var data []byte
	var err error

	ext := filepath.Ext(file)[1:]

	switch ext {
	case "yml", "yaml":
		data, err = yaml.Marshal(cfg)
		if err != nil {
			return err
		}
	case "json":
		data, err = json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return err
		}
	case "toml":
		data, err = toml.Marshal(cfg)
		if err != nil {
			return err
		}
	case "hcl":
		data, err = encode(cfg)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown file extension")
	}

	logger.Infof("Write file: %v", file)
	return ioutil.WriteFile(file, data, fs.ModePerm)
}

func encode(v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	// TODO: use printer.Format? Is the trailing newline an issue?

	ast, err := hcl.Parse(string(b))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	err = printer.Fprint(&buf, ast.Node)
	if err != nil {
		return nil, err
	}

	data, err := printer.Format(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return data, nil
}
