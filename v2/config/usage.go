package config

import (
	"fmt"
	"strings"
)

var command string
var Usage string
var GetUsage string
var DeleteUsage string

func init() {
	if command == "" {
		command = "nct"
	}
	Usage = fmt.Sprintf(`%s - NCR Clean Tag
A tool for using the NCloud Container Registry API

Version:
  %s version

Get:
  %s get [registry|images|tags]

Delete:
  %s delete [image|tags]`, strings.ToUpper(command), command, command, command)

	GetUsage = fmt.Sprintf(`Get Registry or Image or Tag list from NCR

Examples:
  # Get registry list
  %s get registry

  # Get image list
  %s get image [-r registry] [-o json or yaml] [--no-headers]
  %s get images [-r registry] [-o json or yaml] [--no-headers]

  # Get image detail
  %s get image imageName [-r registry] [-o json or yaml] [--no-headers]
	
  # Get tag list
  %s get tag -i imageName [-r registry] [-o json or yaml] [--no-headers]
  %s get tags -i imageName [-r registry] [-o json or yaml] [--no-headers]
	
  # Get tag detail
  %s get tag tagName -i imageName [-r registry] [-o json or yaml] [--no-headers]`, command, command, command, command, command, command, command)

	DeleteUsage = fmt.Sprintf(`Delete Image or Tags from NCR

Examples:
  # Delete image
  %s delete image [-r registry] [--dry-run] [-y]
  
  # Delete tags
  %s delete tag [...tagName] -i imageName [-r registry] [--dry-run] [-y]
  %s delete tags [...tagName] -i imageName [-r registry] [--dry-run] [-y]
  %s delete tags [--all] [--exclude-recent number] [-r registry] [--dry-run] [-y]
  %s delete tags [-f filePath] [--dry-run] [-y]`, command, command, command, command, command)
}
