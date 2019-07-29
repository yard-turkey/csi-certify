# Contributing

## Issue/Feature Tracking
  - CSI Certification [Trello board](https://trello.com/b/qKafI4m4/csi-certificiation-suite) 
  
  View the trello board for user feature stories that have been planned out.
## Development 
  Development of this framework is split into three main features

* [External TestDriver](https://github.com/yard-turkey/csi-certify/blob/master/pkg/certify/external/external.go) (main focus)
  * Running E2E tests against a CSI plugin using only info given from a Driver Definition YAML file
* [TestDriver](https://github.com/yard-turkey/csi-certify/tree/master/pkg/certify/driver)
  * Running E2E tests by implementing your own custom TestDriver in GoLang for your CSI plugin if your plugin does not support Dynamic Provisioning (does not have an external provisioner)
* [External Bash TestDriver](https://github.com/yard-turkey/csi-certify/tree/master/pkg/certify/external-bash)
  * Running E2E tests by implementing your own custom TestDriver in Bash for your CSI plugin if your plugin does not support Dynamic provisioning (does not have an external provisioner)

## Project Tree:

```
.
└── csi-certify
    ├── certify
    ├── cmd
    │   └── certify
    │       └── certify_test.go
    ├── Gopkg.lock
    ├── Gopkg.toml
    ├── LICENSE
    ├── pkg
    │   └── certify
    │       ├── certify.go
    │       ├── driver
    │       │   ├── hostpath_driver.go
    │       │   ├── manifests
    │       │   │   ├── hostpath
    │       │   │   │   ├── attacher-rbac.yaml
    │       │   │   │   ├── csi-hostpath-attacher.yaml
    │       │   │   │   ├── csi-hostpathplugin.yaml
    │       │   │   │   ├── csi-hostpath-provisioner.yaml
    │       │   │   │   ├── driver-registrar-rbac.yaml
    │       │   │   │   ├── e2e-test-rbac.yaml
    │       │   │   │   └── provisioner-rbac.yaml
    │       │   │   └── nfs
    │       │   │       ├── csi-attacher-nfsplugin.yaml
    │       │   │       ├── csi-attacher-rbac.yaml
    │       │   │       ├── csi-nodeplugin-nfsplugin.yaml
    │       │   │       └── csi-nodeplugin-rbac.yaml
    │       │   └── nfs_driver.go
    │       ├── external
    │       │   ├── driver-def.yaml
    │       │   └── external.go
    │       ├── external-bash
    │       │   ├── driver-call.go
    │       │   ├── hostpatheeee
    │       │   ├── hostpath-driver-info.yaml
    │       │   ├── nfs
    │       │   ├── nfs-driver-info.yaml
    │       │   └── server-pod.yaml
    │       ├── test
    │       │   ├── csi_volumes.go
    │       │   └── README.md
    │       └── utils
    │           └── test_utils.go
    ├── README.md
    └── test.txt

```

## Repo Owners:
  - [@screeley44](https://github.com/screeley44)
  - [@childsb](https://github.com/childsb)
  - [@mathu97](https://github.com/mathu97)
