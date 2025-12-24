# zcall

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/zcall.svg)](https://pkg.go.dev/code.hybscloud.com/zcall)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/zcall)](https://goreportcard.com/report/github.com/hayabusa-cloud/zcall)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/zcall/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/zcall)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Primitivas de syscall sin sobrecarga para Linux y Darwin en Go.

Idioma: [English](./README.md) | [简体中文](./README.zh-CN.md) | **Español** | [日本語](./README.ja.md) | [Français](./README.fr.md)

## Descripción

`zcall` proporciona puntos de entrada de syscall sin procesar que evitan la maquinaria de syscall del runtime de Go (`entersyscall`/`exitsyscall`). Esto elimina la latencia de los hooks del scheduler, ideal para rutas de I/O de baja latencia como la submisión de `io_uring`.

### Características Principales

- **Sin Sobrecarga**: Invocación directa al kernel mediante ensamblador
- **Multi-Arquitectura**: Soporta `linux/amd64`, `linux/arm64`, `linux/riscv64`, `linux/loong64`, `darwin/arm64`
- **Semántica Raw**: Retorna resultado del kernel y errno directamente

## Instalación

```bash
go get code.hybscloud.com/zcall
```

## Inicio Rápido

```go
msg := []byte("Hello from zcall!\n")
// Escritura directa al kernel en stdout
zcall.Write(1, msg)
```

## API

### Syscalls Primitivas

```go
// Syscall de 4 argumentos
Syscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)

// Syscall de 6 argumentos
Syscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
```

### Wrappers de Conveniencia

| Categoría | Funciones |
|-----------|-----------|
| I/O Básico | `Read`, `Write`, `Close` |
| I/O Vectorizado | `Readv`, `Writev`, `Preadv`, `Pwritev`, `Preadv2`, `Pwritev2` |
| Socket | `Socket`, `Bind`, `Listen`, `Accept`, `Accept4`, `Connect`, `Shutdown` |
| Socket I/O | `Sendto`, `Recvfrom`, `Sendmsg`, `Recvmsg`, `Sendmmsg`, `Recvmmsg` |
| Memoria | `Mmap`, `Munmap`, `MemfdCreate` |
| Timers | `TimerfdCreate`, `TimerfdSettime`, `TimerfdGettime` |
| Eventos | `Eventfd2`, `Signalfd4` |
| Zero-copy | `Splice`, `Tee`, `Vmsplice`, `Pipe2` |
| io_uring | `IoUringSetup`, `IoUringEnter`, `IoUringRegister` |

## Arquitectura

```
┌─────────────────────────────────────────────────────────┐
│                  Aplicación de Usuario                  │
├─────────────────────────────────────────────────────────┤
│                      zcall API                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │  Syscall4   │  │  Syscall6   │  │ API Conveniencia│  │
│  └──────┬──────┘  └──────┬──────┘  └────────┬────────┘  │
├─────────┴────────────────┴─────────────────┴────────────┤
│                  internal/asm_*.s                       │
│        Ensamblador Raw (SYSCALL / SVC / ECALL)          │
├─────────────────────────────────────────────────────────┤
│                   Sistema Operativo                     │
└─────────────────────────────────────────────────────────┘
```

## Plataformas Soportadas

| Arquitectura | Estado | Instrucción |
|--------------|--------|-------------|
| linux/amd64  | ✅ Soportado | `SYSCALL` |
| linux/arm64  | ✅ Soportado | `SVC #0` |
| linux/riscv64 | ✅ Soportado | `ECALL` |
| linux/loong64 | ✅ Soportado | `SYSCALL` |
| darwin/arm64 | ✅ Soportado | `SVC #0x80` |

## Consideraciones de Seguridad

Dado que `zcall` evita los hooks del scheduler de Go:

1. **Llamadas No Bloqueantes**: Prefiere syscalls no bloqueantes con verificaciones de disponibilidad
2. **Validez de Punteros**: Asegura que los punteros permanezcan válidos durante la ejecución del syscall
3. **Manejo de Errores**: Verifica errno; usa `zcall.Errno(errno)` para conversión de errores

## Licencia

MIT — ver [LICENSE](./LICENSE).

©2025 Hayabusa Cloud Co., Ltd.
