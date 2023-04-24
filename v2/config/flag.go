package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/cornfeedhobo/pflag"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var DefaultRegistry string = "huray-nks-container-registry"
var excludeTags *string

type Flag interface {
	Setup(*[]string) *flag.FlagSet
	Parse(*flag.FlagSet, []string)
}

type GetFlag struct {
	Registry  string
	Image     string
	Output    string
	NoHeaders bool
}

type DeleteConfig struct {
	Registry      string   `yaml:"registry"`
	Image         string   `yaml:"image"`
	Tags          []string `yaml:"tags"`
	ExcludeTags   []string `yaml:"exclude-tags"`
	ExcludeRecent int      `yaml:"exclude-recent"`
	DryRun        bool     `yaml:"dry-run"`
	Enable        bool     `yaml:"all"`
}
type DeleteFlag struct {
	DeleteConfig
	FileName string
	Yes      bool
}

func init() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
	}
	if os.Getenv("DEFAULT_REGISTRY") != "" {
		DefaultRegistry = os.Getenv("DEFAULT_REGISTRY")
	}
}

func (f *GetFlag) Setup(cmd *[]string) *flag.FlagSet {
	flags := flag.NewFlagSet("flag", flag.ExitOnError)
	if (*cmd)[1] != "registry" {
		if (*cmd)[1] != "image" {
			flags.StringVarP(&f.Image, "image", "i", "", "Image name")
		}
		flags.StringVarP(&f.Registry, "registry", "r", DefaultRegistry, "Registry name")
	}
	flags.StringVarP(&f.Output, "output", "o", "", "Available output format:json,yaml")
	flags.BoolVar(&f.NoHeaders, "no-headers", false, "When using the default or custom-column output format, don't print headers (default print headers).")
	flags.Usage = func() {
		fmt.Printf("%s\n\n", GetUsage)
		flags.PrintDefaults()
	}
	return flags
}

func (f *DeleteFlag) Setup(cmd *[]string) *flag.FlagSet {
	flags := flag.NewFlagSet("flag", flag.ExitOnError)
	if (*cmd)[1] == "tag" {
		flags.StringVarP(&f.FileName, "filename", "f", "", "The files that contain the configurations to apply")
		flags.StringVarP(&f.Image, "image", "i", "", "Image name")
		flags.BoolVar(&f.Enable, "all", false, "Delete all tags")
		flags.IntVar(&f.ExcludeRecent, "exclude-recent", 0, "The number of recent tags to be excluded from deletion, only works when --all is true")
	}
	flags.StringVarP(&f.Registry, "registry", "r", DefaultRegistry, "Registry name")
	flags.BoolVarP(&f.Yes, "yes", "y", false, "Delete "+(*cmd)[1]+" without asking.")
	flags.BoolVar(&f.DryRun, "dry-run", false, "Global option. Execute image deletion dry-run.")
	excludeTags = flags.String("exclude-tags", "", excludeTagsUsage)
	flags.Usage = func() {
		fmt.Printf("%s\n\n", DeleteUsage)
		flags.PrintDefaults()
	}
	return flags
}

func (f *GetFlag) Parse(flags *flag.FlagSet, args []string) {
	flags.Parse(args)
}
func (f *DeleteFlag) Parse(flags *flag.FlagSet, args []string) {
	flags.Parse(args)
	if *excludeTags != "" {
		f.ExcludeTags = strings.Split(*excludeTags, ",")
	}
}

func (f *DeleteFlag) GetConfigFromFile() (*[]DeleteConfig, error) {
	var configs []DeleteConfig
	var err error

	filePath, err := filepath.Abs(f.FileName)
	if err != nil {
		return nil, err
	}
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		return nil, err
	}
	return &configs, nil
}
