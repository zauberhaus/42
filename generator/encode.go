package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"

	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/printer"
)

func Marshal(cfg interface{}, file string) ([]byte, error) {

	ext := filepath.Ext(file)[1:]

	switch ext {
	case "yml", "yaml":
		return yaml.Marshal(cfg)
	case "json":
		return json.MarshalIndent(cfg, "", "  ")
	case "toml":
		return toml.Marshal(cfg)
	case "hcl":
		return encode(cfg)
	default:
		return nil, fmt.Errorf("Unknown file extension")
	}
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
