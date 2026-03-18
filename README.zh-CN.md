# zcall

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/zcall.svg)](https://pkg.go.dev/code.hybscloud.com/zcall)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/zcall)](https://goreportcard.com/report/github.com/hayabusa-cloud/zcall)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/zcall/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/zcall)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

面向 Linux、Darwin 和实验性 FreeBSD 的 Go 底层系统调用原语。

语言：[English](./README.md) | **简体中文** | [Español](./README.es.md) | [日本語](./README.ja.md) | [Français](./README.fr.md)

## 概述

`zcall` 提供绕过 Go 运行时系统调用机制（`entersyscall`/`exitsyscall`）的底层系统调用入口点，面向需要直接控制系统调用边界的代码路径，例如 `io_uring` 提交。

### 属性

- 通过各平台对应的汇编桩提供直接系统调用入口。
- 保留接近内核的返回语义：封装函数返回原始结果值以及 `errno`。
- 在内部定义系统调用号和常量，不依赖 `syscall` 或 `x/sys/unix`。
- 通过构建标签选择平台专用实现。

### 操作范围

- 提供 4 参数与 6 参数系统调用的原始入口。
- 为常见 I/O、套接字、内存和套接字选项操作提供封装函数。
- 提供 Linux 专用辅助函数，用于网络接口查询、特殊文件描述符、零拷贝操作以及 `io_uring`。

## 安装

```bash
go get code.hybscloud.com/zcall
```

## 用法

```go
msg := []byte("Hello from zcall!\n")
// 直接内核写入到标准输出
_, _ = zcall.Write(1, msg)
```

## 使用示例

### 直接检查 `errno`

```go
msg := []byte("hello\n")
n, errno := zcall.Write(1, msg)
if errno != 0 {
	return zcall.Errno(errno)
}
```

### 创建非阻塞 socket

```go
fd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
if errno != 0 {
	return zcall.Errno(errno)
}
defer zcall.Close(fd)
```

### Linux：列出网络接口

```go
ifaces, err := zcall.Interfaces()
if err != nil {
	return err
}
```

## 类型/操作

### 原始 syscall

```go
// 4 参数系统调用
Syscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)

// 6 参数系统调用
Syscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
```

### 封装函数可用性

#### Linux 和 Darwin 可用

| 类别 | 函数 |
|------|------|
| 基础 I/O | `Read`、`Write`、`Close`、`Ioctl` |
| 向量 I/O | `Readv`、`Writev`、`Preadv`、`Pwritev` |
| 套接字 | `Socket`、`Bind`、`Listen`、`Accept`、`Connect`、`Shutdown`、`Socketpair` |
| 套接字选项 | `Setsockopt`、`Getsockopt`、`Getsockname`、`Getpeername` |
| 套接字 I/O | `Sendto`、`Recvfrom`、`Sendmsg`、`Recvmsg` |
| 内存 | `Mmap`、`Munmap`、`Pipe2` |

#### 仅 Linux 可用

| 类别 | 函数 |
|------|------|
| 向量 I/O | `Preadv2`、`Pwritev2` |
| 套接字 | `Accept4` |
| 套接字 I/O | `Sendmmsg`、`Recvmmsg` |
| 网络 | `Interfaces`、`InterfaceByName`、`InterfaceByIndex` |
| 内存 | `MemfdCreate` |
| 定时器 | `TimerfdCreate`、`TimerfdSettime`、`TimerfdGettime` |
| 事件 | `Eventfd2`、`Signalfd4` |
| 进程 | `PidfdOpen`、`PidfdGetfd`、`PidfdSendSignal` |
| 零拷贝 | `Splice`、`Tee`、`Vmsplice` |
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

## 平台支持

| 架构 | 状态 | 指令 |
|------|------|------|
| linux/amd64  | ✅ 支持 | `SYSCALL` |
| linux/arm64  | ✅ 支持 | `SVC #0` |
| linux/riscv64 | ✅ 支持 | `ECALL` |
| linux/loong64 | ✅ 支持 | `SYSCALL` |
| darwin/arm64 | ✅ 支持 | `SVC #0x80` |
| freebsd/amd64 | ⚠ 实验性，未经验证 | `SYSCALL` |

## 许可证

MIT — 见 [LICENSE](./LICENSE)。

©2025 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com)
