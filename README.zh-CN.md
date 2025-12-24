# zcall

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/zcall.svg)](https://pkg.go.dev/code.hybscloud.com/zcall)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/zcall)](https://goreportcard.com/report/github.com/hayabusa-cloud/zcall)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/zcall/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/zcall)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go 语言零开销系统调用原语（Linux、Darwin 和 FreeBSD）。

语言：[English](./README.md) | **简体中文** | [Español](./README.es.md) | [日本語](./README.ja.md) | [Français](./README.fr.md)

## 概述

`zcall` 提供绕过 Go 运行时系统调用机制（`entersyscall`/`exitsyscall`）的原始系统调用入口点。这消除了调度器钩子的延迟，非常适合低延迟 I/O 路径，如 `io_uring` 提交。

### 核心特性

- **零开销**：通过原始汇编直接调用内核
- **多架构**：支持 `linux/amd64`、`linux/arm64`、`linux/riscv64`、`linux/loong64`、`darwin/arm64`、`freebsd/amd64`
- **原始语义**：直接返回内核结果和 errno

## 安装

```bash
go get code.hybscloud.com/zcall
```

## 快速开始

```go
msg := []byte("Hello from zcall!\n")
// 直接内核写入到标准输出
zcall.Write(1, msg)
```

## API

### 原始系统调用

```go
// 4 参数系统调用
Syscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)

// 6 参数系统调用
Syscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
```

### 便捷封装

| 类别 | 函数 |
|------|------|
| 基础 I/O | `Read`、`Write`、`Close` |
| 向量 I/O | `Readv`、`Writev`、`Preadv`、`Pwritev`、`Preadv2`、`Pwritev2` |
| 套接字 | `Socket`、`Bind`、`Listen`、`Accept`、`Accept4`、`Connect`、`Shutdown` |
| 套接字 I/O | `Sendto`、`Recvfrom`、`Sendmsg`、`Recvmsg`、`Sendmmsg`、`Recvmmsg` |
| 内存 | `Mmap`、`Munmap`、`MemfdCreate` |
| 定时器 | `TimerfdCreate`、`TimerfdSettime`、`TimerfdGettime` |
| 事件 | `Eventfd2`、`Signalfd4` |
| 零拷贝 | `Splice`、`Tee`、`Vmsplice`、`Pipe2` |
| io_uring | `IoUringSetup`、`IoUringEnter`、`IoUringRegister` |

## 架构

```
┌─────────────────────────────────────────────────────────┐
│                      用户应用程序                        │
├─────────────────────────────────────────────────────────┤
│                      zcall API                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │  Syscall4   │  │  Syscall6   │  │    便捷 API     │  │
│  └──────┬──────┘  └──────┬──────┘  └────────┬────────┘  │
├─────────┴────────────────┴─────────────────┴────────────┤
│                  internal/asm_*.s                       │
│          原始汇编 (SYSCALL / SVC / ECALL)                │
├─────────────────────────────────────────────────────────┤
│                        操作系统                         │
└─────────────────────────────────────────────────────────┘
```

## 支持的平台

| 架构 | 状态 | 指令 |
|------|------|------|
| linux/amd64  | ✅ 支持 | `SYSCALL` |
| linux/arm64  | ✅ 支持 | `SVC #0` |
| linux/riscv64 | ✅ 支持 | `ECALL` |
| linux/loong64 | ✅ 支持 | `SYSCALL` |
| darwin/arm64 | ✅ 支持 | `SVC #0x80` |
| freebsd/amd64 | ✅ 支持 | `SYSCALL` |

## 安全注意事项

由于 `zcall` 绕过了 Go 的调度器钩子：

1. **非阻塞调用**：优先使用非阻塞系统调用并进行适当的就绪检查
2. **指针有效性**：确保指针在系统调用执行期间保持有效
3. **错误处理**：检查 errno；使用 `zcall.Errno(errno)` 进行错误转换

## 许可证

MIT — 见 [LICENSE](./LICENSE)。

©2025 Hayabusa Cloud Co., Ltd.
