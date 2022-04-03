package generator

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/spf13/viper"
)

func EnvBindings() map[string][]string {
	f := reflect.ValueOf(viper.GetViper()).Elem().FieldByName("env")
	rf := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	i := rf.Interface()
	return i.(map[string][]string)
}

func GroupedBindings() ([]map[string]string, error) {
	env := EnvBindings()

	groups := make(map[string]map[string]string)
	group := ""
	var items map[string]string
	for k, l := range env {
		if len(l) > 1 {
			return nil, fmt.Errorf("More than one env binding for %v", k)
		}

		parts := strings.Split(k, ".")
		name := fmt.Sprintf("%v.", len(parts)) + parts[0]
		if len(parts) > 2 {
			name += "." + parts[1]
		}

		if name != group {
			if items != nil {
				groups[group] = items
			}

			group = name
			if i, ok := groups[group]; ok {
				items = i
			} else {
				items = make(map[string]string)
			}

		}

		items[k] = l[0]
	}

	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := []map[string]string{}
	for _, k := range keys {
		result = append(result, groups[k])
	}

	return result, nil
}
