# zcall

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/zcall.svg)](https://pkg.go.dev/code.hybscloud.com/zcall)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/zcall)](https://goreportcard.com/report/github.com/hayabusa-cloud/zcall)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/zcall/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/zcall)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Low-level syscall primitives for Go on Linux, Darwin, and experimental FreeBSD.

Language: **English** | [简体中文](./README.zh-CN.md) | [Español](./README.es.md) | [日本語](./README.ja.md) | [Français](./README.fr.md)

## Overview

`zcall` provides low-level syscall entry points that bypass Go's runtime syscall machinery (`entersyscall`/`exitsyscall`). It targets code paths that need direct control over syscall boundaries, such as `io_uring` submission.

### Properties

- Direct syscall entry points, provided through assembly stubs for each platform.
- Kernel-style return values, with wrappers returning the raw result alongside `errno`.
- Syscall numbers and related constants defined internally, without depending on `syscall` or `x/sys/unix`.
- Platform-specific implementations selected through build tags.

### Operations

- Primitive entry points for 4-argument and 6-argument syscalls.
- Wrappers for common I/O, socket, memory, and socket-option operations.
- Linux-specific helpers for network interface queries, special file descriptors, zero-copy operations, and `io_uring`.

## Installation

```bash
go get code.hybscloud.com/zcall
```

## Usage

```go
msg := []byte("Hello from zcall!\n")
// Direct kernel write to stdout
_, _ = zcall.Write(1, msg)
```

## Examples

### Check `errno` directly

```go
msg := []byte("hello\n")
n, errno := zcall.Write(1, msg)
if errno != 0 {
	return zcall.Errno(errno)
}
```

### Create a nonblocking socket

```go
fd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
if errno != 0 {
	return zcall.Errno(errno)
}
defer zcall.Close(fd)
```

### Linux: list network interfaces

```go
ifaces, err := zcall.Interfaces()
if err != nil {
	return err
}
```

## Types/Operations

### Primitive syscalls

```go
// 4-argument syscall
Syscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)

// 6-argument syscall
Syscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
```

### Wrapper availability

#### Available on Linux and Darwin

| Category | Functions |
|----------|-----------|
| Basic I/O | `Read`, `Write`, `Close`, `Ioctl` |
| Vectored I/O | `Readv`, `Writev`, `Preadv`, `Pwritev` |
| Socket | `Socket`, `Bind`, `Listen`, `Accept`, `Connect`, `Shutdown`, `Socketpair` |
| Socket Options | `Setsockopt`, `Getsockopt`, `Getsockname`, `Getpeername` |
| Socket I/O | `Sendto`, `Recvfrom`, `Sendmsg`, `Recvmsg` |
| Memory | `Mmap`, `Munmap`, `Pipe2` |

#### Linux-only wrappers

| Category | Functions |
|----------|-----------|
| Vectored I/O | `Preadv2`, `Pwritev2` |
| Socket | `Accept4` |
| Socket I/O | `Sendmmsg`, `Recvmmsg` |
| Network | `Interfaces`, `InterfaceByName`, `InterfaceByIndex` |
| Memory | `MemfdCreate` |
| Timers | `TimerfdCreate`, `TimerfdSettime`, `TimerfdGettime` |
| Events | `Eventfd2`, `Signalfd4` |
| Process | `PidfdOpen`, `PidfdGetfd`, `PidfdSendSignal` |
| Zero-copy | `Splice`, `Tee`, `Vmsplice` |
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

## Platform Support

| Architecture | Status | Instruction |
|--------------|--------|-------------|
| linux/amd64  | ✅ Supported | `SYSCALL` |
| linux/arm64  | ✅ Supported | `SVC #0` |
| linux/riscv64 | ✅ Supported | `ECALL` |
| linux/loong64 | ✅ Supported | `SYSCALL` |
| darwin/arm64 | ✅ Supported | `SVC #0x80` |
| freebsd/amd64 | ⚠ Experimental, untested | `SYSCALL` |

## License

MIT — see [LICENSE](./LICENSE).

©2025 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com)
