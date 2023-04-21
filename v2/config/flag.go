package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	flag "github.com/cornfeedhobo/pflag"
	"gopkg.in/yaml.v3"
)

const DefaultRegistry string = "huray-nks-container-registry"

type Flag interface {
	Setup(*[]string) *flag.FlagSet
}

type GetFlag struct {
	Registry  string
	Image     string
	Output    string
	NoHeaders bool
}

type DeleteConfig struct {
	Registry string   `yaml:"registry"`
	Image    string   `yaml:"image"`
	Tags     []string `yaml:"tags"`
	Enable   bool     `yaml:"all"`
	Recent   int      `yaml:"exclude-recent"`
}
type DeleteFlag struct {
	DeleteConfig
	FileName string
	Yes      bool
}

func (f *GetFlag) Setup(cmd *[]string) *flag.FlagSet {
	flags := flag.NewFlagSet("flag", flag.ExitOnError)
	if (*cmd)[1] != "registry" {
		if (*cmd)[1] != "image" {
			flags.StringVarP(&f.Image, "image", "i", "", "image name")
		}
		flags.StringVarP(&f.Registry, "registry", "r", DefaultRegistry, "registry name")
	}
	flags.StringVarP(&f.Output, "output", "o", "", "available output format:json,yaml")
	flags.BoolVar(&f.NoHeaders, "no-headers", false, "When using the default or custom-column output format, don't print headers (default print headers).")
	return flags
}

func (f *DeleteFlag) Setup(cmd *[]string) *flag.FlagSet {
	flags := flag.NewFlagSet("flag", flag.ExitOnError)
	if (*cmd)[1] == "tag" {
		flags.StringVarP(&f.FileName, "filename", "f", "", "The files that contain the configurations to apply")
		flags.StringVarP(&f.Image, "image", "i", "", "")
		flags.BoolVar(&f.Enable, "all", false, "")
		flags.IntVar(&f.Recent, "exclude-recent", 0, "")
	}
	flags.StringVarP(&f.Registry, "registry", "r", DefaultRegistry, "registry name")
	flags.BoolVarP(&f.Yes, "yes", "y", false, "Delete "+(*cmd)[1]+" without asking.")
	return flags
}

func FlagParse(flags *flag.FlagSet, args []string) {
	flags.Parse(args)
}

func (f *DeleteFlag) GetConfigFromFile() (*[]DeleteConfig, error) {
	var configs []DeleteConfig
	var err error

	filePath, err := filepath.Abs(f.FileName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &configs, nil
}
