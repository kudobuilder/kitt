# KUDO Index Transfer Tool

`kitt` synchronizes KUDO repositories from an index of _operator package references_.

Example usage:

```
kitt update --repository /var/kudo/repo $(git diff --name-only --diff-filter=AM)
```
