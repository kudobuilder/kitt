# KUDO Index Transfer Tool

![](https://github.com/kudobuilder/kitt/workflows/Continuous%20Integration/badge.svg)
[![](https://img.shields.io/static/v1?label=go.dev&message=reference&color=informational&logo=Go)](https://pkg.go.dev/mod/github.com/kudobuilder/kitt)

`kitt` synchronizes KUDO repositories from an index of _operator package references_.

Example usage:

```
kitt update --repository /var/kudo/repo $(git diff --name-only --diff-filter=AM)
```
