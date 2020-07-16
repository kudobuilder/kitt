# KUDO Index Transfer Tool

![](https://github.com/kudobuilder/kitt/workflows/Continuous%20Integration/badge.svg)
[![](https://img.shields.io/static/v1?label=go.dev&message=reference&color=informational&logo=Go)](https://pkg.go.dev/mod/github.com/kudobuilder/kitt)

`kitt` synchronizes KUDO repositories from an index of _operator package references_.

A list of YAML files describing _operator package references_ is used to create or update a [KUDO operator repository](https://github.com/kudobuilder/kudo/blob/main/keps/0015-repository-management.md) on the local file system.

Example usage:

```shell
kitt update --repository /var/kudo/repo /var/kudo/operators/*.yaml
```

## Operator package references

`kitt` builds the KUDO repository from references to operator packages. Each reference describes how to retrieve one or more versioned packages of an operator.

Consider an operator that is developed in a Git repository. The actual operator package (`operator.yaml`, `parameters.yaml`, ...) is in the `operator` folder of this repository. Tagged versions of this operator are then referenced by the following YAML:

```yaml
apiVersion: index.kudo.dev/v1alpha1
kind: Operator
name: MyOperator
gitSources:
  - name: my-git-repository
    url: https://github.com/example/myoperator.git
versions:
  - operatorVersion: "1.0.0"
    git:
      source: my-git-repository
      directory: operator
      tag: "v1.0.0"
  - operatorVersion: "2.0.0"
    git:
      source: my-git-repository
      directory: operator
      tag: "v2.0.0"
```

Running `kitt update` with this YAML as an argument will check out the referenced Git repository with the specified tags `v1.0.0` and `v2.0.0`, build tarballs from the operator package in the `operator` folder, and add these tarballs to a KUDO repository.
