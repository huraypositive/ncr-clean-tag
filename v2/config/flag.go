package config

import (
	flag "github.com/cornfeedhobo/pflag"
)

type Flag struct {
	Registry	string
	Image			string
}

func (f *Flag)Setup() *flag.FlagSet {
	flags := flag.NewFlagSet("flag", flag.ExitOnError)
	flags.StringVarP(&f.Registry,"registry","r","huray-nks-container-registry","registry name")
	flags.StringVarP(&f.Image,"image","i","","image name")
	return flags
}

func (f *Flag)Parse(flags *flag.FlagSet,args []string) {
	flags.Parse(args)
}