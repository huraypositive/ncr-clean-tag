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
  nct delete tag [...tagName] -i imageName [--exclude-tags=tag1,tag2,...] [-r registry] [--dry-run] [-y]
  nct delete tags [...tagName] -i imageName [--exclude-tags=tag1,tag2,...] [-r registry] [--dry-run] [-y]
  nct delete tags [--all] [--exclude-recent number] [--exclude-tags=tag1,tag2,...] [-r registry] [--dry-run] [-y]
  nct delete tags [-f filePath] [--exclude-recent number] [--exclude-tags=tag1,tag2,...] [-r registry] [--dry-run] [-y]
```

## Configmap example
```yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: delete-tag-list
data:
  delete-tag-list.yaml: |-
    ---
    - registry: <registry> # -r flag or DEFAULT_REGISTRY env can be used instead.
      image: <image>
      tags:
      - <tag>
      - <tag>
    - image: <image>
      all: yes
      dry-run: true
    - image: <image>
      all: yes
      exclude-recent: 1
    - image: <image>
      all: yes
      exclude-tags:
      - <excludeTag>
      - <excludeTag>
    - registry: <registry>
      image: <image>
      all: yes
      exclude-recent: 5
      exclude-tags:
      - <excludeTag>
      - <excludeTag>
      dry-run: true
```

## Secret example
```yaml
---
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: ncr-clean-tag
data:
  DEFAULT_REGISTRY: <defaultRegistry>
  NCR_ACCESS_KEY: <accessKey>
  NCR_SECRET_KEY: <secretKey>
```

## Cronjob example
```yaml
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: ncr-clean-tag
spec:
  schedule: "30 15 * * 0"
  jobTemplate:
    spec:
      activeDeadlineSeconds: 300
      ttlSecondsAfterFinished: 86400
      template:
        spec:
          restartPolicy: Never
          containers:
          - image: nct:v2.0.1
            name: ncr-clean-tag
            envFrom:
            - secretRef:
                name: ncr-clean-tag
            args:
            - delete
            - tags
            - -f
            - /var/run/configmaps/delete-tag-list.yaml
            - -y
            volumeMounts:
            - mountPath: /var/run/configmaps
              name: delete-tag-list
              readOnly: true
          volumes:
          - configMap:
              name: delete-tag-list
            name: delete-tag-list
```