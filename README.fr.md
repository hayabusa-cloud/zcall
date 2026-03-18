# zcall

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/zcall.svg)](https://pkg.go.dev/code.hybscloud.com/zcall)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/zcall)](https://goreportcard.com/report/github.com/hayabusa-cloud/zcall)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/zcall/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/zcall)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Primitives syscall de bas niveau pour Go sur Linux, Darwin et FreeBSD expérimental.

Langue : [English](./README.md) | [简体中文](./README.zh-CN.md) | [Español](./README.es.md) | [日本語](./README.ja.md) | **Français**

## Aperçu

`zcall` fournit des points d'entrée syscall de bas niveau qui contournent la mécanique syscall du runtime Go (`entersyscall`/`exitsyscall`). Il vise les chemins de code qui nécessitent un contrôle direct des frontières de syscall, comme la soumission `io_uring`.

### Propriétés

- Points d'entrée syscall directs, exposés via des stubs assembleur pour chaque plateforme.
- Sémantique de retour proche du noyau : les wrappers renvoient le résultat brut avec `errno`.
- Numéros de syscall et constantes définis en interne, sans dépendre de `syscall` ni de `x/sys/unix`.
- Implémentations spécifiques à chaque plateforme sélectionnées via des build tags.

### Portée opérationnelle

- Entrées primitives pour les syscalls à 4 et 6 arguments.
- Wrappers pour les opérations courantes d'I/O, de socket, de mémoire et d'options de socket.
- Helpers Linux pour l'interrogation des interfaces réseau, les descripteurs spéciaux, le zero-copy et `io_uring`.

## Installation

```bash
go get code.hybscloud.com/zcall
```

## Utilisation

```go
msg := []byte("Hello from zcall!\n")
// Écriture directe au kernel sur stdout
_, _ = zcall.Write(1, msg)
```

## Exemples d'utilisation

### Vérifier `errno` directement

```go
msg := []byte("hello\n")
n, errno := zcall.Write(1, msg)
if errno != 0 {
	return zcall.Errno(errno)
}
```

### Créer un socket non bloquant

```go
fd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
if errno != 0 {
	return zcall.Errno(errno)
}
defer zcall.Close(fd)
```

### Linux : lister les interfaces réseau

```go
ifaces, err := zcall.Interfaces()
if err != nil {
	return err
}
```

## Types/Opérations

### Syscalls primitives

```go
// Syscall à 4 arguments
Syscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)

// Syscall à 6 arguments
Syscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
```

### Disponibilité des wrappers

#### Disponible sur Linux et Darwin

| Catégorie | Fonctions |
|-----------|-----------|
| I/O Basique | `Read`, `Write`, `Close`, `Ioctl` |
| I/O Vectorisé | `Readv`, `Writev`, `Preadv`, `Pwritev` |
| Socket | `Socket`, `Bind`, `Listen`, `Accept`, `Connect`, `Shutdown`, `Socketpair` |
| Options de Socket | `Setsockopt`, `Getsockopt`, `Getsockname`, `Getpeername` |
| Socket I/O | `Sendto`, `Recvfrom`, `Sendmsg`, `Recvmsg` |
| Mémoire | `Mmap`, `Munmap`, `Pipe2` |

#### Linux uniquement

| Catégorie | Fonctions |
|-----------|-----------|
| I/O Vectorisé | `Preadv2`, `Pwritev2` |
| Socket | `Accept4` |
| Socket I/O | `Sendmmsg`, `Recvmmsg` |
| Réseau | `Interfaces`, `InterfaceByName`, `InterfaceByIndex` |
| Mémoire | `MemfdCreate` |
| Timers | `TimerfdCreate`, `TimerfdSettime`, `TimerfdGettime` |
| Événements | `Eventfd2`, `Signalfd4` |
| Processus | `PidfdOpen`, `PidfdGetfd`, `PidfdSendSignal` |
| Zero-copy | `Splice`, `Tee`, `Vmsplice` |
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

## Support des plateformes

| Architecture | Statut | Instruction |
|--------------|--------|-------------|
| linux/amd64  | ✅ Supporté | `SYSCALL` |
| linux/arm64  | ✅ Supporté | `SVC #0` |
| linux/riscv64 | ✅ Supporté | `ECALL` |
| linux/loong64 | ✅ Supporté | `SYSCALL` |
| darwin/arm64 | ✅ Supporté | `SVC #0x80` |
| freebsd/amd64 | ⚠ Expérimental, non testé | `SYSCALL` |

## Licence

MIT — voir [LICENSE](./LICENSE).

©2025 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com)
