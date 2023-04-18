package config

import (
	flag "github.com/cornfeedhobo/pflag"
)

type Flag struct {
	Registry  string
	Image     string
	Output    string
	NoHeaders bool
}

func (f *Flag) Setup(cmd *[]string) *flag.FlagSet {
	flags := flag.NewFlagSet("flag", flag.ExitOnError)
	switch (*cmd)[0] {
	case "get":
		if (*cmd)[1] != "registry" {
			if (*cmd)[1] != "image" {
				flags.StringVarP(&f.Image, "image", "i", "", "image name")
			}
			flags.StringVarP(&f.Registry, "registry", "r", "huray-nks-container-registry", "registry name")
		}
		flags.StringVarP(&f.Output, "output", "o", "", "available output format:json,yaml")
		flags.BoolVar(&f.NoHeaders, "no-headers", false, "When using the default or custom-column output format, don't print headers (default print headers).")
	case "delete":
	}
	return flags
}

func (f *Flag) Parse(flags *flag.FlagSet, args []string) {
	flags.Parse(args)
}
