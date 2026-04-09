# tflint-ruleset-ht

A custom tflint ruleset enforcing HT Terraform module conventions.

## Rules

| Rule | Severity | Description |
|------|----------|-------------|
| `ht_key_attributes` | ERROR | Identifying resource attributes must appear first, followed by remaining attributes sorted A-Z |
| `ht_variable_field_order` | ERROR | Fields within each variable block must be ordered: `type`, `default`, `description` |
| `ht_variable_location` | ERROR | All `variable` blocks must be defined in `inputs.tf` |
| `ht_variable_order` | ERROR | Variables in `inputs.tf` must be sorted A-Z within each section (required first, then optional) |
| `ht_variable_section_order` | ERROR | In `inputs.tf`, required and optional variables must be separated by a `DEFAULTS` comment, with required variables before it and optional variables after it |

## Rule Details

### `ht_key_attributes`

For resource types with a configured identifying attribute (e.g. `name` for `aws_s3_bucket`, `bucket` + `key` for `aws_s3_object`), the key attribute(s) must appear before all other attributes. Non-key attributes after the key section must be sorted A-Z.

Violation: a non-key attribute appears on an earlier line than the last key attribute, or non-key attributes after the key section are out of alphabetical order.

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
