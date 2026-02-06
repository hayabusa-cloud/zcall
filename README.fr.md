# zcall

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/zcall.svg)](https://pkg.go.dev/code.hybscloud.com/zcall)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/zcall)](https://goreportcard.com/report/github.com/hayabusa-cloud/zcall)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/zcall/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/zcall)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Primitives syscall sans surcharge pour Linux, Darwin et FreeBSD en Go.

Langue : [English](./README.md) | [简体中文](./README.zh-CN.md) | [Español](./README.es.md) | [日本語](./README.ja.md) | **Français**

## Aperçu

`zcall` fournit des points d'entrée syscall bruts qui contournent la machinerie syscall du runtime Go (`entersyscall`/`exitsyscall`). Cela élimine la latence des hooks du scheduler, idéal pour les chemins I/O à faible latence comme la soumission `io_uring`.

### Caractéristiques Principales

- **Zéro Surcharge** : Invocation directe du kernel via assembleur brut
- **Multi-Architecture** : Supporte `linux/amd64`, `linux/arm64`, `linux/riscv64`, `linux/loong64`, `darwin/arm64`, `freebsd/amd64`
- **Sémantique Brute** : Retourne le résultat kernel et errno directement

## Installation

```bash
go get code.hybscloud.com/zcall
```

## Exemple

### I/O Basique

```go
// Écrire sur stdout
msg := []byte("Hello from zcall!\n")
n, errno := zcall.Write(1, msg)
if errno != 0 {
    fmt.Printf("write failed: %v\n", zcall.Errno(errno))
}
```

### Socket Non-Bloquant

```go
// Créer un socket TCP non-bloquant
fd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK, 0)
if errno != 0 {
    return zcall.Errno(errno)
}
defer zcall.Close(fd)
```

## API

### Syscalls Primitives

```go
// Syscall à 4 arguments
Syscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)

// Syscall à 6 arguments
Syscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
```

### Wrappers de Commodité

| Catégorie | Fonctions |
|-----------|-----------|
| I/O Basique | `Read`, `Write`, `Close` |
| I/O Vectorisé | `Readv`, `Writev`, `Preadv`, `Pwritev`, `Preadv2`, `Pwritev2` |
| Socket | `Socket`, `Bind`, `Listen`, `Accept`, `Accept4`, `Connect`, `Shutdown` |
| Socket I/O | `Sendto`, `Recvfrom`, `Sendmsg`, `Recvmsg`, `Sendmmsg`, `Recvmmsg` |
| Mémoire | `Mmap`, `Munmap`, `MemfdCreate` |
| Timers | `TimerfdCreate`, `TimerfdSettime`, `TimerfdGettime` |
| Événements | `Eventfd2`, `Signalfd4` |
| Zero-copy | `Splice`, `Tee`, `Vmsplice`, `Pipe2` |
| io_uring | `IoUringSetup`, `IoUringEnter`, `IoUringRegister` |

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                 Application Utilisateur                 │
├─────────────────────────────────────────────────────────┤
│                      zcall API                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │  Syscall4   │  │  Syscall6   │  │ API Commodité   │  │
│  └──────┬──────┘  └──────┬──────┘  └────────┬────────┘  │
├─────────┴────────────────┴─────────────────┴────────────┤
│                  internal/asm_*.s                       │
│        Assembleur Brut (SYSCALL / SVC / ECALL)          │
├─────────────────────────────────────────────────────────┤
│                   Système d'Exploitation                │
└─────────────────────────────────────────────────────────┘
```

## Plateformes Supportées

| Architecture | Statut | Instruction |
|--------------|--------|-------------|
| linux/amd64  | ✅ Supporté | `SYSCALL` |
| linux/arm64  | ✅ Supporté | `SVC #0` |
| linux/riscv64 | ✅ Supporté | `ECALL` |
| linux/loong64 | ✅ Supporté | `SYSCALL` |
| darwin/arm64 | ✅ Supporté | `SVC #0x80` |
| freebsd/amd64 | ✅ Supporté | `SYSCALL` |

## Considérations de Sécurité

Puisque `zcall` contourne les hooks du scheduler Go :

1. **Appels Non-Bloquants** : Préférez les syscalls non-bloquants avec des vérifications de disponibilité
2. **Validité des Pointeurs** : Assurez-vous que les pointeurs restent valides pendant l'exécution du syscall
3. **Gestion des Erreurs** : Vérifiez errno ; utilisez `zcall.Errno(errno)` pour la conversion d'erreurs

## Licence

MIT — voir [LICENSE](./LICENSE).

©2025 Hayabusa Cloud Co., Ltd.
