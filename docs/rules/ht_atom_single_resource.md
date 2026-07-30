# `ht_atom_single_resource`

**Severity:** ERROR
**Scope:** modules under `modules/*/atoms/` (atom tier)

## Summary

Enforces that every atom module contains exactly one `resource` block. Atoms are 1-to-1 wrappers around a single provider resource — multiple resources indicate molecule-level coupling that belongs at the molecule tier.

## Rule checks

| Check | Trigger | Message |
|-------|---------|---------|
| Wrong resource count | atom path + ≠1 resource blocks | `atom must contain exactly 1 resource block, found N` |

## Notes

- Detection is path-based: a module is treated as an atom when either the process working directory or any filename reported by the runner contains `/atoms/`. Both signals are needed — tflint changes into the module directory before running rules, so in real invocations the filenames are module-relative (`main.tf`) and carry no tier segment, while the plugin-sdk test helper uses its map keys verbatim and so supplies full paths. See `isAtomModule` in `rules/atom_helpers.go`.
- `dynamic` blocks are not resources and are not counted.
- Molecules (`/molecules/`) and data modules (`/data/`) are not checked.
