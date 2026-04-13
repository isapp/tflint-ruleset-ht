# `ht_module_source`

**Severity:** ERROR  
**Scope:** module `source` values referencing `github.com/isapp/ht-terraform-modules`

## Summary

Enforces three conventions for module sources that reference `ht-terraform-modules`:

1. Must use `git::https://` — not `git::ssh://`
2. Must not use `//subdir` notation
3. Must be pinned with `?ref=`

## Background

`terraform-module-releaser` creates flat-commit tags: only the changed module's files are committed at the tag root — there is no directory tree. The `//subdir` notation fails at `terraform init` because the subdirectory does not exist in the flat commit.

SSH sources (`git::ssh://`) require a deploy key in CI and are not supported in our pipelines. All sources must use HTTPS.

Floating sources (no `?ref=`) resolve to HEAD at init time, producing non-reproducible builds.

## Rule checks

| Check | Trigger | Message |
|-------|---------|----------|
| Wrong scheme | source contains `git::ssh://` | `module source must use git::https://, not git::ssh://` |
| Subdir notation | source contains `//` after stripping scheme | `module source must not use //subdir notation — ht-terraform-modules tags are flat commits` |
| Missing pin | source does not contain `?ref=` | `module source must be pinned with ?ref=<tag>` |

All three checks are independent — a single source can emit up to three issues.

## Scope

Only fires on sources containing `github.com/isapp/ht-terraform-modules`. The following are **ignored**:

- Local relative paths (`../../`, `./`)
- Terraform registry sources (`namespace/module/provider`)
- Other git repos

## Examples

### Valid

```hcl
# atom
module "topic" {
  source = "git::https://github.com/isapp/ht-terraform-modules?ref=modules/atoms/sns-topic/v1.0.0"
}

# molecule
module "alarm" {
  source = "git::https://github.com/isapp/ht-terraform-modules?ref=modules/molecules/cloudwatch-email-alarm/v1.1.0"
}
```

### Invalid — SSH scheme

```hcl
module "topic" {
  source = "git::ssh://git@github.com/isapp/ht-terraform-modules?ref=modules/atoms/sns-topic/v1.0.0"
  # ERROR: module source must use git::https://, not git::ssh://
}
```

### Invalid — `//subdir` notation

```hcl
module "topic" {
  source = "git::https://github.com/isapp/ht-terraform-modules//modules/atoms/sns-topic?ref=modules/atoms/sns-topic/v1.0.0"
  # ERROR: module source must not use //subdir notation — ht-terraform-modules tags are flat commits
}
```

### Invalid — missing `?ref=`

```hcl
module "topic" {
  source = "git::https://github.com/isapp/ht-terraform-modules"
  # ERROR: module source must be pinned with ?ref=<tag>
}
```
