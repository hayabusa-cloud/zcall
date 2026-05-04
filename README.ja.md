# zcall

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/zcall.svg)](https://pkg.go.dev/code.hybscloud.com/zcall)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/zcall)](https://goreportcard.com/report/github.com/hayabusa-cloud/zcall)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/zcall/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/zcall)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Linux、Darwin、および実験的な FreeBSD 向けの Go 低レベル syscall プリミティブ。

言語: [English](./README.md) | [简体中文](./README.zh-CN.md) | [Español](./README.es.md) | **日本語** | [Français](./README.fr.md)

## 概要

`zcall` は Go ランタイムの syscall 機構（`entersyscall`/`exitsyscall`）をバイパスする低レベルな syscall エントリポイントを提供します。`io_uring` サブミッションのように syscall 境界を直接制御する必要があるコードパスを対象にしています。

### 特性

- 各プラットフォーム向けのアセンブリスタブで提供される直接 syscall エントリポイント。
- カーネルに近い戻り値の形を保ち、ラッパーは生の結果値と `errno` を返します。
- `syscall` や `x/sys/unix` に依存せず、syscall 番号と定数を内部で定義します。
- ビルドタグでプラットフォーム別の実装を切り替えます。

### 操作範囲

- 4 引数および 6 引数 syscall のプリミティブ入口。
- 一般的な I/O、ソケット、メモリ、ソケットオプション向けのラッパー。
- Linux 専用のネットワークリンク参照、特殊 FD、ゼロコピー、`io_uring` 用ヘルパー。

## インストール

```bash
go get code.hybscloud.com/zcall
```

## Usage

```go
msg := []byte("Hello from zcall!\n")
// stdout への直接カーネル書き込み
_, _ = zcall.Write(1, msg)
```

## 使用例

### `errno` を直接確認する

```go
msg := []byte("hello\n")
n, errno := zcall.Write(1, msg)
if errno != 0 {
	return zcall.Errno(errno)
}
```

### nonblocking ソケットを作成する

```go
fd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
if errno != 0 {
	return zcall.Errno(errno)
}
defer zcall.Close(fd)
```

### Linux: ネットワークリンクを列挙する

```go
links, err := zcall.Links()
if err != nil {
	return err
}
```

## Types/Operations

### プリミティブ Syscall

```go
// 4 引数 syscall
Syscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)

// 6 引数 syscall
Syscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
```

### ラッパーの提供状況

#### Linux、Darwin、FreeBSD で利用可能

| カテゴリ | 関数 |
|----------|------|
| 基本 I/O | `Read`、`Write`、`Close` |
| ベクタ I/O | `Readv`、`Writev`、`Preadv`、`Pwritev` |
| ソケット | `Socket`、`Bind`、`Listen`、`Accept`、`Connect`、`Shutdown`、`Socketpair` |
| ソケットオプション | `Setsockopt`、`Getsockopt`、`Getsockname`、`Getpeername` |
| ソケット I/O | `Sendto`、`Recvfrom`、`Sendmsg`、`Recvmsg` |
| メモリ | `Mmap`、`Munmap` |

#### Linux と FreeBSD で利用可能

| カテゴリ | 関数 |
|----------|------|
| ソケット | `Accept4` |
| ソケット I/O | `Sendmmsg`、`Recvmmsg` |
| パイプ | `Pipe2` |

#### Linux 専用

| カテゴリ | 関数 |
|----------|------|
| 基本 I/O | `Ioctl` |
| ベクタ I/O | `Preadv2`、`Pwritev2` |
| ネットワーク | `Links`、`LinkByName`、`LinkByIndex` |
| メモリ | `MemfdCreate` |
| タイマー | `TimerfdCreate`、`TimerfdSettime`、`TimerfdGettime` |
| イベント | `Eventfd2`、`Signalfd4` |
| プロセス | `PidfdOpen`、`PidfdGetfd`、`PidfdSendSignal` |
| ゼロコピー | `Splice`、`Tee`、`Vmsplice` |
| io_uring | `IoUringSetup`、`IoUringEnter`、`IoUringRegister` |

## アーキテクチャ

```
┌─────────────────────────────────────────────────────────┐
│                   ユーザーアプリケーション                 │
├─────────────────────────────────────────────────────────┤
│                      zcall API                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │  Syscall4   │  │  Syscall6   │  │  便利な API     │  │
│  └──────┬──────┘  └──────┬──────┘  └────────┬────────┘  │
├─────────┴────────────────┴─────────────────┴────────────┤
│                  internal/asm_*.s                       │
│         生アセンブリ (SYSCALL / SVC / ECALL)             │
├─────────────────────────────────────────────────────────┤
│               オペレーティングシステム                  │
└─────────────────────────────────────────────────────────┘
```

## Platform Support

| アーキテクチャ | 状態 | 命令 |
|----------------|------|------|
| linux/amd64  | ✅ サポート | `SYSCALL` |
| linux/arm64  | ✅ サポート | `SVC #0` |
| linux/riscv64 | ✅ サポート | `ECALL` |
| linux/loong64 | ✅ サポート | `SYSCALL` |
| darwin/amd64 | ✅ サポート | `SYSCALL` |
| darwin/arm64 | ✅ サポート | `SVC #0x80` |
| freebsd/amd64 | ⚠ 実験的、未検証 | `SYSCALL` |

Darwin では、`Preadv` と `Pwritev` は呼び出し側が渡した iovec 配列に対する raw `pread`/`pwrite` ループとして実装されています。`zcall` が使う Darwin syscall テーブルは raw `pread`/`pwrite` を公開しますが、raw `preadv`/`pwritev` の trap 番号は公開しません。

## ライセンス

MIT — [LICENSE](./LICENSE) を参照。

©2025 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com)
