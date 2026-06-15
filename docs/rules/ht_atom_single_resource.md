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

- Detection is path-based: any module whose file paths contain `/atoms/` is treated as an atom.
- `dynamic` blocks are not resources and are not counted.
- Molecules (`/molecules/`) and data modules (`/data/`) are not checked.
