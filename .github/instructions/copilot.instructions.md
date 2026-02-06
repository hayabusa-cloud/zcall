# zcall — Copilot Instructions

## Return Convention

All syscall wrappers return `(result, errno uintptr)`:

| errno | Meaning |
|-------|---------|
| `0` | Success |
| `!= 0` | Kernel error number |

Always check `errno != 0` before using result.

## Pointer Safety

- All buffer pointers use `noescape()` to prevent heap escape
- Caller ensures pointer validity for entire syscall duration
- `zcall` makes no copies — raw kernel interface

## Assembly Conventions

| Directive | Purpose |
|-----------|---------|
| `NOSPLIT` | No stack split checks (hot path) |
| `//go:noescape` | Hide pointers from escape analysis |
| `//go:nosplit` | No stack growth in Go stubs |

## Platform Differences

| Platform | Error Detection |
|----------|-----------------|
| Linux | Negative return = -errno (check vs -4095) |
| Darwin | Carry flag set = error |
| FreeBSD | Carry flag set = error |

## Known Vet Warning

`zcall.go:279` — `unsafe.Pointer(r1)` in `Mmap()` triggers vet warning. This is a **false positive**: mmap returns kernel-managed address outside Go heap.

## Review Guidelines

**Only report true mistakes.**

Correct patterns (do NOT flag):
- `noescape(unsafe.Pointer(&buf[0]))` — intentional escape prevention
- `return unsafe.Pointer(r1), errno` in Mmap — kernel address conversion
- Assembly using `NOSPLIT` flag — required for zero-overhead
