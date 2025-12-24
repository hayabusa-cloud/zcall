# zcall

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/zcall.svg)](https://pkg.go.dev/code.hybscloud.com/zcall)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/zcall)](https://goreportcard.com/report/github.com/hayabusa-cloud/zcall)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/zcall/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/zcall)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Zero-overhead syscall primitives for Linux and Darwin in Go.

Language: **English** | [简体中文](./README.zh-CN.md) | [Español](./README.es.md) | [日本語](./README.ja.md) | [Français](./README.fr.md)

## Overview

`zcall` provides raw syscall entry points that bypass Go's runtime syscall machinery (`entersyscall`/`exitsyscall`). This eliminates scheduler hook latency, making it ideal for low-latency I/O paths such as `io_uring` submission.

### Key Features

- **Zero Overhead**: Direct kernel invocation via raw assembly
- **Multi-Architecture**: Supports `linux/amd64`, `linux/arm64`, `linux/riscv64`, `linux/loong64`, `darwin/arm64`
- **Raw Semantics**: Returns kernel result and errno directly

## Installation

```bash
go get code.hybscloud.com/zcall
```

## Quick Start

```go
msg := []byte("Hello from zcall!\n")
// Direct kernel write to stdout
zcall.Write(1, msg)
```

## API

### Primitive Syscalls

```go
// 4-argument syscall
Syscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)

// 6-argument syscall
Syscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
```

### Convenience Wrappers

| Category | Functions |
|----------|-----------|
| Basic I/O | `Read`, `Write`, `Close` |
| Vectored I/O | `Readv`, `Writev`, `Preadv`, `Pwritev`, `Preadv2`, `Pwritev2` |
| Socket | `Socket`, `Bind`, `Listen`, `Accept`, `Accept4`, `Connect`, `Shutdown` |
| Socket I/O | `Sendto`, `Recvfrom`, `Sendmsg`, `Recvmsg`, `Sendmmsg`, `Recvmmsg` |
| Memory | `Mmap`, `Munmap`, `MemfdCreate` |
| Timers | `TimerfdCreate`, `TimerfdSettime`, `TimerfdGettime` |
| Events | `Eventfd2`, `Signalfd4` |
| Zero-copy | `Splice`, `Tee`, `Vmsplice`, `Pipe2` |
| io_uring | `IoUringSetup`, `IoUringEnter`, `IoUringRegister` |

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    User Application                     │
├─────────────────────────────────────────────────────────┤
│                      zcall API                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │  Syscall4   │  │  Syscall6   │  │ Convenience API │  │
│  └──────┬──────┘  └──────┬──────┘  └────────┬────────┘  │
├─────────┴────────────────┴─────────────────┴────────────┤
│                  internal/asm_*.s                       │
│         Raw Assembly (SYSCALL / SVC / ECALL)            │
├─────────────────────────────────────────────────────────┤
│                     Operating System                    │
└─────────────────────────────────────────────────────────┘
```

## Supported Platforms

| Architecture | Status | Instruction |
|--------------|--------|-------------|
| linux/amd64  | ✅ Supported | `SYSCALL` |
| linux/arm64  | ✅ Supported | `SVC #0` |
| linux/riscv64 | ✅ Supported | `ECALL` |
| linux/loong64 | ✅ Supported | `SYSCALL` |
| darwin/arm64 | ✅ Supported | `SVC #0x80` |

## Safety Considerations

Since `zcall` bypasses Go's scheduler hooks:

1. **Non-blocking Calls**: Prefer non-blocking syscalls with proper readiness checks
2. **Pointer Validity**: Ensure pointers remain valid during syscall execution
3. **Error Handling**: Check errno; use `zcall.Errno(errno)` for error conversion

## License

MIT — see [LICENSE](./LICENSE).

©2025 Hayabusa Cloud Co., Ltd.
