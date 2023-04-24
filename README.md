# NCT - NCR Clean Tag
A tool for using the NCloud Container Registry API

## Version
```sh
nct version
```

## Get
```sh
# Get Registry or Image or Tag list from NCR
nct get [registry|images|tags]

# Example:
  # Get registry list
  nct get registry

  # Get image list
  nct get image [-r registry] [-o json or yaml] [--no-headers]
  nct get images [-r registry] [-o json or yaml] [--no-headers]

  # Get image detail
  nct get image imageName [-r registry] [-o json or yaml] [--no-headers]

  # Get tag list
  nct get tag -i imageName [-r registry] [-o json or yaml] [--no-headers]
  nct get tags -i imageName [-r registry] [-o json or yaml] [--no-headers]

  # Get tag detail
  nct get tag tagName -i imageName [-r registry] [-o json or yaml] [--no-headers]
```

## Delete
```sh
# Delete Image or Tags from NCR
nct delete [image|tags]

Examples:
  # Delete image
  nct delete image [-r registry] [--dry-run] [-y]
  
  # Delete tags
  nct delete tag [...tagName] -i imageName [-r registry] [--dry-run] [-y]
  nct delete tags [...tagName] -i imageName [-r registry] [--dry-run] [-y]
  nct delete tags [--all] [--exclude-recent number] [-r registry] [--dry-run] [-y]
  nct delete tags [-f filePath] [--dry-run] [-y]
```

## Delete list file example
```yaml
- image: <image>
  tags:
  - <tag>
  - <tag>
- image: <image>
  all: yes
  exclude-recent: 1
- image: <image>
  all: yes
  dry-run: true
```