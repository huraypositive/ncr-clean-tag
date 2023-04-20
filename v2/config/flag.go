package config

import (
	flag "github.com/cornfeedhobo/pflag"
)

type Flag interface{
	Setup(*[]string) *flag.FlagSet
}

type GetFlag struct {
	Registry  string
	Image     string
	Output    string
	NoHeaders bool
}

type DeleteFlag struct {
	Recent		int
	Registry  string
	Image			string
	Yes				bool
	All				bool
}

func (f *GetFlag) Setup(cmd *[]string) *flag.FlagSet {
	flags := flag.NewFlagSet("flag", flag.ExitOnError)
	if (*cmd)[1] != "registry" {
		if (*cmd)[1] != "image" {
			flags.StringVarP(&f.Image, "image", "i", "", "image name")
		}
		flags.StringVarP(&f.Registry, "registry", "r", "huray-nks-container-registry", "registry name")
	}
	flags.StringVarP(&f.Output, "output", "o", "", "available output format:json,yaml")
	flags.BoolVar(&f.NoHeaders, "no-headers", false, "When using the default or custom-column output format, don't print headers (default print headers).")
	return flags
}

func (f *DeleteFlag) Setup(cmd *[]string) *flag.FlagSet {
	flags := flag.NewFlagSet("flag", flag.ExitOnError)
	if (*cmd)[1] == "tag" {
		flags.StringVarP(&f.Image, "image", "i", "", "")
		flags.BoolVar(&f.All, "all", false, "")
		flags.IntVar(&f.Recent, "exclude-recent", 0, "")
	}
	flags.StringVarP(&f.Registry, "registry", "r", "huray-nks-container-registry", "registry name")
	flags.BoolVarP(&f.Yes, "yes", "y", false, "Delete " + (*cmd)[1] + " without asking.")
	return flags
}

func FlagParse(flags *flag.FlagSet, args []string) {
	flags.Parse(args)
}
