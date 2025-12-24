# zcall

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/zcall.svg)](https://pkg.go.dev/code.hybscloud.com/zcall)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/zcall)](https://goreportcard.com/report/github.com/hayabusa-cloud/zcall)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/zcall/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/zcall)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go 言語向けゼロオーバーヘッド syscall プリミティブ（Linux・Darwin）。

言語: [English](./README.md) | [简体中文](./README.zh-CN.md) | [Español](./README.es.md) | **日本語** | [Français](./README.fr.md)

## 概要

`zcall` は Go ランタイムの syscall 機構（`entersyscall`/`exitsyscall`）をバイパスする生の syscall エントリポイントを提供します。これによりスケジューラフックのレイテンシが排除され、`io_uring` サブミッションなどの低レイテンシ I/O パスに最適です。

### 主な特徴

- **ゼロオーバーヘッド**: 生アセンブリによる直接カーネル呼び出し
- **マルチアーキテクチャ**: `linux/amd64`、`linux/arm64`、`linux/riscv64`、`linux/loong64`、`darwin/arm64` をサポート
- **生セマンティクス**: カーネル結果と errno を直接返却

## インストール

```bash
go get code.hybscloud.com/zcall
```

## クイックスタート

```go
msg := []byte("Hello from zcall!\n")
// stdout への直接カーネル書き込み
zcall.Write(1, msg)
```

## API

### プリミティブ Syscall

```go
// 4 引数 syscall
Syscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)

// 6 引数 syscall
Syscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
```

### 便利なラッパー

| カテゴリ | 関数 |
|----------|------|
| 基本 I/O | `Read`、`Write`、`Close` |
| ベクタ I/O | `Readv`、`Writev`、`Preadv`、`Pwritev`、`Preadv2`、`Pwritev2` |
| ソケット | `Socket`、`Bind`、`Listen`、`Accept`、`Accept4`、`Connect`、`Shutdown` |
| ソケット I/O | `Sendto`、`Recvfrom`、`Sendmsg`、`Recvmsg`、`Sendmmsg`、`Recvmmsg` |
| メモリ | `Mmap`、`Munmap`、`MemfdCreate` |
| タイマー | `TimerfdCreate`、`TimerfdSettime`、`TimerfdGettime` |
| イベント | `Eventfd2`、`Signalfd4` |
| ゼロコピー | `Splice`、`Tee`、`Vmsplice`、`Pipe2` |
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

## サポートプラットフォーム

| アーキテクチャ | 状態 | 命令 |
|----------------|------|------|
| linux/amd64  | ✅ サポート | `SYSCALL` |
| linux/arm64  | ✅ サポート | `SVC #0` |
| linux/riscv64 | ✅ サポート | `ECALL` |
| linux/loong64 | ✅ サポート | `SYSCALL` |
| darwin/arm64 | ✅ サポート | `SVC #0x80` |

## 安全性に関する注意事項

`zcall` は Go のスケジューラフックをバイパスするため：

1. **ノンブロッキング呼び出し**: 適切な準備状態チェックを伴うノンブロッキング syscall を優先
2. **ポインタの有効性**: syscall 実行中はポインタが有効であることを確認
3. **エラー処理**: errno をチェック；エラー変換には `zcall.Errno(errno)` を使用

## ライセンス

MIT — [LICENSE](./LICENSE) を参照。

©2025 Hayabusa Cloud Co., Ltd.
