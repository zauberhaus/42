package cmd

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func (r *RootCommand) AutoBindEnv(config interface{}) {
	r.ParseTags(viper.GetViper(), reflect.ValueOf(config), []string{}, []string{})
}

func (r *RootCommand) ParseTags(viper *viper.Viper, value reflect.Value, path []string, envpath []string) {
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	for i := 0; i < value.NumField(); i++ {
		v := value.Field(i)
		f := value.Type().Field(i)
		tag := f.Tag.Get("env")
		t := f.Type

		switch t.Kind() {
		case reflect.Struct:
			subEnvPath := envpath
			if tag != "" {
				subEnvPath = append(envpath, strings.ToUpper(tag))
			}

			subPath := append(path, f.Name)

			r.ParseTags(viper, v, subPath, subEnvPath)
		default:
			if tag != "" {
				tmp := append(path, f.Name)
				name := strings.Join(tmp, ".")

				tmp = append(envpath, strings.ToUpper(tag))
				tag := strings.Join(tmp, "_")

				viper.BindEnv(name, tag)
			}
		}
	}
}
