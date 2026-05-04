# zcall

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/zcall.svg)](https://pkg.go.dev/code.hybscloud.com/zcall)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/zcall)](https://goreportcard.com/report/github.com/hayabusa-cloud/zcall)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/zcall/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/zcall)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Primitivas de syscall de bajo nivel para Go en Linux, Darwin y FreeBSD experimental.

Idioma: [English](./README.md) | [简体中文](./README.zh-CN.md) | **Español** | [日本語](./README.ja.md) | [Français](./README.fr.md)

## Descripción

`zcall` proporciona puntos de entrada de syscall de bajo nivel que evitan la maquinaria de syscall del runtime de Go (`entersyscall`/`exitsyscall`). Está dirigido a rutas de código que necesitan control directo sobre los límites de syscall, como la sumisión de `io_uring`.

### Propiedades

- Puntos de entrada directos de syscall, expuestos mediante stubs de ensamblador para cada plataforma.
- Semántica de retorno de estilo kernel: los wrappers devuelven el resultado bruto junto con `errno`.
- Números de syscall y constantes definidos internamente, sin depender de `syscall` ni de `x/sys/unix`.
- Implementaciones específicas por plataforma seleccionadas mediante build tags.

### Alcance operativo

- Entradas primitivas para syscalls de 4 y 6 argumentos.
- Wrappers para operaciones comunes de I/O, sockets, memoria y opciones de socket.
- Helpers solo para Linux para consulta de enlaces de red, descriptores especiales, operaciones de zero-copy e `io_uring`.

## Instalación

```bash
go get code.hybscloud.com/zcall
```

## Uso

```go
msg := []byte("Hello from zcall!\n")
// Escritura directa al kernel en stdout
_, _ = zcall.Write(1, msg)
```

## Ejemplos de uso

### Comprobar `errno` directamente

```go
msg := []byte("hello\n")
n, errno := zcall.Write(1, msg)
if errno != 0 {
	return zcall.Errno(errno)
}
```

### Crear un socket nonblocking

```go
fd, errno := zcall.Socket(zcall.AF_INET, zcall.SOCK_STREAM|zcall.SOCK_NONBLOCK|zcall.SOCK_CLOEXEC, 0)
if errno != 0 {
	return zcall.Errno(errno)
}
defer zcall.Close(fd)
```

### Linux: listar enlaces de red

```go
links, err := zcall.Links()
if err != nil {
	return err
}
```

## Tipos/Operaciones

### Syscalls primitivas

```go
// Syscall de 4 argumentos
Syscall4(num, a1, a2, a3, a4 uintptr) (r1, errno uintptr)

// Syscall de 6 argumentos
Syscall6(num, a1, a2, a3, a4, a5, a6 uintptr) (r1, errno uintptr)
```

### Disponibilidad de wrappers

#### Disponible en Linux, Darwin y FreeBSD

| Categoría | Funciones |
|-----------|-----------|
| I/O Básico | `Read`, `Write`, `Close` |
| I/O Vectorizado | `Readv`, `Writev`, `Preadv`, `Pwritev` |
| Socket | `Socket`, `Bind`, `Listen`, `Accept`, `Connect`, `Shutdown`, `Socketpair` |
| Opciones de Socket | `Setsockopt`, `Getsockopt`, `Getsockname`, `Getpeername` |
| Socket I/O | `Sendto`, `Recvfrom`, `Sendmsg`, `Recvmsg` |
| Memoria | `Mmap`, `Munmap` |

#### Disponible en Linux y FreeBSD

| Categoría | Funciones |
|-----------|-----------|
| Socket | `Accept4` |
| Socket I/O | `Sendmmsg`, `Recvmmsg` |
| Pipe | `Pipe2` |

#### Solo Linux

| Categoría | Funciones |
|-----------|-----------|
| I/O Básico | `Ioctl` |
| I/O Vectorizado | `Preadv2`, `Pwritev2` |
| Red | `Links`, `LinkByName`, `LinkByIndex` |
| Memoria | `MemfdCreate` |
| Timers | `TimerfdCreate`, `TimerfdSettime`, `TimerfdGettime` |
| Eventos | `Eventfd2`, `Signalfd4` |
| Proceso | `PidfdOpen`, `PidfdGetfd`, `PidfdSendSignal` |
| Zero-copy | `Splice`, `Tee`, `Vmsplice` |
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
| darwin/amd64 | ✅ Soportado | `SYSCALL` |
| darwin/arm64 | ✅ Soportado | `SVC #0x80` |
| freebsd/amd64 | ⚠ Experimental, sin probar | `SYSCALL` |

En Darwin, `Preadv` y `Pwritev` usan bucles raw `pread`/`pwrite` sobre el arreglo iovec proporcionado por el llamador porque la tabla de syscalls de Darwin usada por `zcall` expone `pread`/`pwrite` raw, no números de trap raw para `preadv`/`pwritev`.

## Licencia

MIT — ver [LICENSE](./LICENSE).

©2025 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com)
