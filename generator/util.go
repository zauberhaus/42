package generator

import (
	"reflect"
	"unsafe"

	"github.com/spf13/viper"
)

func GetEnvBindings() map[string][]string {
	f := reflect.ValueOf(viper.GetViper()).Elem().FieldByName("env")
	rf := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
	i := rf.Interface()
	return i.(map[string][]string)
}
