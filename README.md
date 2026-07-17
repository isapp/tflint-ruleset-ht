# tflint-ruleset-ht

A custom tflint ruleset enforcing HT Terraform module conventions.

## Rules

| Rule | Severity | Description |
|------|----------|-------------|
| `ht_key_attributes` | ERROR | Identifying resource attributes must appear first, followed by remaining attributes sorted A-Z |
| `ht_managed_repo_ref` | WARNING | `repository` values that match a `managed-repo` module's `repo_name` should reference the module output, not hardcode the string |
| `ht_module_source` | ERROR | Module sources referencing `ht-terraform-modules` must use `git::https://`, include `?ref=`, and omit `//subdir` notation |
| `ht_no_hardcoded_account_id` | WARNING | String literals shaped like a 12-digit AWS account ID should be a variable or `data.aws_caller_identity` |
| `ht_no_hardcoded_arn` | WARNING | Hardcoded ARNs should reference the owning resource's `.arn` (AWS-managed `arn:aws:iam::aws:*` excepted) |
| `ht_no_hardcoded_resource_id` | WARNING | String literals shaped like an AWS resource ID (`vpc-`, `subnet-`, `sg-`, …) should come from a data source or reference |
| `ht_variable_field_order` | ERROR | Fields within each variable block must be ordered: `type`, `default`, `description` |
| `ht_variable_location` | ERROR | All `variable` blocks must be defined in `inputs.tf` |
| `ht_variable_order` | ERROR | Variables in `inputs.tf` must be sorted A-Z within each section (required first, then optional) |
| `ht_variable_section_order` | ERROR | In `inputs.tf`, required and optional variables must be separated by a `DEFAULTS` comment, with required variables before it and optional variables after it |

## Rule Details

### `ht_key_attributes`

For resource types with a configured identifying attribute (e.g. `name` for `aws_s3_bucket`, `bucket` + `key` for `aws_s3_object`), the key attribute(s) must appear before all other attributes. Non-key attributes after the key section must be sorted A-Z.

Violation: a non-key attribute appears on an earlier line than the last key attribute, or non-key attributes after the key section are out of alphabetical order.

### `ht_module_source`

Module `source` values referencing `github.com/isapp/ht-terraform-modules` must follow three conventions:

1. **HTTPS scheme** — `git::https://` is required; `git::ssh://` breaks CI without a deploy key
2. **No `//subdir` notation** — `terraform-module-releaser` creates flat-commit tags with module files at repo root; subdirectories don't exist at tag time and `terraform init` will fail
3. **Pinned `?ref=`** — floating sources are non-reproducible and must be rejected

Sources that do not reference `github.com/isapp/ht-terraform-modules` (local paths, registry sources, other git repos) are ignored by this rule.

See [`docs/rules/ht_module_source.md`](docs/rules/ht_module_source.md) for full details and examples.

### Prefer references over hardcoded values (WARNING)

Four rules discourage literals where a Terraform *reference* belongs. All are `WARNING` severity, so with tflint's default `--minimum-failure-severity=error` they are surfaced but do not fail CI.

- **`ht_managed_repo_ref`** — a `repository` string literal that matches the `repo_name` of a `managed-repo` module in the same configuration should reference `split("/", module.<x>.repository_full_name)[1]` instead. See [`docs/rules/ht_managed_repo_ref.md`](docs/rules/ht_managed_repo_ref.md).
- **`ht_no_hardcoded_account_id`** — flags exact 12-digit strings (AWS account ID shape). See [`docs/rules/ht_no_hardcoded_account_id.md`](docs/rules/ht_no_hardcoded_account_id.md).
- **`ht_no_hardcoded_arn`** — flags `arn:aws...` strings, excepting AWS-managed `arn:aws:iam::aws:*`. See [`docs/rules/ht_no_hardcoded_arn.md`](docs/rules/ht_no_hardcoded_arn.md).
- **`ht_no_hardcoded_resource_id`** — flags AWS resource-ID-shaped strings (`vpc-`, `subnet-`, `sg-`, `ami-`, …). See [`docs/rules/ht_no_hardcoded_resource_id.md`](docs/rules/ht_no_hardcoded_resource_id.md).

The account-id / ARN / resource-id rules are format smells: tflint cannot resolve the correct reference (it usually lives in remote state), so they nudge toward a variable or data source rather than resolving it.

#### Suppressing a finding

There is no per-rule config in v1. Suppress a specific finding with tflint's native annotation on the preceding line:

```hcl
# tflint-ignore: ht_no_hardcoded_account_id
account_map = { prod = "123456789012" }
```

### `ht_variable_field_order`

Within any `variable` block, attributes must appear in the order: `type` → `default` → `description`.

Violation: `type` appears after `default` or `description`; `default` appears after `description`.

### `ht_variable_location`

All `variable` blocks must be declared in a file named `inputs.tf`.

Violation: a `variable` block is found in any file other than `inputs.tf`.

### `ht_variable_order`

In `inputs.tf`, required variables (no `default`) must all appear before optional variables (have `default`). Within each group, variables must be sorted alphabetically.

Violation: a required variable appears after an optional one, or variables within a group are out of alphabetical order.

### `ht_variable_section_order`

In `inputs.tf` files containing both required and optional variables, a separator comment containing the word `DEFAULTS` (e.g. `### DEFAULTS ###`) must be present. All required variables must appear before this separator and all optional variables must appear after it.

Violation: the separator comment is missing; a required variable appears below the separator; an optional variable appears above the separator.
