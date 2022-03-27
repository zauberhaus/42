package cmd

import (
	"fmt"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Version struct {
	BuildDate    string `yaml:"buildDate,omitempty" json:"buildDate,omitempty" xml:"buildDate,omitempty"`
	Compiler     string `yaml:"compiler" json:"compiler" xml:"compiler"`
	GitCommit    string `yaml:"gitCommit,omitempty" json:"gitCommit,omitempty" xml:"gitCommit,omitempty"`
	GitTreeState string `yaml:"gitTreeState,omitempty" json:"gitTreeState,omitempty" xml:"gitTreeState,omitempty"`
	GitVersion   string `yaml:"gitVersion,omitempty" json:"gitVersion,omitempty" xml:"gitVersion,omitempty"`
	GoVersion    string `yaml:"goVersion" json:"goVersion" xml:"goVersion"`
	Platform     string `yaml:"platform" json:"platform" xml:"platform"`
}

func (v *Version) String() string {
	data, _ := yaml.Marshal(v)
	return string(data)
}

// NewVersion creates a new version object
func NewVersion(buildDate string, gitCommit string, tag string, treeState string) *Version {
	return &Version{
		BuildDate:    buildDate,
		Compiler:     runtime.Compiler,
		GitCommit:    gitCommit,
		GitTreeState: treeState,
		GitVersion:   tag,
		GoVersion:    runtime.Version(),
		Platform:     fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH),
	}
}
